// Package usecase holds operations shared by the CLI and MCP adapters.
package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"path/filepath"
	"time"

	"github.com/projectious-work/aibox/internal/contract"
	"github.com/projectious-work/aibox/internal/inspection"
	"github.com/projectious-work/aibox/internal/logging"
	"github.com/projectious-work/aibox/internal/redact"
	"github.com/projectious-work/aibox/internal/settings"
)

// NewRequestID creates a correlation identifier without granting authority.
func NewRequestID() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}

// Inspector is the port to the read-only workspace inspection; it returns
// the canonical root, the input digest and the payload.
type Inspector func(ctx context.Context, root string) (string, string, contract.WorkspaceData, error)

// Workspace is the inspect_workspace use case with its ports.
type Workspace struct {
	// Inspect reads the native declaration; inspection.Workspace by default.
	Inspect Inspector
	// Health reports the file sink state.
	Health SinkHealth
	// Redactor redacts the effective-configuration view.
	Redactor *redact.Redactor
}

// InspectRequest is one inspect_workspace request from the CLI or MCP.
type InspectRequest struct {
	RequestID string
	// Initiator is cli or mcp; an attribution label, never a grant.
	Initiator string
	// ProjectRoot is the project path as given, made absolute.
	ProjectRoot string
	// CanonicalOnly requires a canonical absolute root (MCP).
	CanonicalOnly bool
	// EffectiveConfig adds the settings entries and settingsSources.
	EffectiveConfig bool
}

// InspectWorkspace executes the same read-only operation for CLI and MCP
// with the settings resolved before it starts, and emits the events of
// logging.eventOrder: operation.started, operation.validated on success,
// operation.cancelled on interruption, operation.finished.
func (w Workspace) InspectWorkspace(ctx context.Context, resolved *settings.Resolved, request InspectRequest, events Events) contract.Result {
	inspect := w.Inspect
	if inspect == nil {
		inspect = inspection.Workspace
	}
	member := func(name string, severity logging.Severity) logging.Event {
		return logging.Event{Severity: severity, Name: name, Component: "usecase", Operation: string(contract.InspectWorkspace),
			Initiator: request.Initiator}
	}
	started := member("operation.started", logging.Info)
	started.Fields = map[string]string{"projectRoot": request.ProjectRoot}
	events.Emit(started)
	operation := Operation{
		Name: contract.InspectWorkspace, Class: Read, RequestID: request.RequestID, Initiator: request.Initiator,
		TimeLimit:       time.Duration(resolved.TimeLimitSeconds(settings.ReadClassLimitSeconds)) * time.Second,
		LoggingRequired: resolved.Flag(settings.LoggingRequired),
	}
	result := Execute(ctx, operation, w.Health, func(ctx context.Context) contract.Result {
		return w.inspect(ctx, inspect, resolved, request)
	})
	if result.Outcome == contract.Succeeded {
		events.Emit(member("operation.validated", logging.Debug))
	}
	if result.Outcome == contract.Cancelled {
		cancelled := member("operation.cancelled", logging.Warn)
		cancelled.Outcome, cancelled.ErrorCategory = string(result.Outcome), result.Error.Category
		events.Emit(cancelled)
	}
	finished := member("operation.finished", logging.Info)
	finished.Outcome = string(result.Outcome)
	if result.Error != nil {
		finished.ErrorCategory = result.Error.Category
	}
	events.Emit(finished)
	if w.Health.FileSinkFailed() {
		result.AddWarning(contract.WarningLogSinkDegraded)
	}
	return result
}

func (w Workspace) inspect(ctx context.Context, inspect Inspector, resolved *settings.Resolved, request InspectRequest) contract.Result {
	result := contract.NewResult(contract.InspectWorkspace, request.RequestID, request.Initiator)
	var err error
	projectRoot := request.ProjectRoot
	if request.CanonicalOnly && (!filepath.IsAbs(projectRoot) || filepath.Clean(projectRoot) != projectRoot) {
		err = errors.New("MCP project root must be a canonical absolute path")
	}
	if err == nil {
		var root, digest string
		var data contract.WorkspaceData
		root, digest, data, err = inspect(ctx, projectRoot)
		if err == nil && request.CanonicalOnly && root != projectRoot {
			err = errors.New("MCP project root must resolve to itself")
		}
		if err == nil {
			data.Configuration = redactNative(w.Redactor, data.Configuration)
			if request.EffectiveConfig {
				entries, sources := EffectiveConfiguration(resolved, w.Redactor, settings.ReadClassLimitSeconds)
				data.Configuration = append(data.Configuration, entries...)
				data.SettingsSources = sources
			}
			result.Target = &contract.Target{Scope: "local", WorkspaceRoot: root}
			result.InputDigest = digest
			result.Outcome = contract.Succeeded
			result.Data = data
		}
	}
	if err != nil {
		result.Outcome = contract.Failed
		result.Error = contract.InvalidInput("cannot inspect the native project declaration")
	}
	return result
}
