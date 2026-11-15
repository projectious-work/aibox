package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func minimalProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".devcontainer"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".devcontainer", "devcontainer.json"), []byte(`{"image":"example.invalid/base:one"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestMCPViewParityAndEvents(t *testing.T) {
	// CFG-007, CFG-055, CFG-065, LOG-004, LOG-011, LOG-012 (J-CFG-17, J-LOG-04).
	binary := buildPreview(t)
	home, project := t.TempDir(), minimalProject(t)
	settingsFile := filepath.Join(home, ".config", "aibox", "settings.json")
	if err := os.MkdirAll(filepath.Dir(settingsFile), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(settingsFile, []byte(`{"schemaVersion":"1","logging":{"level":"warn"},"output":{"color":"never"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	environ := []string{"HOME=" + home, "PATH=" + os.Getenv("PATH")}
	server := exec.Command(binary, "mcp", "serve", "--context", "local", "--log-level", "debug", "--log-file", "m.log")
	server.Env = environ
	serverStderr := &lockedBuffer{}
	server.Stderr = serverStderr
	client := mcp.NewClient(&mcp.Implementation{Name: "aibox-v1-06-test", Version: "1"}, nil)
	session, err := client.Connect(t.Context(), &mcp.CommandTransport{Command: server}, nil)
	if err != nil {
		t.Fatal(err)
	}
	var mcpData map[string]any
	for index, id := range []string{"call-1", "call-2"} {
		response, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "inspect_workspace",
			Arguments: map[string]any{"requestId": id, "projectRoot": project, "effectiveConfig": index == 0}})
		if err != nil || response.IsError {
			t.Fatalf("%v %+v", err, response)
		}
		if index == 0 {
			encoded, _ := json.Marshal(response.StructuredContent)
			validate(t, "schemas/operation-result.schema.json", encoded)
			var result map[string]any
			_ = json.Unmarshal(encoded, &result)
			mcpData = result["data"].(map[string]any)
		}
	}
	for _, extra := range []string{"logLevel", "logging", "outputFormat"} {
		response, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "inspect_workspace",
			Arguments: map[string]any{"requestId": "bad", "projectRoot": project, extra: "x"}})
		if err == nil && !response.IsError {
			t.Fatalf("input member %s accepted", extra)
		}
	}
	_ = session.Close()
	deadline := time.Now().Add(5 * time.Second)
	for !strings.Contains(serverStderr.String(), "process.finished") && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}

	// Events: process.started, settings.resolved first with the process ID;
	// each call has its own request ID with sequence 0 to 2.
	lines := strings.Split(strings.TrimSuffix(serverStderr.String(), "\n"), "\n")
	sequences := map[string][]float64{}
	var first, last map[string]any
	for index, line := range lines {
		validate(t, "schemas/log-event.schema.json", []byte(line))
		var event map[string]any
		_ = json.Unmarshal([]byte(line), &event)
		id := event["requestId"].(string)
		sequences[id] = append(sequences[id], event["sequence"].(float64))
		if index == 0 {
			first = event
		}
		last = event
		if id != first["requestId"] && event["initiator"] != "mcp" {
			t.Fatalf("%v", event)
		}
	}
	processID := first["requestId"].(string)
	if first["event"] != "process.started" || last["event"] != "process.finished" || last["requestId"] != processID ||
		len(sequences) != 3 || jsonText(sequences[processID]) != "[0,1,2]" ||
		jsonText(sequences["call-1"]) != "[0,1,2]" || jsonText(sequences["call-2"]) != "[0,1,2]" {
		t.Fatalf("%v\n%s", sequences, serverStderr.String())
	}
	logged, err := os.ReadFile(filepath.Join(home, ".local", "state", "aibox", "logs", "m.log"))
	if err != nil || string(logged) != serverStderr.String() {
		t.Fatalf("m.log %q %v", logged, err)
	}

	// CLI view with the same layers: equal settingsSources and entries except
	// the two formats that the view itself passes as flags.
	cli := exec.Command(binary, "inspect", "--context", "local", "--project", project, "--effective-config", "--format", "json", "--log-format", "jsonl")
	cli.Env = environ
	stdout, err := cli.Output()
	if err != nil {
		t.Fatal(err)
	}
	var cliResult map[string]any
	_ = json.Unmarshal(stdout, &cliResult)
	cliData := cliResult["data"].(map[string]any)
	if jsonText(cliData["settingsSources"]) != jsonText(mcpData["settingsSources"]) {
		t.Fatalf("sources differ:\n%v\n%v", cliData["settingsSources"], mcpData["settingsSources"])
	}
	cliEntries, mcpEntries := cliData["configuration"].([]any), mcpData["configuration"].([]any)
	for index, entry := range cliEntries {
		cliEntry, mcpEntry := entry.(map[string]any), mcpEntries[index].(map[string]any)
		switch cliEntry["key"] {
		case "output.format", "logging.format":
			if mcpEntry["layer"] != "built_in" {
				t.Fatalf("%v", mcpEntry)
			}
		case "logging.level":
			if mcpEntry["value"] != "debug" || mcpEntry["layer"] != "invocation" || cliEntry["value"] != "warn" {
				t.Fatalf("cli %v mcp %v", cliEntry, mcpEntry)
			}
		case "logging.file":
		default:
			if jsonText(cliEntry) != jsonText(mcpEntry) {
				t.Fatalf("cli %v mcp %v", cliEntry, mcpEntry)
			}
		}
	}
}

// lockedBuffer collects a child's stderr while the test reads it.
type lockedBuffer struct {
	mu     sync.Mutex
	buffer bytes.Buffer
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buffer.Write(p)
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buffer.String()
}

func jsonText(value any) string {
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

func TestClosedStderrKeepsTheResult(t *testing.T) {
	// LOG-032 (J-LOG-11 (a)): stderr closed; the result and the log file are complete.
	binary := buildPreview(t)
	home, project := t.TempDir(), minimalProject(t)
	shell := exec.Command("sh", "-c", `"$0" inspect --project "$1" --format json --log-level trace --log-file e.log 2>&-`, binary, project)
	shell.Env = []string{"HOME=" + home, "PATH=" + os.Getenv("PATH")}
	stdout, err := shell.Output()
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(stdout, &result); err != nil || result["outcome"] != "succeeded" || len(result["warnings"].([]any)) != 0 {
		t.Fatalf("%s %v", stdout, err)
	}
	logged, _ := os.ReadFile(filepath.Join(home, ".local", "state", "aibox", "logs", "e.log"))
	if lines := strings.Split(strings.TrimSuffix(string(logged), "\n"), "\n"); len(lines) != 4 {
		t.Fatalf("e.log %q", logged)
	}
}

func TestStaticBinaryWithoutUserDatabase(t *testing.T) {
	// CFG-024: HOME only. The binary needs no cgo and no os/user.
	list := exec.Command("go", "list", "-deps", ".")
	list.Env = append(os.Environ(), "CGO_ENABLED=0")
	output, err := list.Output()
	if err != nil {
		t.Fatal(err)
	}
	for _, dependency := range strings.Split(string(output), "\n") {
		if dependency == "os/user" || dependency == "runtime/cgo" {
			t.Fatalf("the binary depends on %s", dependency)
		}
	}
}
