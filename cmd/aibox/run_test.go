package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/projectious-work/aibox/internal/contract"
	"github.com/projectious-work/aibox/internal/mcpserver"
	"github.com/projectious-work/aibox/internal/schematest"
	"github.com/projectious-work/aibox/internal/settings"
	"github.com/projectious-work/aibox/internal/usecase"
)

// harness runs the composition root in process with a fresh HOME, working
// directory and project, and buffers for stdout and stderr.
type harness struct {
	t                   *testing.T
	home, work, project string
	environ             []string
	inspect             usecase.Inspector
	stdout, stderr      bytes.Buffer
	stdoutTerminal      bool
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	base := t.TempDir()
	h := &harness{t: t, home: filepath.Join(base, "home"), work: filepath.Join(base, "work"), project: filepath.Join(base, "p")}
	for _, dir := range []string{h.home, h.work, filepath.Join(h.project, ".devcontainer")} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	h.write(filepath.Join(h.project, ".devcontainer", "devcontainer.json"), `{"image":"example.invalid/base:one"}`)
	h.environ = []string{"HOME=" + h.home, "PATH=/usr/bin:/bin"}
	return h
}

func (h *harness) write(path, contents string) {
	h.t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		h.t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		h.t.Fatal(err)
	}
}

func (h *harness) userSettings(contents string) {
	h.write(filepath.Join(h.home, ".config", "aibox", "settings.json"), contents)
}

func (h *harness) run(args ...string) int {
	h.stdout.Reset()
	h.stderr.Reset()
	return run(host{args: args, environ: h.environ, workDir: h.work, stdout: &h.stdout, stderr: &h.stderr,
		stdoutTerminal: h.stdoutTerminal, files: settings.OSFileSystem{}, inspect: h.inspect, now: time.Now,
		goos: "linux", serveMCP: mcpserver.Serve})
}

func (h *harness) view(extra ...string) []string {
	return append([]string{"inspect", "--context", "local", "--project", h.project, "--effective-config", "--format", "json", "--log-format", "jsonl"}, extra...)
}

func (h *harness) events() []map[string]any {
	h.t.Helper()
	var events []map[string]any
	for _, line := range strings.Split(strings.TrimSuffix(h.stderr.String(), "\n"), "\n") {
		if !strings.HasPrefix(line, "{") {
			continue
		}
		var event map[string]any
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			h.t.Fatalf("event %q: %v", line, err)
		}
		events = append(events, event)
	}
	return events
}

func (h *harness) result() map[string]any {
	h.t.Helper()
	var result map[string]any
	if err := json.Unmarshal(h.stdout.Bytes(), &result); err != nil {
		h.t.Fatalf("stdout %q: %v", h.stdout.String(), err)
	}
	return result
}

func validate(t *testing.T, schemaName string, document []byte) {
	t.Helper()
	schema, err := schematest.Load(schemaName)
	if err != nil {
		t.Fatal(err)
	}
	if err := schema.Validate(document); err != nil {
		t.Fatalf("%v: %s", err, document)
	}
}

var bootstrapPattern = regexp.MustCompile(`^aibox: invalid_input: ([a-z_]+): .+$`)

func (h *harness) expectRejected(reason string) map[string]any {
	h.t.Helper()
	events := h.events()
	if len(events) != 1 {
		h.t.Fatalf("expected one settings.rejected event, stderr %q", h.stderr.String())
	}
	event := events[0]
	envelope := h.result()
	validate(h.t, "schemas/operation-result.schema.json", h.stdout.Bytes())
	validate(h.t, "schemas/log-event.schema.json", bytes.TrimSpace(h.stderr.Bytes()))
	if event["event"] != "settings.rejected" || event["fields"].(map[string]any)["reason"] != reason || event["sequence"] != float64(0) ||
		event["requestId"] != envelope["requestId"] || envelope["operation"] != "invalid_request" || envelope["outcome"] != "failed" ||
		envelope["error"].(map[string]any)["code"] != "invalid_input" || envelope["error"].(map[string]any)["nextAction"] != "edit_configuration" {
		h.t.Fatalf("event %v envelope %v", event, envelope)
	}
	if _, present := envelope["target"]; present {
		h.t.Fatal("rejection envelope has a target")
	}
	return event
}

