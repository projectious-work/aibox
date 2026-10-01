package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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
	for _, unavailable := range []string{"build", "up"} {
		output, err := exec.Command(binary, unavailable).CombinedOutput()
		if err == nil || !strings.Contains(string(output), "not available") {
			t.Fatalf("%s unexpectedly available: %v %s", unavailable, err, output)
		}
	}
}

func callInspect(t *testing.T, session *mcp.ClientSession, root string) (contract.Result, bool) {
	t.Helper()
	response, err := session.CallTool(t.Context(), &mcp.CallToolParams{
		Name: "inspect_workspace", Arguments: map[string]any{"requestId": "parity-1", "projectRoot": root},
	})
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(response.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var result contract.Result
	if err := json.Unmarshal(encoded, &result); err != nil {
		t.Fatalf("decode structured MCP result: %v: %s", err, encoded)
	}
	if err := result.ValidateEnvelope(); err != nil {
		t.Fatal(err)
	}
	return result, response.IsError
}

func TestPreviewBinaryMCPParity(t *testing.T) {
	binary := buildPreview(t)
	client := mcp.NewClient(&mcp.Implementation{Name: "aibox-test", Version: "1"}, nil)
	session, err := client.Connect(t.Context(), &mcp.CommandTransport{Command: exec.Command(binary, "mcp", "serve", "--context", "local")}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = session.Close() })
	listed, err := session.ListTools(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(listed.Tools) != 1 || listed.Tools[0].Name != "inspect_workspace" || listed.Tools[0].Annotations == nil || !listed.Tools[0].Annotations.ReadOnlyHint {
		t.Fatalf("unexpected V1-03 MCP registry: %+v", listed.Tools)
	}
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".devcontainer"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".devcontainer", "devcontainer.json"), []byte(`{"image":"example.invalid/base:one"}`), 0600); err != nil {
		t.Fatal(err)
	}
	mcpResult, isError := callInspect(t, session, root)
	if isError || mcpResult.Outcome != contract.Succeeded || mcpResult.Target.WorkspaceRoot != root {
		t.Fatalf("MCP inspect: %+v error=%v", mcpResult, isError)
	}
	output, err := exec.Command(binary, "inspect", "--context", "local", "--project", root, "--format", "json").Output()
	if err != nil {
		t.Fatal(err)
	}
	var cliResult contract.Result
	if err := json.Unmarshal(output, &cliResult); err != nil {
		t.Fatal(err)
	}
	if cliResult.Operation != mcpResult.Operation || cliResult.Outcome != mcpResult.Outcome || cliResult.InputDigest != mcpResult.InputDigest || *cliResult.Target != *mcpResult.Target {
		t.Fatalf("CLI/MCP result mismatch: cli=%+v mcp=%+v", cliResult, mcpResult)
	}
	cliData, _ := json.Marshal(cliResult.Data)
	mcpData, _ := json.Marshal(mcpResult.Data)
	var cliFields, mcpFields map[string]any
	if err := json.Unmarshal(cliData, &cliFields); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(mcpData, &mcpFields); err != nil {
		t.Fatal(err)
	}
	delete(cliFields, "observedAt")
	delete(mcpFields, "observedAt")
	cliData, _ = json.Marshal(cliFields)
	mcpData, _ = json.Marshal(mcpFields)
	if !bytes.Equal(cliData, mcpData) {
		t.Fatalf("CLI/MCP data mismatch: cli=%s mcp=%s", cliData, mcpData)
	}
	missing := filepath.Join(root, "missing")
	mcpFailure, isError := callInspect(t, session, missing)
	if !isError || mcpFailure.Outcome != contract.Failed || mcpFailure.Error.Code != "invalid_input" || mcpFailure.Target != nil {
		t.Fatalf("MCP refusal: %+v error=%v", mcpFailure, isError)
	}
	command := exec.Command(binary, "inspect", "--context", "local", "--project", missing, "--format", "json")
	output, err = command.Output()
	if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 2 {
		t.Fatalf("CLI refusal exit: %v", err)
	}
	if err := json.Unmarshal(output, &cliResult); err != nil {
		t.Fatal(err)
	}
	if cliResult.Outcome != mcpFailure.Outcome || *cliResult.Error != *mcpFailure.Error {
		t.Fatalf("CLI/MCP refusal mismatch: cli=%+v mcp=%+v", cliResult, mcpFailure)
	}
	symlink := filepath.Join(t.TempDir(), "linked")
	if err := os.Symlink(root, symlink); err != nil {
		t.Fatal(err)
	}
	symlinkFailure, isError := callInspect(t, session, symlink)
	if !isError || symlinkFailure.Error == nil || symlinkFailure.Error.Code != "invalid_input" || symlinkFailure.Target != nil {
		t.Fatalf("MCP symlink refusal: %+v error=%v", symlinkFailure, isError)
	}
	t.Logf("stdio MCP registry=%s success=%s refusal=%s parity=passed", listed.Tools[0].Name, mcpResult.Outcome, mcpFailure.Error.Code)
}

