package output

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/projectious-work/aibox/internal/contract"
	"github.com/projectious-work/aibox/internal/redact"
)

func sampleResult() contract.Result {
	result := contract.NewResult(contract.InspectWorkspace, "req-1", "cli")
	result.Target = &contract.Target{Scope: "local", WorkspaceRoot: "/w/p-canary-0123456789"}
	result.InputDigest = "sha256:" + strings.Repeat("a", 64)
	result.Outcome = contract.Succeeded
	result.Data = contract.WorkspaceData{
		State: "unknown", Resources: []any{}, ObservedAt: "2026-11-15T10:00:00Z",
		Configuration: []contract.ConfigurationEntry{
			{Key: "devcontainer.features", Value: []string{"a", "b"}, Layer: "project", Source: ".devcontainer/devcontainer.json", Domain: "native"},
			{Key: "logging.level", Value: "warn", Layer: "user", Source: "~/.config/aibox/settings.json", Domain: "process",
				Overridden: []contract.OverriddenValue{{Layer: "system", Source: "/etc/aibox/settings.json", Value: "debug"}}},
			{Key: "execution.timeoutSeconds", Value: 30, Layer: "built_in", Source: "built-in", Domain: "process"},
			{Key: "logging.file", Value: nil, Layer: "built_in", Source: "built-in", Domain: "process"},
			{Key: "odd", Value: "quote\" back\\ tab\t nl\n del\x7f nel\u0085 ls  é", Layer: "built_in", Source: "built-in", Domain: "process"},
		},
		SettingsSources: []contract.SettingsSource{{Layer: "system", Path: "/etc/aibox/settings.json", State: "absent"}},
	}
	return result
}

func render(t *testing.T, result contract.Result, options Options) string {
	t.Helper()
	var out bytes.Buffer
	if err := Write(&out, result, options); err != nil {
		t.Fatal(err)
	}
	return out.String()
}

func TestJSONIsOneDocumentAndANewline(t *testing.T) {
	// OUT-002.
	text := render(t, sampleResult(), Options{Format: JSON})
	if strings.Count(text, "\n") != 1 || !strings.HasSuffix(text, "}\n") || !strings.Contains(text, `"overridden":[{"layer":"system"`) {
		t.Fatalf("%s", text)
	}
}

func TestYAMLBlockStyle(t *testing.T) {
	// OUT-003 and output.yaml; the acceptance suite parses it with a YAML 1.2 parser.
	want := `"schemaVersion": "aibox.operation-result/v1"
"operation": "inspect_workspace"
"requestId": "req-1"
"actors":
  "initiator": "cli"
  "executor": "aibox"
"target":
  "scope": "local"
  "workspaceRoot": "/w/p-canary-0123456789"
"inputDigest": "sha256:` + strings.Repeat("a", 64) + `"
"outcome": "succeeded"
"changedResources": []
"warnings": []
"evidence": []
"data":
  "state": "unknown"
  "resources": []
  "configuration":
    - "key": "devcontainer.features"
      "value":
        - "a"
        - "b"
      "redacted": false
      "layer": "project"
      "source": ".devcontainer/devcontainer.json"
      "domain": "native"
      "invocationOnly": false
    - "key": "logging.level"
      "value": "warn"
      "redacted": false
      "layer": "user"
      "source": "~/.config/aibox/settings.json"
      "domain": "process"
      "invocationOnly": false
      "overridden":
        - "layer": "system"
          "source": "/etc/aibox/settings.json"
          "value": "debug"
          "redacted": false
    - "key": "execution.timeoutSeconds"
      "value": 30
      "redacted": false
      "layer": "built_in"
      "source": "built-in"
      "domain": "process"
      "invocationOnly": false
    - "key": "logging.file"
      "value": null
      "redacted": false
      "layer": "built_in"
      "source": "built-in"
      "domain": "process"
      "invocationOnly": false
    - "key": "odd"
      "value": "quote\" back\\ tab\t nl\n del\u007f nel\u0085 ls\u2028 é"
      "redacted": false
      "layer": "built_in"
      "source": "built-in"
      "domain": "process"
      "invocationOnly": false
  "observedAt": "2026-11-15T10:00:00Z"
  "settingsSources":
    - "layer": "system"
      "path": "/etc/aibox/settings.json"
      "state": "absent"
`
	if got := render(t, sampleResult(), Options{Format: YAML}); got != want {
		gotLines, wantLines := strings.Split(got, "\n"), strings.Split(want, "\n")
		for index := range min(len(gotLines), len(wantLines)) {
			if gotLines[index] != wantLines[index] {
				t.Fatalf("line %d:\ngot  %q\nwant %q", index+1, gotLines[index], wantLines[index])
			}
		}
		t.Fatalf("got %d lines, want %d", len(gotLines), len(wantLines))
	}
}

