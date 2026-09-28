// Command aibox is the v1 Go CLI. V1-03 intentionally exposes only local,
// read-only inspection; host lifecycle operations arrive in later phases.
package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/projectious-work/aibox/internal/contract"
	"github.com/projectious-work/aibox/internal/inspection"
)

// version may be set at build time; the default identifies a development CLI.
var version = "v1.0.0-dev"

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

// run keeps command routing testable without granting an implicit host client.
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprintln(stdout, "aibox v1 preview\nUsage: aibox inspect --context local --project PATH --format json\n       aibox version")
		return 0
	}
	if args[0] == "version" || args[0] == "--version" {
		fmt.Fprintln(stdout, version)
		return 0
	}
	if args[0] != "inspect" {
		fmt.Fprintf(stderr, "aibox: %q is not available in this preview\n", args[0])
		return 2
	}
	flags := flag.NewFlagSet("inspect", flag.ContinueOnError)
	flags.SetOutput(stderr)
	contextValue := flags.String("context", "local", "inspection context")
	projectValue := flags.String("project", ".", "native project directory")
	formatValue := flags.String("format", "json", "result format")
	if err := flags.Parse(args[1:]); err != nil {
		return 2
	}
	if flags.NArg() != 0 || *formatValue != "json" {
		fmt.Fprintln(stderr, "aibox: inspect accepts no positional arguments and only --format json")
		return 2
	}
	if *contextValue != "local" {
		fmt.Fprintln(stderr, "aibox: operator inspection is not available in this preview")
		return 2
	}
	requestBytes := make([]byte, 16)
	if _, err := rand.Read(requestBytes); err != nil {
		fmt.Fprintln(stderr, "aibox: cannot create a request ID")
		return 5
	}
	requestID := hex.EncodeToString(requestBytes)
	result := contract.Result{
		SchemaVersion:    contract.ResultSchemaVersion,
		Operation:        contract.InspectWorkspace,
		RequestID:        requestID,
		Actors:           contract.Actors{Initiator: "cli", Executor: "aibox"},
		ChangedResources: []string{}, Warnings: []string{}, Evidence: []string{},
	}
	absProject, err := filepath.Abs(*projectValue)
	if err == nil {
		var root, digest string
		var data inspection.WorkspaceData
		root, digest, data, err = inspection.Workspace(absProject)
		if err == nil {
			result.Target = &contract.Target{Scope: "local", WorkspaceRoot: root}
			result.InputDigest = digest
			result.Outcome = contract.Succeeded
			result.Data = data
		}
	}
	if err != nil {
		result.Outcome = contract.Failed
		result.Error = &contract.Error{
			Code: "invalid_input", Category: "invalid_input",
			Message:   "cannot inspect the native project declaration",
			Retryable: false, NextAction: "edit_configuration",
		}
	}
	if err := result.ValidateEnvelope(); err != nil {
		fmt.Fprintf(stderr, "aibox: invalid operation result: %v\n", err)
		return 5
	}
	if err := json.NewEncoder(stdout).Encode(result); err != nil {
		fmt.Fprintf(stderr, "aibox: cannot write operation result: %v\n", err)
		return 5
	}
	return result.ExitCode()
}
