package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/projectious-work/aibox/internal/contract"
)

// buildPreview exercises the actual executable, not only command helpers.
func buildPreview(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "aibox")
	build := exec.Command("go", "build", "-o", binary, ".")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build preview: %v\n%s", err, output)
	}
	return binary
}

func TestPreviewBinaryInspection(t *testing.T) {
	binary := buildPreview(t)
	for _, args := range [][]string{{"--version"}, {"--help"}} {
		if output, err := exec.Command(binary, args...).CombinedOutput(); err != nil || len(output) == 0 {
			t.Fatalf("%v: %v %s", args, err, output)
		}
	}
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".devcontainer"), 0700); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(root, ".devcontainer", "devcontainer.json")
	if err := os.WriteFile(config, []byte(`{"image":"example.invalid/base:one"}`), 0600); err != nil {
		t.Fatal(err)
	}
	command := exec.Command(binary, "inspect", "--context", "local", "--project", root, "--format", "json")
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err != nil || stderr.Len() != 0 {
		t.Fatalf("inspect failed: %v stderr=%q", err, stderr.String())
	}
	var result contract.Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if err := result.ValidateEnvelope(); err != nil {
		t.Fatal(err)
	}
	if result.Operation != contract.InspectWorkspace || result.Target.WorkspaceRoot != root || result.InputDigest == "" {
		t.Fatalf("unexpected inspect result: %s", stdout.String())
	}
	for _, unavailable := range []string{"build", "up", "mcp"} {
		output, err := exec.Command(binary, unavailable).CombinedOutput()
		if err == nil || !strings.Contains(string(output), "not available") {
			t.Fatalf("%s unexpectedly available: %v %s", unavailable, err, output)
		}
	}
}

func TestPreviewBinaryMissingProject(t *testing.T) {
	binary := buildPreview(t)
	command := exec.Command(binary, "inspect", "--project", filepath.Join(t.TempDir(), "missing"), "--format", "json")
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	err := command.Run()
	if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 2 || stderr.Len() != 0 {
		t.Fatalf("unexpected invalid-project result: %v stderr=%q", err, stderr.String())
	}
	var result contract.Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Outcome != contract.Failed || result.Error.Code != "invalid_input" || result.Target != nil {
		t.Fatalf("invalid project leaked success or target: %s", stdout.String())
	}
}
