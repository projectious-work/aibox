package usecase

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/projectious-work/aibox/internal/contract"
	"github.com/projectious-work/aibox/internal/logging"
	"github.com/projectious-work/aibox/internal/redact"
	"github.com/projectious-work/aibox/internal/schematest"
	"github.com/projectious-work/aibox/internal/settings"
)

// recordedEvents is an Events fake.
type recordedEvents struct {
	mu     sync.Mutex
	events []logging.Event
}

func (r *recordedEvents) Emit(event logging.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, event)
}

func (r *recordedEvents) names() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	var names []string
	for _, event := range r.events {
		names = append(names, event.Name)
	}
	return strings.Join(names, ",")
}

// sinkState is a SinkHealth stub.
type sinkState struct{ failed bool }

func (s sinkState) FileSinkFailed() bool { return s.failed }

type fixture struct {
	home, project string
	environ       []string
	flags         []settings.Flag
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	base := t.TempDir()
	f := &fixture{home: filepath.Join(base, "home"), project: filepath.Join(base, "project")}
	if err := os.MkdirAll(filepath.Join(f.project, ".devcontainer"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(f.home, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f.project, ".devcontainer", "devcontainer.json"), []byte(`{"image":"x"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	f.environ = []string{"HOME=" + f.home}
	return f
}

func (f *fixture) resolve(t *testing.T) *settings.Resolved {
	t.Helper()
	loaded := settings.Load(settings.OSFileSystem{}, settings.Request{Mode: settings.ModeCLI, GOOS: "linux", Environ: f.environ,
		WorkDir: f.project, Flags: f.flags})
	resolved, rejection := loaded.Resolve(settings.ReadProject(f.project))
	if rejection != nil {
		t.Fatalf("%+v", rejection)
	}
	return resolved
}

func (f *fixture) request(effective bool) InspectRequest {
	return InspectRequest{RequestID: "req-1", Initiator: "cli", ProjectRoot: f.project, EffectiveConfig: effective}
}

func TestInspectEventsAndResult(t *testing.T) {
	// LOG-001 (event lists of inspect), CFG-067, CFG-068.
	f := newFixture(t)
	events := &recordedEvents{}
	result := Workspace{Health: sinkState{}}.InspectWorkspace(t.Context(), f.resolve(t), f.request(false), events)
	if result.Outcome != contract.Succeeded || events.names() != "operation.started,operation.validated,operation.finished" {
		t.Fatalf("%+v %s", result, events.names())
	}
	if events.events[0].Fields["projectRoot"] != f.project || events.events[2].Outcome != "succeeded" || events.events[2].ErrorCategory != "" {
		t.Fatalf("%+v", events.events)
	}
	encoded, _ := json.Marshal(result)
	validateResult(t, encoded)
	if strings.Contains(string(encoded), "settingsSources") || strings.Count(string(encoded), `"key"`) != 1 {
		t.Fatalf("view without --effective-config: %s", encoded)
	}

	events = &recordedEvents{}
	missing := InspectRequest{RequestID: "req-2", Initiator: "cli", ProjectRoot: filepath.Join(f.project, "missing")}
	result = Workspace{Health: sinkState{}}.InspectWorkspace(t.Context(), f.resolve(t), missing, events)
	if result.Outcome != contract.Failed || result.ExitCode() != 2 || events.names() != "operation.started,operation.finished" ||
		events.events[1].ErrorCategory != "invalid_input" {
		t.Fatalf("%+v %s", result, events.names())
	}
}

func validateResult(t *testing.T, encoded []byte) {
	t.Helper()
	schema, err := schematest.Load("schemas/operation-result.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := schema.Validate(encoded); err != nil {
		t.Fatalf("%v: %s", err, encoded)
	}
}

func TestEffectiveConfigurationView(t *testing.T) {
	// CFG-060 to CFG-064, EXE-003, RED-001, RED-003 (J-CFG-22 at component level).
	f := newFixture(t)
	canary := "canary-0123456789"
	f.environ = append(f.environ, "EXAMPLE_API_TOKEN="+canary)
	envFile := filepath.Join(f.home, canary+".env")
	if err := os.WriteFile(envFile, []byte("AIBOX_WORKSPACE__THEME=nord\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	f.environ = append(f.environ, "AIBOX_ENV_FILE="+envFile)
	f.flags = []settings.Flag{{Name: "--format", Value: "json"}, {Name: "--log-file", Value: "run.log"}}
	resolved := f.resolve(t)
	redactor := redact.New(f.environ, f.home)
	result := Workspace{Health: sinkState{}, Redactor: redactor}.InspectWorkspace(t.Context(), resolved, f.request(true), &recordedEvents{})
	encoded, _ := json.Marshal(result)
	validateResult(t, encoded)
	data := result.Data.(contract.WorkspaceData)
	if len(data.Configuration) != 16 || data.Configuration[0].Key != "devcontainer.config" {
		t.Fatalf("%+v", data.Configuration)
	}
	byKey := map[string]contract.ConfigurationEntry{}
	for index, entry := range data.Configuration[1:] {
		if entry.Key != settings.Keys[index].Name {
			t.Fatalf("entry %d is %s", index, entry.Key)
		}
		byKey[entry.Key] = entry
	}
	if entry := byKey[settings.OutputFormat]; entry.Value != "json" || entry.Layer != "invocation" || entry.Source != "--format" || !entry.InvocationOnly {
		t.Fatalf("%+v", entry)
	}
	if entry := byKey[settings.ExecutionTimeout]; entry.Value != 30 || entry.Layer != "built_in" || entry.Source != "built-in" {
		t.Fatalf("%+v", entry)
	}
	if entry := byKey[settings.LoggingFile]; entry.Value != "~/.local/state/aibox/logs/run.log" || entry.Redacted {
		t.Fatalf("%+v", entry)
	}
	theme := byKey[settings.WorkspaceTheme]
	if theme.Value != "nord" || theme.Redacted || theme.Source != "~/[redacted].env" {
		t.Fatalf("%+v", theme)
	}
	last := data.SettingsSources[len(data.SettingsSources)-1]
	if last.Layer != "env_file" || last.Path != "~/[redacted].env" || data.SettingsSources[1].Path != "~/.config/aibox/settings.json" {
		t.Fatalf("%+v", data.SettingsSources)
	}
	if strings.Contains(string(encoded), canary) || strings.Count(string(encoded), f.home) != 0 {
		t.Fatalf("view leaks: %s", encoded)
	}
	f.flags = append(f.flags, settings.Flag{Name: "--timeout", Value: "5"})
	entries, _ := EffectiveConfiguration(f.resolve(t), redactor, settings.ReadClassLimitSeconds)
	for _, entry := range entries {
		if entry.Key == settings.ExecutionTimeout && (entry.Value != 5 || entry.Layer != "invocation") {
			t.Fatalf("%+v", entry)
		}
	}
}

// blockingInspector is an Inspector fake that records its deadline and
// waits until its context ends.
type blockingInspector struct {
	deadline chan time.Time
}

func (b blockingInspector) inspect(ctx context.Context, _ string) (string, string, contract.WorkspaceData, error) {
	deadline, _ := ctx.Deadline()
	b.deadline <- deadline
	<-ctx.Done()
	return "", "", contract.WorkspaceData{}, ctx.Err()
}

func TestReadClassLimitAppliesWhenUnset(t *testing.T) {
	// EXE-001: without execution.timeoutSeconds, inspect gets 30 seconds.
	f := newFixture(t)
	inspector := blockingInspector{deadline: make(chan time.Time, 1)}
	ctx, cancel := context.WithCancel(t.Context())
	resolved := f.resolve(t)
	started := time.Now()
	done := make(chan contract.Result, 1)
	go func() {
		done <- Workspace{Inspect: inspector.inspect, Health: sinkState{}}.InspectWorkspace(ctx, resolved, f.request(false), &recordedEvents{})
	}()
	deadline := <-inspector.deadline
	cancel()
	<-done
	if limit := deadline.Sub(started); limit < 29*time.Second || limit > 31*time.Second {
		t.Fatalf("limit %v", limit)
	}
}

func TestTimeoutOutcome(t *testing.T) {
	// EXE-002: past its limit an operation is timed_out with exit 124.
	f := newFixture(t)
	f.flags = []settings.Flag{{Name: "--timeout", Value: "1"}}
	events := &recordedEvents{}
	inspector := blockingInspector{deadline: make(chan time.Time, 1)}
	started := time.Now()
	result := Workspace{Inspect: inspector.inspect, Health: sinkState{}}.InspectWorkspace(t.Context(), f.resolve(t), f.request(false), events)
	if elapsed := time.Since(started); elapsed < time.Second || elapsed > 3*time.Second {
		t.Fatalf("took %v", elapsed)
	}
	if result.Outcome != contract.TimedOut || result.Error.Code != "operation_timeout" || result.Error.Category != "timeout" ||
		result.ExitCode() != 124 || result.ValidateEnvelope() != nil {
		t.Fatalf("%+v", result)
	}
	if events.names() != "operation.started,operation.finished" || events.events[1].Outcome != "timed_out" || events.events[1].ErrorCategory != "timeout" {
		t.Fatalf("%s %+v", events.names(), events.events)
	}
}

func TestCancellationOutcome(t *testing.T) {
	// LOG-008 at use-case level; the CLI test sends a real SIGINT.
	f := newFixture(t)
	events := &recordedEvents{}
	inspector := blockingInspector{deadline: make(chan time.Time, 1)}
	ctx, cancel := context.WithCancel(t.Context())
	go func() {
		<-inspector.deadline
		cancel()
	}()
	result := Workspace{Inspect: inspector.inspect, Health: sinkState{}}.InspectWorkspace(ctx, f.resolve(t), f.request(false), events)
	if result.Outcome != contract.Cancelled || result.Error.Code != "operation_cancelled" || result.ExitCode() != 130 {
		t.Fatalf("%+v", result)
	}
	if events.names() != "operation.started,operation.cancelled,operation.finished" || events.events[1].Severity != logging.Warn ||
		events.events[1].ErrorCategory != "interrupted" {
		t.Fatalf("%s", events.names())
	}
}

func TestRequiredSinkRefusesMutatingOperationBeforeAnyEffect(t *testing.T) {
	// LOG-031, mutating part (D16): a mutating operation behind the use-case
	// port is refused with operation_record_unavailable and its body never
	// runs; a read-only operation continues with log_sink_degraded.
	effects := 0
	body := func(context.Context) contract.Result {
		effects++
		result := contract.NewResult(contract.RefreshWorkspace, "req", "cli")
		result.Outcome = contract.NoChange
		return result
	}
	operation := Operation{Name: contract.RefreshWorkspace, Class: Mutating, TimeLimit: time.Minute, RequestID: "req", Initiator: "cli", LoggingRequired: true}
	result := Execute(t.Context(), operation, sinkState{failed: true}, body)
	if effects != 0 || result.Outcome != contract.Refused || result.Error.Code != "operation_record_unavailable" ||
		result.Error.Category != "internal" || result.ExitCode() != 5 || result.ValidateEnvelope() != nil {
		t.Fatalf("effects %d, %+v", effects, result)
	}
	operation.LoggingRequired = false
	if Execute(t.Context(), operation, sinkState{failed: true}, body); effects != 1 {
		t.Fatal("not required: the operation must run")
	}

	f := newFixture(t)
	if err := os.MkdirAll(filepath.Join(f.home, ".config", "aibox"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f.home, ".config", "aibox", "settings.json"), []byte(`{"schemaVersion":"1","logging":{"required":true}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	read := Workspace{Health: sinkState{failed: true}}.InspectWorkspace(t.Context(), f.resolve(t), f.request(false), &recordedEvents{})
	if read.Outcome != contract.Succeeded || len(read.Warnings) != 1 || read.Warnings[0] != contract.WarningLogSinkDegraded {
		t.Fatalf("%+v", read)
	}
}
