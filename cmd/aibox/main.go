// Command aibox is the v1 Go CLI. It exposes local, read-only inspection
// through the CLI and MCP, with layered settings and diagnostic sinks
// (settings-and-logging@v1); host lifecycle operations arrive in later
// phases. This file is the composition root: it wires the adapters (files,
// environment, stdio, signals) to the use cases.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/projectious-work/aibox/internal/contract"
	"github.com/projectious-work/aibox/internal/logging"
	"github.com/projectious-work/aibox/internal/mcpserver"
	"github.com/projectious-work/aibox/internal/output"
	"github.com/projectious-work/aibox/internal/redact"
	"github.com/projectious-work/aibox/internal/settings"
	"github.com/projectious-work/aibox/internal/usecase"
)

// version may be set at build time; the default identifies a development CLI.
var version = "v1.0.0-dev"

const usageText = `aibox v1 preview
Usage: aibox inspect [--context local] [--project PATH] [--effective-config]
                     [--format human|json|yaml] [--color auto|always|never]
                     [--log-level trace|debug|info|warn|error|off]
                     [--log-format plain|jsonl] [--log-file PATH] [--timeout SECONDS]
                     [--settings-file PATH] [--env-file PATH]...
       aibox mcp serve [--context local] [--log-level LEVEL] [--log-format jsonl]
                     [--log-file PATH] [--settings-file PATH] [--env-file PATH]...
       aibox version
       aibox help`

// host holds the process's adapters, so that run is testable without
// touching the real process state.
type host struct {
	args    []string
	environ []string
	// workDir is the working directory at process start; aibox never
	// changes directory (CFG-053).
	workDir        string
	stdout         io.Writer
	stderr         io.Writer
	stdoutTerminal bool
	files          settings.FileSystem
	inspect        usecase.Inspector
	now            func() time.Time
	goos           string
	// serveMCP runs the MCP server; tests replace the stdio transport.
	serveMCP func(context.Context, mcpserver.Dependencies) error
}

func main() {
	// A write to a stdout or stderr whose reader has gone returns EPIPE
	// instead of killing the process (LOG-032).
	signal.Notify(make(chan os.Signal, 1), syscall.SIGPIPE)
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "aibox: cannot determine the working directory")
		os.Exit(5)
	}
	os.Exit(run(host{
		args: os.Args[1:], environ: os.Environ(), workDir: workDir,
		stdout: os.Stdout, stderr: os.Stderr, stdoutTerminal: output.IsTerminal(os.Stdout.Fd()),
		files: settings.OSFileSystem{}, now: time.Now, goos: runtime.GOOS, serveMCP: mcpserver.Serve,
	}))
}

// run routes one command line and returns the exit code.
func run(h host) int {
	line := parseCommandLine(h.args)
	if line.usage == nil {
		switch line.command {
		case commandHelp:
			fmt.Fprintln(h.stdout, usageText)
			return 0
		case commandVersion:
			fmt.Fprintln(h.stdout, version)
			return 0
		}
	}
	redactor := redact.New(h.environ, environValue(h.environ, "HOME"))
	if line.command == commandHelp || line.command == commandVersion {
		// version and help resolve no settings and emit no event (CFG-009, LOG-006).
		_, _ = io.WriteString(h.stderr, output.BootstrapDiagnostic(string(line.usage.Reason), line.usage.Message, redactor))
		return 2
	}
	requestID, err := usecase.NewRequestID()
	if err != nil {
		fmt.Fprintln(h.stderr, "aibox: cannot create a request ID")
		return 5
	}
	if line.usage != nil {
		return h.reject(line, requestID, line.usage, redactor)
	}
	mode := settings.ModeCLI
	if line.command == commandMCPServe {
		mode = settings.ModeMCP
	}
	loaded := settings.Load(h.files, settings.Request{
		Mode: mode, GOOS: h.goos, Environ: h.environ, WorkDir: h.workDir,
		SettingsFile: line.settingsFile, SettingsFileGiven: line.settingsFileGiven, EnvFiles: line.envFiles, Flags: line.flags,
	})
	if line.command == commandMCPServe {
		return h.serve(line, requestID, loaded, redactor)
	}
	return h.inspectWorkspace(line, requestID, loaded, redactor)
}

// reject reports a failure before the operation starts: one settings.rejected
// event on stderr (LOG-006), then either the invalid_request envelope
// (OUT-006) or the bootstrap diagnostic as the last stderr line (OUT-007).
func (h host) reject(line commandLine, requestID string, rejection *settings.Rejection, redactor *redact.Redactor) int {
	format := logging.Plain
	if line.logFormatJSONL || line.command == commandMCPServe {
		format = logging.JSONL
	}
	if !line.logLevelOff {
		logging.WriteRejection(h.stderr, format, requestID, rejection.Fields(), redactor, h.now())
	}
	if line.command == commandInspect && line.envelopeFormat != "" {
		result := contract.InvalidRequestResult(requestID, "cli", redactor.All(rejectionMessage(rejection)))
		if err := output.Write(h.stdout, result, output.Options{Format: line.envelopeFormat, Redactor: redactor}); err != nil {
			return 5
		}
		return result.ExitCode()
	}
	_, _ = io.WriteString(h.stderr, output.BootstrapDiagnostic(string(rejection.Reason), rejectionMessage(rejection), redactor))
	return 2
}