func TestYAMLNestedCollections(t *testing.T) {
	got, err := ToYAML([]byte(`{"a":[[1,[]],{"b":{}},[]],"c":{"d":{"e":true}},"f":{}}`))
	if err != nil {
		t.Fatal(err)
	}
	want := "\"a\":\n  -\n    - 1\n    - []\n  - \"b\": {}\n  - []\n\"c\":\n  \"d\":\n    \"e\": true\n\"f\": {}\n"
	if string(got) != want {
		t.Fatalf("got:\n%s", got)
	}
}

func TestHumanSummary(t *testing.T) {
	// OUT-004, RED-001/RED-002 in human output (including the workspace root).
	redactor := redact.New([]string{"EXAMPLE_API_TOKEN=canary-0123456789"}, "")
	text := render(t, sampleResult(), Options{Format: Human, Redactor: redactor})
	if !strings.Contains(text, "inspect_workspace") || !strings.Contains(text, "succeeded") || strings.Contains(text, "canary") ||
		strings.Contains(text, "\x1b") {
		t.Fatalf("%s", text)
	}
	failed := contract.NewResult(contract.InspectWorkspace, "req-2", "cli")
	failed.Outcome, failed.Error = contract.Failed, contract.InvalidInput("cannot inspect the native project declaration")
	text = render(t, failed, Options{Format: Human, Color: true})
	if !strings.Contains(text, "failed") || !strings.Contains(text, "invalid_input") || !strings.Contains(text, "\x1b[31m") {
		t.Fatalf("%q", text)
	}
}

func TestColorRule(t *testing.T) {
	// OUT-005, output.colorRule.
	cases := []struct {
		format, color, noColor string
		terminal, want         bool
	}{
		{"human", "auto", "", false, false},
		{"human", "auto", "", true, true},
		{"human", "auto", "1", true, false},
		{"human", "always", "", false, true},
		{"human", "always", "1", false, false},
		{"human", "never", "", true, false},
		{"json", "always", "", true, false},
		{"yaml", "always", "", true, false},
	}
	for _, c := range cases {
		if got := ColorEnabled(c.format, c.color, c.noColor, c.terminal); got != c.want {
			t.Fatalf("%+v: %v", c, got)
		}
	}
}

func TestBootstrapDiagnostic(t *testing.T) {
	// OUT-007; the home prefix and secrets are redacted.
	pattern := regexp.MustCompile(`^aibox: invalid_input: [a-z_]+: .+$`)
	redactor := redact.New([]string{"MY_TOKEN=canary-0123456789"}, "/home/u")
	line := BootstrapDiagnostic("malformed", "the settings file /home/u/.config/aibox/settings.json\nis not canary-0123456789", redactor)
	if !strings.HasSuffix(line, "\n") || strings.Count(line, "\n") != 1 || !pattern.MatchString(strings.TrimSuffix(line, "\n")) ||
		strings.Contains(line, "/home/u") || strings.Contains(line, "canary") {
		t.Fatalf("%q", line)
	}
	if line := BootstrapDiagnostic("usage_error", "", nil); !pattern.MatchString(strings.TrimSuffix(line, "\n")) {
		t.Fatalf("%q", line)
	}
}

func TestTerminalDetection(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	defer writer.Close()
	null, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer null.Close()
	if IsTerminal(writer.Fd()) || IsTerminal(null.Fd()) {
		t.Fatal("a pipe or /dev/null counted as a terminal")
	}
}