func TestV104StarterCLIMCPParity(t *testing.T) {
	binary := buildPreview(t)
	client := mcp.NewClient(&mcp.Implementation{Name: "aibox-v1-04-test", Version: "1"}, nil)
	session, err := client.Connect(t.Context(), &mcp.CommandTransport{Command: exec.Command(binary, "mcp", "serve", "--context", "local")}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = session.Close() })
	listed, err := session.ListTools(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(listed.Tools) != 1 || listed.Tools[0].Name != "inspect_workspace" {
		t.Fatalf("unexpected V1-04 MCP registry: %+v", listed.Tools)
	}
	working, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(working, "..", ".."))
	for _, name := range []string{"minimal", "customized", "custom-user", "bind-home"} {
		root := filepath.Join(repoRoot, "spec", "v1", "examples", name)
		cliOutput, err := exec.Command(binary, "inspect", "--context", "local", "--project", root, "--format", "json").Output()
		if err != nil {
			t.Fatalf("%s CLI inspection: %v", name, err)
		}
		var cliResult contract.Result
		if err := json.Unmarshal(cliOutput, &cliResult); err != nil {
			t.Fatalf("%s CLI result: %v", name, err)
		}
		mcpResult, isError := callInspect(t, session, root)
		if isError || cliResult.Outcome != contract.Succeeded || mcpResult.Outcome != cliResult.Outcome ||
			mcpResult.Operation != cliResult.Operation || mcpResult.InputDigest != cliResult.InputDigest ||
			mcpResult.Target == nil || cliResult.Target == nil || *mcpResult.Target != *cliResult.Target {
			t.Fatalf("%s CLI/MCP envelope mismatch: cli=%+v mcp=%+v error=%v", name, cliResult, mcpResult, isError)
		}
		cliData, _ := json.Marshal(cliResult.Data)
		mcpData, _ := json.Marshal(mcpResult.Data)
		var cliFields, mcpFields map[string]any
		if err := json.Unmarshal(cliData, &cliFields); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(mcpData, &mcpFields); err != nil {
			t.Fatal(err)
		}
		delete(cliFields, "observedAt")
		delete(mcpFields, "observedAt")
		cliData, _ = json.Marshal(cliFields)
		mcpData, _ = json.Marshal(mcpFields)
		if !bytes.Equal(cliData, mcpData) {
			t.Fatalf("%s CLI/MCP data mismatch: cli=%s mcp=%s", name, cliData, mcpData)
		}
		if name == "customized" {
			configuration := cliFields["configuration"].([]any)
			if len(configuration) != 2 || configuration[1].(map[string]any)["key"] != "devcontainer.features" {
				t.Fatalf("customized selection is not visible through CLI/MCP inspection: %s", cliData)
			}
		}
		t.Logf("starter=%s operation=%s outcome=%s digest=%s cli-mcp-parity=passed", name, cliResult.Operation, cliResult.Outcome, cliResult.InputDigest)
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
