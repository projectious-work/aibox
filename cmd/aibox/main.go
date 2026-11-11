// Command aibox is the v1 Go CLI. V1-03 exposes local, read-only inspection
// through the CLI and MCP; host lifecycle operations arrive in later phases.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/projectious-work/aibox/internal/mcpserver"
	"github.com/projectious-work/aibox/internal/usecase"
)

// version may be set at build time; the default identifies a development CLI.
var version = "v1.0.0-dev"

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

// run keeps command routing testable without granting an implicit host client.
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprintln(stdout, "aibox v1 preview\nUsage: aibox inspect --context local --project PATH --format json\n       aibox mcp serve --context local\n       aibox version")
		return 0
	}
	if args[0] == "version" || args[0] == "--version" {
		fmt.Fprintln(stdout, version)
		return 0
	}
	if args[0] == "mcp" {
		flags := flag.NewFlagSet("mcp serve", flag.ContinueOnError)
		flags.SetOutput(stderr)
		contextValue := flags.String("context", "local", "MCP context")
		if len(args) < 2 || args[1] != "serve" || flags.Parse(args[2:]) != nil || flags.NArg() != 0 || *contextValue != "local" {
			fmt.Fprintln(stderr, "aibox: only mcp serve --context local is available in this preview")
			return 2
		}
		if err := mcpserver.Serve(context.Background(), version); err != nil {
			fmt.Fprintf(stderr, "aibox: MCP server failed: %v\n", err)
			return 5
		}
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
	requestID, err := usecase.NewRequestID()
	if err != nil {
		fmt.Fprintln(stderr, "aibox: cannot create a request ID")
		return 5
	}
	result := usecase.InspectWorkspace(requestID, "cli", *projectValue, false)
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