func TestRejectionWithEnvelope(t *testing.T) {
	// OUT-006, OUT-008, LOG-006, CFG-006 (J-OUT-04).
	h := newHarness(t)
	h.userSettings("not json")
	if code := h.run("inspect", "--project", h.project, "--format", "json", "--log-format", "jsonl"); code != 2 {
		t.Fatalf("exit %d", code)
	}
	event := h.expectRejected("malformed")
	if event["fields"].(map[string]any)["source"] != "~/.config/aibox/settings.json" || event["fields"].(map[string]any)["layer"] != "user" {
		t.Fatalf("%v", event)
	}
	if code := h.run("inspect", "--project", h.project, "--format", "yaml", "--log-format", "jsonl"); code != 2 ||
		!strings.Contains(h.stdout.String(), `"operation": "invalid_request"`) {
		t.Fatalf("exit %d stdout %q", code, h.stdout.String())
	}
	h = newHarness(t)
	h.run("inspect", "--project", h.project, "--format", "json", "--log-format", "jsonl", "--log-level", "loud")
	h.expectRejected("invalid_value")
	h.run("inspect", "--project", h.project, "--format", "json", "--log-format", "jsonl", "--format", "json")
	h.expectRejected("usage_error")
	for _, extra := range [][]string{{"--set", "a=b"}, {"extra"}, {"--effective-config=true"}, {"--log-level"}} {
		h.run(append([]string{"inspect", "--project", h.project, "--format", "json", "--log-format", "jsonl"}, extra...)...)
		h.expectRejected("usage_error")
	}
	var many []string
	for range 17 {
		many = append(many, "--env-file", "a.env")
	}
	h.run(append([]string{"inspect", "--project", h.project, "--format", "json", "--log-format", "jsonl"}, many...)...)
	h.expectRejected("usage_error")
}

func TestRejectionWithBootstrapDiagnostic(t *testing.T) {
	// OUT-007, LOG-006 rejection format (J-OUT-05).
	h := newHarness(t)
	h.userSettings("not json")
	for _, environ := range [][]string{nil, {"AIBOX_OUTPUT__FORMAT=json"}} {
		h.environ = append([]string{"HOME=" + h.home}, environ...)
		if code := h.run("inspect", "--project", h.project); code != 2 || h.stdout.Len() != 0 {
			t.Fatalf("exit %d stdout %q", code, h.stdout.String())
		}
		lines := strings.Split(strings.TrimSuffix(h.stderr.String(), "\n"), "\n")
		if len(lines) != 2 || !strings.Contains(lines[0], " ERROR settings.rejected ") ||
			bootstrapPattern.FindStringSubmatch(lines[1])[1] != "malformed" {
			t.Fatalf("%q", lines)
		}
	}
	h = newHarness(t)
	if code := h.run("inspect", "--format", "xml", "--log-level", "off"); code != 2 {
		t.Fatalf("exit %d", code)
	}
	if lines := strings.Split(strings.TrimSuffix(h.stderr.String(), "\n"), "\n"); len(lines) != 1 || bootstrapPattern.FindStringSubmatch(lines[0])[1] != "invalid_value" {
		t.Fatalf("%q", h.stderr.String())
	}
	if code := h.run("build"); code != 2 || !strings.Contains(h.stderr.String(), "not available") {
		t.Fatalf("exit %d %q", code, h.stderr.String())
	}
}

func TestVersionAndHelpIgnoreSettings(t *testing.T) {
	// CFG-009, D20 (J-CFG-23 (a), J-CFG-15 (g)).
	h := newHarness(t)
	h.userSettings("not json")
	h.environ = append(h.environ, "AIBOX_NOPE=1")
	for _, command := range []string{"version", "help"} {
		if code := h.run(command); code != 0 || h.stdout.Len() == 0 || h.stderr.Len() != 0 {
			t.Fatalf("%s: exit %d stderr %q", command, code, h.stderr.String())
		}
	}
	if code := h.run("version", "--format", "json"); code != 2 || h.stdout.Len() != 0 ||
		strings.Count(h.stderr.String(), "\n") != 1 || bootstrapPattern.FindStringSubmatch(strings.TrimSpace(h.stderr.String()))[1] != "usage_error" {
		t.Fatalf("exit %d stderr %q", code, h.stderr.String())
	}
}

