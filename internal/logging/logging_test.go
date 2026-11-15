package logging

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/projectious-work/aibox/internal/redact"
	"github.com/projectious-work/aibox/internal/schematest"
)

// syncBuffer is a goroutine-safe stderr fake.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) lines() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	text := strings.TrimSuffix(b.buf.String(), "\n")
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}

// recordingFile is a FileAppender fake that records lines, fails on demand,
// or blocks until released.
type recordingFile struct {
	mu      sync.Mutex
	lines   []string
	fail    error
	release chan struct{}
	entered chan struct{}
}

func (f *recordingFile) Append(line []byte) error {
	if f.release != nil {
		select {
		case f.entered <- struct{}{}:
		default:
		}
		<-f.release
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.fail != nil {
		return f.fail
	}
	f.lines = append(f.lines, strings.TrimSuffix(string(line), "\n"))
	return nil
}

var fixedTime = time.Date(2026, 11, 15, 10, 0, 0, 123000000, time.UTC)

func fixedClock() time.Time { return fixedTime }

func decode(t *testing.T, line string) map[string]any {
	t.Helper()
	var event map[string]any
	if err := json.Unmarshal([]byte(line), &event); err != nil {
		t.Fatalf("not JSON: %q", line)
	}
	return event
}

func validateEvents(t *testing.T, lines []string) {
	t.Helper()
	schema, err := schematest.Load("schemas/log-event.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range lines {
		if err := schema.Validate([]byte(line)); err != nil {
			t.Fatalf("%v: %s", err, line)
		}
	}
}

func TestEventsValidateAndCarryRequestSequence(t *testing.T) {
	// LOG-002, LOG-003, LOG-010 (jsonl), LOG-012.
	stderr, file := &syncBuffer{}, &recordingFile{}
	logger := Start(Config{Level: "trace", Format: JSONL, Stderr: stderr, File: file, Now: fixedClock})
	scope := logger.Scope("req-1")
	scope.Emit(Event{Severity: Debug, Name: "settings.resolved", Component: "settings", Fields: map[string]string{"layers": ""}})
	scope.Emit(Event{Severity: Info, Name: "operation.started", Component: "usecase", Operation: "inspect_workspace", Initiator: "cli",
		Fields: map[string]string{"projectRoot": "/p"}})
	scope.Emit(Event{Severity: Info, Name: "operation.finished", Component: "usecase", Operation: "inspect_workspace", Initiator: "cli", Outcome: "succeeded"})
	logger.Close(scope)
	lines := stderr.lines()
	validateEvents(t, lines)
	if len(lines) != 3 || strings.Join(file.lines, "\n") != strings.Join(lines, "\n") {
		t.Fatalf("stderr %v file %v", lines, file.lines)
	}
	for index, line := range lines {
		event := decode(t, line)
		if event["requestId"] != "req-1" || event["sequence"] != float64(index) || event["executor"] != "aibox" ||
			event["schemaVersion"] != SchemaVersion || event["timestamp"] != "2026-11-15T10:00:00.123Z" {
			t.Fatalf("%v", event)
		}
	}
	if _, present := decode(t, lines[2])["fields"]; present {
		t.Fatal("empty fields encoded")
	}
}

func TestSchemaCheckerRejectsBadEvents(t *testing.T) {
	schema, err := schematest.Load("schemas/log-event.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{`{}`, `{"schemaVersion":"aibox.log-event/v1","timestamp":"x","severity":"info","event":"e","component":"c","requestId":"r","sequence":0}`,
		`{"schemaVersion":"aibox.log-event/v1","timestamp":"2026-01-01T00:00:00Z","severity":"info","event":"e","component":"c","requestId":"r","sequence":0,"fields":{"n":1}}`} {
		if schema.Validate([]byte(bad)) == nil {
			t.Fatalf("accepted %s", bad)
		}
	}
}

var plainLine = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?Z (TRACE|DEBUG|INFO|WARN|ERROR) [a-z][a-z0-9_.]*( [A-Za-z][A-Za-z0-9]*=.*)?$`)

func TestPlainLines(t *testing.T) {
	// LOG-010: logging.plainLine; a value with spaces is JSON-quoted.
	stderr, file := &syncBuffer{}, &recordingFile{}
	logger := Start(Config{Level: "trace", Format: Plain, Stderr: stderr, File: file, Now: fixedClock})
	scope := logger.Scope("req-2")
	scope.Emit(Event{Severity: Info, Name: "operation.started", Component: "usecase", Operation: "inspect_workspace",
		Initiator: "cli", Fields: map[string]string{"projectRoot": "/a b/c", "layers": ""}})
	logger.Close(scope)
	lines := stderr.lines()
	if len(lines) != 1 || !plainLine.MatchString(lines[0]) {
		t.Fatalf("%q", lines)
	}
	want := `2026-11-15T10:00:00.123Z INFO operation.started component=usecase operation=inspect_workspace requestId=req-2 sequence=0 initiator=cli executor=aibox layers="" projectRoot="/a b/c"`
	if lines[0] != want {
		t.Fatalf("got  %s\nwant %s", lines[0], want)
	}
	validateEvents(t, file.lines)
}

func TestLevelFilter(t *testing.T) {
	// LOG-005: at or above logging.level; off emits nothing.
	for level, want := range map[string]int{"trace": 5, "debug": 4, "info": 3, "warn": 2, "error": 1, "off": 0} {
		stderr := &syncBuffer{}
		logger := Start(Config{Level: level, Format: JSONL, Stderr: stderr, Now: fixedClock})
		scope := logger.Scope("r")
		for severity := Trace; severity <= Error; severity++ {
			scope.Emit(Event{Severity: severity, Name: "test.event", Component: "cli"})
		}
		logger.Close(scope)
		lines := stderr.lines()
		if len(lines) != want {
			t.Fatalf("%s: %d events", level, len(lines))
		}
		for index, line := range lines {
			if decode(t, line)["sequence"] != float64(index) {
				t.Fatalf("%s: sequence gap in %v", level, lines)
			}
		}
	}
}

func TestRedactionBeforeAnySink(t *testing.T) {
	// RED-001 to RED-003 for events: secrets in every member and field, the
	// home prefix in fields.
	home := t.TempDir()
	stderr, file := &syncBuffer{}, &recordingFile{}
	redactor := redact.New([]string{"EXAMPLE_API_TOKEN=canary-0123456789"}, home)
	logger := Start(Config{Level: "trace", Format: JSONL, Stderr: stderr, File: file, Redactor: redactor, Now: fixedClock})
	scope := logger.Scope("canary-0123456789")
	scope.Emit(Event{Severity: Trace, Name: "operation.started", Component: "usecase", Initiator: "ghp_" + strings.Repeat("a", 36),
		Fields: map[string]string{"projectRoot": home + "/p-canary-0123456789"}})
	logger.Close(scope)
	all := strings.Join(stderr.lines(), "\n") + strings.Join(file.lines, "\n")
	if strings.Contains(all, "canary") || strings.Contains(all, "ghp_") || strings.Contains(all, home) {
		t.Fatalf("leaked: %s", all)
	}
	if event := decode(t, stderr.lines()[0]); event["fields"].(map[string]any)["projectRoot"] != "~/p-[redacted]" {
		t.Fatalf("%v", event)
	}
}

func TestLongLinesLoseFieldsLastNameFirst(t *testing.T) {
	big := strings.Repeat("x", MaxLineBytes)
	event := Event{Timestamp: fixedTime, Severity: Info, Name: "e", Component: "c", RequestID: "r", Executor: "aibox",
		Fields: map[string]string{"a": "keep", "y": big, "z": big}}.fit()
	line := event.encodeJSON()
	if len(line) > MaxLineBytes || event.Fields["a"] != "keep" || event.Fields["truncated"] != "true" || event.Fields["y"] != "" {
		t.Fatalf("%d bytes, fields %v", len(line), len(event.Fields))
	}
}

func TestFileWriteFailureDegradesOnce(t *testing.T) {
	// LOG-030 (write failure): one logging.sink_degraded on stderr with the
	// next sequence of the request; stderr output continues.
	stderr, file := &syncBuffer{}, &recordingFile{fail: errors.New("disk full")}
	logger := Start(Config{Level: "info", Format: JSONL, Stderr: stderr, File: file, Now: fixedClock})
	scope := logger.Scope("r")
	scope.Emit(Event{Severity: Info, Name: "first", Component: "usecase"})
	waitFor(t, logger.FileSinkFailed)
	scope.Emit(Event{Severity: Info, Name: "second", Component: "usecase"})
	logger.Close(scope)
	var names []string
	for _, line := range stderr.lines() {
		event := decode(t, line)
		names = append(names, event["event"].(string))
		if event["event"] == "logging.sink_degraded" {
			fields := event["fields"].(map[string]any)
			if fields["sink"] != "file" || fields["reason"] != "write_failed" || event["sequence"] != float64(1) {
				t.Fatalf("%v", event)
			}
		}
	}
	if strings.Join(names, ",") != "first,logging.sink_degraded,second" {
		t.Fatalf("%v", names)
	}
}

func TestOpenFailureIsReportedOnce(t *testing.T) {
	stderr := &syncBuffer{}
	logger := Start(Config{Level: "warn", Format: JSONL, Stderr: stderr, FileFailure: InsecureMode, Now: fixedClock})
	scope := logger.Scope("r")
	logger.ReportOpenFailure(scope)
	logger.ReportOpenFailure(scope)
	logger.Close(scope)
	lines := stderr.lines()
	if len(lines) != 1 || !logger.FileSinkFailed() || decode(t, lines[0])["fields"].(map[string]any)["reason"] != "insecure_mode" {
		t.Fatalf("%v", lines)
	}
}

// failingWriter is a stderr whose reader has gone.
type failingWriter struct{ calls int }

func (w *failingWriter) Write([]byte) (int, error) { w.calls++; return 0, errors.New("broken pipe") }

func TestClosedStderrIsSilent(t *testing.T) {
	// LOG-032: the file still gets every event; no event about stderr.
	stderr, file := &failingWriter{}, &recordingFile{}
	logger := Start(Config{Level: "trace", Format: JSONL, Stderr: stderr, File: file, Now: fixedClock})
	scope := logger.Scope("r")
	for range 4 {
		scope.Emit(Event{Severity: Info, Name: "operation.started", Component: "usecase"})
	}
	logger.Close(scope)
	if len(file.lines) != 4 || stderr.calls != 1 || logger.FileSinkFailed() {
		t.Fatalf("file %d lines, stderr calls %d", len(file.lines), stderr.calls)
	}
	for _, line := range file.lines {
		if decode(t, line)["event"] != "operation.started" {
			t.Fatalf("unexpected event in file: %s", line)
		}
	}
}

func TestFullQueueDropsDebugAndWritesWarnDirectly(t *testing.T) {
	// LOG-033 with the contract's limits and a file sink that blocks on demand.
	stderr := &syncBuffer{}
	file := &recordingFile{release: make(chan struct{}), entered: make(chan struct{}, 1)}
	logger := Start(Config{Level: "trace", Format: JSONL, Stderr: stderr, File: file, Now: fixedClock})
	scope := logger.Scope("r")
	scope.Emit(Event{Severity: Info, Name: "blocker", Component: "usecase"})
	<-file.entered
	for range DefaultLimits.Capacity {
		scope.Emit(Event{Severity: Trace, Name: "fill", Component: "usecase"})
	}
	scope.Emit(Event{Severity: Debug, Name: "dropped.one", Component: "usecase"})
	scope.Emit(Event{Severity: Trace, Name: "dropped.two", Component: "usecase"})
	started := time.Now()
	scope.Emit(Event{Severity: Warn, Name: "urgent", Component: "usecase"})
	if elapsed := time.Since(started); elapsed > DefaultLimits.DirectWrite+50*time.Millisecond {
		t.Fatalf("warn took %v", elapsed)
	}
	if !strings.Contains(strings.Join(stderr.lines(), "\n"), `"event":"urgent"`) {
		t.Fatal("warn event not written directly to stderr")
	}
	started = time.Now()
	scope.Emit(Event{Severity: Info, Name: "dropped.info", Component: "usecase"})
	if elapsed := time.Since(started); elapsed > DefaultLimits.DirectWrite+50*time.Millisecond {
		t.Fatalf("info blocked for %v", elapsed)
	}
	close(file.release)
	waitFor(t, func() bool { return len(logger.queue) == 0 })
	scope.Emit(Event{Severity: Info, Name: "after", Component: "usecase"})
	logger.Close(scope)
	var dropped []map[string]any
	for _, line := range stderr.lines() {
		event := decode(t, line)
		if strings.HasPrefix(event["event"].(string), "dropped.") {
			t.Fatalf("dropped event written: %v", event)
		}
		if event["event"] == "logging.dropped" {
			dropped = append(dropped, event)
		}
	}
	if len(dropped) != 1 || dropped[0]["fields"].(map[string]any)["count"] != "3" {
		t.Fatalf("logging.dropped: %v", dropped)
	}
}

func TestShutdownFlushIsBounded(t *testing.T) {
	// LOG-034: Close waits at most flushSeconds and then reports
	// logging.flush_incomplete directly on stderr.
	stderr := &syncBuffer{}
	file := &recordingFile{release: make(chan struct{}), entered: make(chan struct{}, 1)}
	defer close(file.release)
	logger := Start(Config{Level: "info", Format: JSONL, Stderr: stderr, File: file, Now: fixedClock})
	scope := logger.Scope("r")
	scope.Emit(Event{Severity: Info, Name: "stuck", Component: "usecase"})
	<-file.entered
	started := time.Now()
	logger.Close(scope)
	elapsed := time.Since(started)
	if elapsed < DefaultLimits.Flush || elapsed > DefaultLimits.Flush+DefaultLimits.DirectWrite+200*time.Millisecond {
		t.Fatalf("Close took %v", elapsed)
	}
	lines := stderr.lines()
	last := decode(t, lines[len(lines)-1])
	if last["event"] != "logging.flush_incomplete" || last["fields"].(map[string]any)["sink"] != "file" || last["sequence"] != float64(1) {
		t.Fatalf("%v", lines)
	}
	validateEvents(t, lines)
}

func TestNormalShutdownReportsNoFlush(t *testing.T) {
	// D21: a normal run has a fixed event list.
	stderr := &syncBuffer{}
	logger := Start(Config{Level: "trace", Format: JSONL, Stderr: stderr, File: &recordingFile{}, Now: fixedClock})
	scope := logger.Scope("r")
	scope.Emit(Event{Severity: Info, Name: "only", Component: "usecase"})
	logger.Close(scope)
	logger.Close(scope)
	scope.Emit(Event{Severity: Info, Name: "late", Component: "usecase"})
	if lines := stderr.lines(); len(lines) != 1 {
		t.Fatalf("%v", lines)
	}
}

func TestRejectionEvent(t *testing.T) {
	// LOG-006: settings.rejected with sequence 0, written directly.
	var out bytes.Buffer
	home := "/home/someone"
	WriteRejection(&out, JSONL, "req", map[string]string{"reason": "malformed", "source": home + "/.config/aibox/settings.json"},
		redact.New(nil, home), fixedTime)
	event := decode(t, out.String())
	validateEvents(t, []string{strings.TrimSuffix(out.String(), "\n")})
	if event["event"] != "settings.rejected" || event["sequence"] != float64(0) || event["severity"] != "error" ||
		event["outcome"] != "failed" || event["errorCategory"] != "invalid_input" ||
		event["fields"].(map[string]any)["source"] != "~/.config/aibox/settings.json" {
		t.Fatalf("%v", event)
	}
	out.Reset()
	WriteRejection(&out, Plain, "req", map[string]string{"reason": "usage_error"}, nil, fixedTime)
	if !plainLine.MatchString(strings.TrimSuffix(out.String(), "\n")) {
		t.Fatalf("%q", out.String())
	}
}

func waitFor(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatal("condition not reached")
		}
		time.Sleep(time.Millisecond)
	}
}

func TestContractConstants(t *testing.T) {
	// Drift test for the queue limits, severities and line limits.
	contents, err := os.ReadFile(schematest.ContractPath("interfaces.json"))
	if err != nil {
		t.Fatal(err)
	}
	var contract struct {
		Logging struct {
			Severities []string `json:"severities"`
			Queue      struct {
				Capacity          int `json:"capacity"`
				FlushSeconds      int `json:"flushSeconds"`
				DirectWriteMillis int `json:"directWriteMillis"`
			} `json:"queue"`
			FileSink struct {
				MaxLineBytes int `json:"maxLineBytes"`
			} `json:"fileSink"`
			PlainLine string `json:"plainLine"`
		} `json:"logging"`
	}
	if err := json.Unmarshal(contents, &contract); err != nil {
		t.Fatal(err)
	}
	queue := contract.Logging.Queue
	if strings.Join(contract.Logging.Severities, ",") != strings.Join(severityNames, ",") || queue.Capacity != DefaultLimits.Capacity ||
		time.Duration(queue.FlushSeconds)*time.Second != DefaultLimits.Flush ||
		time.Duration(queue.DirectWriteMillis)*time.Millisecond != DefaultLimits.DirectWrite ||
		contract.Logging.FileSink.MaxLineBytes != MaxLineBytes || contract.Logging.PlainLine != plainLine.String() {
		t.Fatalf("logging constants drifted: %+v", contract.Logging)
	}
}