func rejectionMessage(rejection *settings.Rejection) string {
	message := rejection.Message
	if rejection.Suggestion != "" {
		message += "; did you mean " + rejection.Suggestion + "?"
	}
	return message
}

// startLogger opens the file sink, if logging.file is set, and the logger.
func (h host) startLogger(resolved *settings.Resolved, redactor *redact.Redactor) (*logging.Logger, func()) {
	config := logging.Config{
		Level: resolved.Text(settings.LoggingLevel), Format: logging.Format(resolved.Text(settings.LoggingFormat)),
		Stderr: h.stderr, Redactor: redactor, Now: h.now, Limits: logging.DefaultLimits,
	}
	release := func() {}
	if path := resolved.Text(settings.LoggingFile); path != "" {
		sink, failure := logging.OpenFileSink(logging.FileSinkConfig{
			Path: path, RotationMiB: resolved.Integer(settings.LoggingRotationMiB),
			RetentionFiles: resolved.Integer(settings.LoggingRetentionFiles), RetentionDays: resolved.Integer(settings.LoggingRetentionDays),
			LockWait: logging.DefaultLimits.Flush, Now: h.now,
		})
		config.FileFailure = failure
		if sink != nil {
			config.File = sink
			release = func() { _ = sink.Close() }
		}
	}
	return logging.Start(config), release
}

func resolvedEvent(resolved *settings.Resolved) logging.Event {
	layers := make([]string, 0, len(resolved.LayersSet))
	for _, layer := range resolved.LayersSet {
		layers = append(layers, string(layer))
	}
	return logging.Event{Severity: logging.Debug, Name: "settings.resolved", Component: "settings",
		Fields: map[string]string{"layers": strings.Join(layers, ",")}}
}

// inspectWorkspace runs aibox inspect: settings first (CFG-006), then the
// operation under SIGINT and SIGTERM handling, then the result on stdout.
func (h host) inspectWorkspace(line commandLine, requestID string, loaded *settings.Loaded, redactor *redact.Redactor) int {
	project := line.project
	if !filepath.IsAbs(project) {
		project = filepath.Join(h.workDir, project)
	}
	resolved, rejection := loaded.Resolve(settings.ReadProject(project))
	if rejection != nil {
		return h.reject(line, requestID, rejection, redactor)
	}
	logger, release := h.startLogger(resolved, redactor)
	defer release()
	scope := logger.Scope(requestID)
	scope.Emit(resolvedEvent(resolved))
	logger.ReportOpenFailure(scope)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	workspace := usecase.Workspace{Inspect: h.inspect, Health: logger, Redactor: redactor}
	result := workspace.InspectWorkspace(ctx, resolved, usecase.InspectRequest{
		RequestID: requestID, Initiator: "cli", ProjectRoot: project, EffectiveConfig: line.effectiveConfig,
	}, scope)
	stop()
	logger.Close(scope)
	if logger.FileSinkFailed() {
		result.AddWarning(contract.WarningLogSinkDegraded)
	}
	if err := result.ValidateEnvelope(); err != nil {
		fmt.Fprintf(h.stderr, "aibox: invalid operation result: %v\n", err)
		return 5
	}
	format := resolved.Text(settings.OutputFormat)
	color := output.ColorEnabled(format, resolved.Text(settings.OutputColor), environValue(h.environ, "NO_COLOR"), h.stdoutTerminal)
	if err := output.Write(h.stdout, result, output.Options{Format: format, Color: color, Redactor: redactor}); err != nil {
		return 5
	}
	return result.ExitCode()
}

// serve runs mcp serve with the settings resolved at start (CFG-055).
func (h host) serve(line commandLine, requestID string, loaded *settings.Loaded, redactor *redact.Redactor) int {
	resolved, rejection := loaded.Resolve(nil)
	if rejection != nil {
		return h.reject(line, requestID, rejection, redactor)
	}
	logger, release := h.startLogger(resolved, redactor)
	defer release()
	process := logger.Scope("process-" + requestID)
	process.Emit(logging.Event{Severity: logging.Info, Name: "process.started", Component: "mcp"})
	process.Emit(resolvedEvent(resolved))
	logger.ReportOpenFailure(process)
	err := h.serveMCP(context.Background(), mcpserver.Dependencies{
		Version: version, Settings: loaded, Logger: logger,
		Workspace: usecase.Workspace{Inspect: h.inspect, Health: logger, Redactor: redactor},
	})
	process.Emit(logging.Event{Severity: logging.Info, Name: "process.finished", Component: "mcp"})
	logger.Close(process)
	if err != nil {
		fmt.Fprintln(h.stderr, "aibox: MCP server failed: "+redactor.All(err.Error()))
		return 5
	}
	return 0
}

func environValue(environ []string, name string) string {
	for _, entry := range environ {
		if key, value, found := strings.Cut(entry, "="); found && key == name {
			return value
		}
	}
	return ""
}