func TestNoHome(t *testing.T) {
	// CFG-024 (J-CFG-23 (b)).
	h := newHarness(t)
	h.environ = []string{"PATH=/usr/bin"}
	h.run("inspect", "--project", h.project, "--format", "json", "--log-format", "jsonl")
	h.expectRejected("no_home")
}

func TestInspectEventsAndView(t *testing.T) {
	// LOG-001 to LOG-003, LOG-007, CFG-060, CFG-061, CFG-063, CFG-068 (J-LOG-01, J-CFG-01).
	h := newHarness(t)
	if code := h.run("inspect", "--project", h.project, "--format", "json", "--log-level", "trace", "--log-format", "jsonl"); code != 0 {
		t.Fatalf("exit %d %s", code, h.stderr.String())
	}
	result := h.result()
	validate(t, "schemas/operation-result.schema.json", h.stdout.Bytes())
	var names []string
	for index, event := range h.events() {
		names = append(names, event["event"].(string))
		if event["requestId"] != result["requestId"] || event["sequence"] != float64(index) {
			t.Fatalf("%v", event)
		}
	}
	if strings.Join(names, ",") != "settings.resolved,operation.started,operation.validated,operation.finished" {
		t.Fatalf("%v", names)
	}
	if strings.Contains(h.stdout.String(), "operation.started") {
		t.Fatal("event on stdout")
	}
	if events := h.events(); events[0]["fields"].(map[string]any)["layers"] != "invocation" ||
		events[1]["fields"].(map[string]any)["projectRoot"] != h.project {
		t.Fatalf("%v", events)
	}

	h.run(h.view()...)
	validate(t, "schemas/operation-result.schema.json", h.stdout.Bytes())
	data := h.result()["data"].(map[string]any)
	configuration := data["configuration"].([]any)
	if len(configuration) != 16 {
		t.Fatalf("%d entries", len(configuration))
	}
	sources := data["settingsSources"].([]any)
	want := []any{
		map[string]any{"layer": "system", "path": "/etc/aibox/settings.json", "state": "absent"},
		map[string]any{"layer": "user", "path": "~/.config/aibox/settings.json", "state": "absent"},
		map[string]any{"layer": "project", "path": filepath.Join(h.project, ".devcontainer", "devcontainer.json"), "state": "loaded"},
	}
	if !jsonEqual(sources, want) {
		t.Fatalf("%v", sources)
	}
	for _, entry := range configuration[1:] {
		fields := entry.(map[string]any)
		switch fields["key"] {
		case "output.format", "logging.format":
			if fields["layer"] != "invocation" || fields["invocationOnly"] != true {
				t.Fatalf("%v", fields)
			}
		default:
			if fields["layer"] != "built_in" || fields["source"] != "built-in" || fields["invocationOnly"] != false || fields["redacted"] != false {
				t.Fatalf("%v", fields)
			}
		}
	}
}

func jsonEqual(a, b any) bool {
	left, _ := json.Marshal(a)
	right, _ := json.Marshal(b)
	return bytes.Equal(left, right)
}

func TestDefaultHumanOutputAndColour(t *testing.T) {
	// OUT-004, OUT-005, D6: human is the default; colour by output.colorRule.
	h := newHarness(t)
	if code := h.run("inspect", "--project", h.project, "--log-level", "off"); code != 0 ||
		!strings.Contains(h.stdout.String(), "inspect_workspace") || !strings.Contains(h.stdout.String(), "succeeded") ||
		strings.Contains(h.stdout.String(), "\x1b") {
		t.Fatalf("%q", h.stdout.String())
	}
	h.stdoutTerminal = true
	if h.run("inspect", "--project", h.project, "--log-level", "off"); !strings.Contains(h.stdout.String(), "\x1b") {
		t.Fatal("no colour on a terminal")
	}
	h.environ = append(h.environ, "NO_COLOR=1")
	if h.run("inspect", "--project", h.project, "--log-level", "off", "--color", "always"); strings.Contains(h.stdout.String(), "\x1b") {
		t.Fatal("colour despite NO_COLOR")
	}
	h = newHarness(t)
	code := h.run("inspect", "--project", filepath.Join(h.project, "missing"), "--log-level", "off", "--color", "always")
	if code != 2 || !strings.Contains(h.stdout.String(), "failed") || !strings.Contains(h.stdout.String(), "invalid_input") ||
		!strings.Contains(h.stdout.String(), "\x1b") {
		t.Fatalf("exit %d %q", code, h.stdout.String())
	}
}

func TestPlainEventsByDefault(t *testing.T) {
	// LOG-010: the CLI default logging.format is plain.
	h := newHarness(t)
	h.run("inspect", "--project", h.project, "--format", "json", "--log-level", "trace")
	pattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?Z (TRACE|DEBUG|INFO|WARN|ERROR) [a-z][a-z0-9_.]*( [A-Za-z][A-Za-z0-9]*=.*)?$`)
	lines := strings.Split(strings.TrimSuffix(h.stderr.String(), "\n"), "\n")
	if len(lines) != 4 {
		t.Fatalf("%q", lines)
	}
	for _, line := range lines {
		if !pattern.MatchString(line) {
			t.Fatalf("%q", line)
		}
	}
}

func TestLogFileAndDegradedSink(t *testing.T) {
	// LOG-012, LOG-013, LOG-015, LOG-020, LOG-030 (J-LOG-05, J-LOG-10 (c)).
	h := newHarness(t)
	if code := h.run("inspect", "--project", h.project, "--format", "json", "--log-format", "jsonl", "--log-file", "run.log"); code != 0 {
		t.Fatalf("exit %d", code)
	}
	logs := filepath.Join(h.home, ".local", "state", "aibox", "logs")
	contents, err := os.ReadFile(filepath.Join(logs, "run.log"))
	if err != nil || string(contents) != h.stderr.String() {
		t.Fatalf("file %q stderr %q %v", contents, h.stderr.String(), err)
	}
	if _, err := os.Stat(filepath.Join(h.work, "run.log")); err == nil {
		t.Fatal("log file in the working directory")
	}
	wide := filepath.Join(logs, "wide.log")
	h.write(wide, "keep\n")
	if err := os.Chmod(wide, 0o644); err != nil {
		t.Fatal(err)
	}
	if code := h.run("inspect", "--project", h.project, "--format", "json", "--log-format", "jsonl", "--log-file", "wide.log"); code != 0 {
		t.Fatalf("exit %d", code)
	}
	result := h.result()
	if !slices.Contains(result["warnings"].([]any), any("log_sink_degraded")) {
		t.Fatalf("%v", result["warnings"])
	}
	degraded := 0
	for _, event := range h.events() {
		if event["event"] == "logging.sink_degraded" {
			degraded++
			if event["fields"].(map[string]any)["reason"] != "insecure_mode" {
				t.Fatalf("%v", event)
			}
		}
	}
	if contents, _ := os.ReadFile(wide); degraded != 1 || string(contents) != "keep\n" {
		t.Fatalf("degraded %d, wide.log %q", degraded, contents)
	}
	h.run("inspect", "--project", h.project, "--format", "json", "--log-format", "jsonl", "--log-file", "/tmp/x.log")
	h.expectRejected("path_escape")
}

func TestWorkingDirectoryAnchorsSelectorsAndStaysUnchanged(t *testing.T) {
	// CFG-017, CFG-031, CFG-052, CFG-053: relative selectors are anchored to
	// the working directory at process start, which aibox never changes.
	h := newHarness(t)
	h.write(filepath.Join(h.work, "other.json"), `{"schemaVersion":"1","logging":{"level":"warn"}}`)
	h.write(filepath.Join(h.work, "a.env"), "AIBOX_WORKSPACE__THEME=nord\n")
	sub := filepath.Join(h.work, "sub", "p", ".devcontainer")
	h.write(filepath.Join(sub, "devcontainer.json"), `{"image":"x"}`)
	before, _ := os.Getwd()
	code := h.run("inspect", "--project", "sub/p", "--effective-config", "--format", "json", "--log-level", "off",
		"--settings-file", "other.json", "--env-file", "a.env")
	after, _ := os.Getwd()
	if code != 0 || before != after {
		t.Fatalf("exit %d, cwd %s -> %s", code, before, after)
	}
	result := h.result()
	if result["target"].(map[string]any)["workspaceRoot"] != filepath.Join(h.work, "sub", "p") {
		t.Fatalf("%v", result["target"])
	}
	sources := result["data"].(map[string]any)["settingsSources"].([]any)
	if sources[1].(map[string]any)["path"] != filepath.Join(h.work, "other.json") || sources[3].(map[string]any)["path"] != filepath.Join(h.work, "a.env") {
		t.Fatalf("%v", sources)
	}
}

func TestSignalCancelsInspect(t *testing.T) {
	// LOG-008: SIGINT during inspect gives operation.cancelled, outcome
	// cancelled with operation_cancelled, and exit 130. The inspection port
	// blocks so that the signal arrives while the operation runs.
	h := newHarness(t)
	entered := make(chan struct{})
	h.inspect = func(ctx context.Context, _ string) (string, string, contract.WorkspaceData, error) {
		close(entered)
		<-ctx.Done()
		return "", "", contract.WorkspaceData{}, ctx.Err()
	}
	go func() {
		<-entered
		_ = syscall.Kill(os.Getpid(), syscall.SIGINT)
	}()
	code := h.run("inspect", "--project", h.project, "--format", "json", "--log-format", "jsonl")
	result := h.result()
	if code != 130 || result["outcome"] != "cancelled" || result["error"].(map[string]any)["code"] != "operation_cancelled" {
		t.Fatalf("exit %d %v", code, result)
	}
	validate(t, "schemas/operation-result.schema.json", h.stdout.Bytes())
	var names []string
	for _, event := range h.events() {
		names = append(names, event["event"].(string))
	}
	if strings.Join(names, ",") != "operation.started,operation.cancelled,operation.finished" {
		t.Fatalf("%v", names)
	}
}

func TestMCPServeRejectsPlainEvents(t *testing.T) {
	// CFG-054, CFG-057 (J-CFG-18 (b)).
	h := newHarness(t)
	if code := h.run("mcp", "serve", "--log-format", "plain"); code != 2 || h.stdout.Len() != 0 {
		t.Fatalf("exit %d", code)
	}
	lines := strings.Split(strings.TrimSuffix(h.stderr.String(), "\n"), "\n")
	if len(lines) != 2 || !strings.Contains(lines[0], `"reason":"invalid_value"`) || bootstrapPattern.FindStringSubmatch(lines[1])[1] != "invalid_value" {
		t.Fatalf("%q", lines)
	}
}

func TestFlagTableMatchesTheContract(t *testing.T) {
	contents, err := os.ReadFile(schematest.ContractPath("interfaces.json"))
	if err != nil {
		t.Fatal(err)
	}
	var contract struct {
		Settings struct {
			Selectors []struct {
				Name       string   `json:"name"`
				Repeatable bool     `json:"repeatable"`
				Commands   []string `json:"commands"`
			} `json:"selectors"`
		} `json:"settings"`
		CLI struct {
			Flags []struct {
				Flag     string   `json:"flag"`
				Key      *string  `json:"key"`
				Commands []string `json:"commands"`
			} `json:"flags"`
		} `json:"cli"`
	}
	if err := json.Unmarshal(contents, &contract); err != nil {
		t.Fatal(err)
	}
	type row struct {
		name       string
		repeatable bool
		setting    bool
		commands   string
	}
	var want []row
	for _, flag := range contract.CLI.Flags {
		want = append(want, row{flag.Flag, false, flag.Key != nil, strings.Join(flag.Commands, ",")})
	}
	for _, selector := range contract.Settings.Selectors {
		if selector.Commands != nil {
			want = append(want, row{selector.Name, selector.Repeatable, false, strings.Join(selector.Commands, ",")})
		}
	}
	var got []row
	for _, spec := range flagTable {
		got = append(got, row{spec.name, spec.repeatable, spec.setting, strings.Join(spec.commands, ",")})
	}
	if !slices.Equal(got, want) {
		t.Fatalf("flag table drifted:\ngot  %v\nwant %v", got, want)
	}
}
