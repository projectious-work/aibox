// Package contract contains wire-facing operation values shared by future CLI
// and MCP adapters. It has no dependency on runtime or presentation packages.
package contract

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ResultSchemaVersion identifies the common operation-result wire contract.
const ResultSchemaVersion = "aibox.operation-result/v1"

// Operation identifies the use case that produced a result. Its values match
// the operation union in spec/v1/operation-result.schema.json.
type Operation string

const (
	// BuildEnvironment builds a Dev Container image without starting it.
	BuildEnvironment Operation = "build_environment"
	// StartEnvironment starts or reuses an exact environment.
	StartEnvironment Operation = "start_environment"
	// StopEnvironment stops an exact environment without removing retained data.
	StopEnvironment Operation = "stop_environment"
	// RemoveEnvironment removes an exact environment without broad pruning.
	RemoveEnvironment Operation = "remove_environment"
	// RebuildEnvironment replaces an exact environment after a build preflight.
	RebuildEnvironment Operation = "rebuild_environment"
	// InspectEnvironment reads declared and observed operator-side state.
	InspectEnvironment Operation = "inspect_environment"
	// CheckEnvironment runs operator-side diagnostics.
	CheckEnvironment Operation = "check_environment"
	// InspectWorkspace reads workspace-local state without host authority.
	InspectWorkspace Operation = "inspect_workspace"
	// CheckWorkspace runs workspace-local diagnostics.
	CheckWorkspace Operation = "check_workspace"
	// RefreshWorkspace updates only aibox-managed local workspace output.
	RefreshWorkspace Operation = "refresh_workspace"
	// ReadLogs reads a declared, bounded diagnostic source.
	ReadLogs Operation = "read_logs"
	// InspectOperation reads a durable operation record.
	InspectOperation Operation = "inspect_operation"
	// MigratePreview previews a one-time v0-to-v1 conversion.
	MigratePreview Operation = "migrate_preview"
	// MigrateApply applies an explicitly selected migration plan.
	MigrateApply Operation = "migrate_apply"
	// MigrateRollback rolls back an explicitly selected migration plan.
	MigrateRollback Operation = "migrate_rollback"
	// InvalidRequest represents a request rejected before its operation is known.
	InvalidRequest Operation = "invalid_request"
)

// Outcome records the observed result of a use case, not the child's exit
// status alone. In particular, a failed postcondition cannot be succeeded.
type Outcome string

const (
	// Succeeded means the requested postcondition was observed.
	Succeeded Outcome = "succeeded"
	// NoChange means the requested state was already present.
	NoChange Outcome = "no_change"
	// RebuildRequired means the current workspace needs an explicit rebuild.
	RebuildRequired Outcome = "rebuild_required"
	// Failed means the operation failed without a distinct partial outcome.
	Failed Outcome = "failed"
	// Cancelled means cancellation left no observed effects.
	Cancelled Outcome = "cancelled"
	// Partial means some effects occurred or their absence cannot be proven.
	Partial Outcome = "partial"
	// Refused means authorization or a safety precondition denied the request.
	Refused Outcome = "refused"
	// TimedOut means the deadline passed without observed effects.
	TimedOut Outcome = "timed_out"
)

// Actors identifies the attributed initiator and executor. These identifiers
// are evidence labels, not authorization grants.
type Actors struct {
	// Initiator names the externally established requester.
	Initiator string `json:"initiator"`
	// Executor names the process or principal performing the use case.
	Executor string `json:"executor"`
}

// Target identifies the resolved scope and workspace. RuntimeContext and
// ResourceID are omitted until those identities are actually known.
type Target struct {
	// Scope is local or operator; local scope has no host runtime authority.
	Scope string `json:"scope"`
	// WorkspaceRoot is the canonical absolute workspace path.
	WorkspaceRoot string `json:"workspaceRoot"`
	// RuntimeContext is the selected operator-side runtime context, when known.
	RuntimeContext string `json:"runtimeContext,omitempty"`
	// ResourceID is an inspected native resource ID, when known.
	ResourceID string `json:"resourceId,omitempty"`
}

// Error is a structured domain failure with a closed recovery action. Message
// is explanatory text and must never be interpreted as a command to execute.
type Error struct {
	// Code is the stable machine-readable error code.
	Code string `json:"code"`
	// Category groups the error into a documented CLI exit class.
	Category string `json:"category"`
	// Message is a bounded, redacted human explanation.
	Message string `json:"message"`
	// Retryable says whether a later attempt may be meaningful, not automatic.
	Retryable bool `json:"retryable"`
	// NextAction names a documented recovery choice, never arbitrary code.
	NextAction string `json:"nextAction"`
}

// Effect records a resource action whose status was inspected or remains
// unknown after a partial operation.
type Effect struct {
	// Resource is the specific native or managed resource reference.
	Resource string `json:"resource"`
	// Action describes the attempted effect.
	Action string `json:"action"`
	// Status is confirmed or unknown.
	Status string `json:"status"`
}

// Result follows operation-result.schema.json. Data is deliberately opaque at
// this layer; operation-specific constructors and schema checks are still due.
// The three common lists must be non-nil so JSON contains arrays, not null.
type Result struct {
	// SchemaVersion is ResultSchemaVersion.
	SchemaVersion string `json:"schemaVersion"`
	// Operation identifies the use case or an invalid bootstrap request.
	Operation Operation `json:"operation"`
	// RequestID correlates requests, results and records without granting access.
	RequestID string `json:"requestId"`
	// Actors attributes the request and execution.
	Actors Actors `json:"actors"`
	// Target is absent when discovery failed before identity resolution.
	Target *Target `json:"target,omitempty"`
	// InputDigest is absent when the declared input set could not be resolved.
	InputDigest string `json:"inputDigest,omitempty"`
	// Outcome reports observed state, including uncertainty after effects.
	Outcome Outcome `json:"outcome"`
	// ChangedResources contains bounded resource references, not display names.
	ChangedResources []string `json:"changedResources"`
	// Warnings contains bounded, non-secret human messages.
	Warnings []string `json:"warnings"`
	// Evidence contains references to redacted evidence, not raw child output.
	Evidence []string `json:"evidence"`
	// Data is an operation-specific object on successful outcomes.
	Data any `json:"data,omitempty"`
	// Error is required for failed, refused, cancelled, timed-out or partial outcomes.
	Error *Error `json:"error,omitempty"`
	// CompletedEffects is present, even if empty, only for partial outcomes.
	CompletedEffects *[]Effect `json:"completedEffects,omitempty"`
	// UnknownEffects is present, even if empty, only for partial outcomes.
	UnknownEffects *[]Effect `json:"unknownEffects,omitempty"`
}

var requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
var digestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

// ValidateEnvelope checks the common schema invariants before an adapter
// serializes a result. Full operation-specific validation remains separate.
func (r Result) ValidateEnvelope() error {
	if r.SchemaVersion != ResultSchemaVersion {
		return errors.New("unsupported result schema version")
	}
	if !validOperation(r.Operation) || !requestIDPattern.MatchString(r.RequestID) {
		return errors.New("invalid operation or request ID")
	}
	if r.Actors.Initiator == "" || r.Actors.Executor == "" {
		return errors.New("both actors are required")
	}
	if len(r.Actors.Initiator) > 4096 || len(r.Actors.Executor) > 4096 {
		return errors.New("actor identifier exceeds limit")
	}
	if r.ChangedResources == nil || r.Warnings == nil || r.Evidence == nil {
		return errors.New("result lists must be present")
	}
	if r.Target != nil && (r.Target.Scope != "local" && r.Target.Scope != "operator" || !strings.HasPrefix(r.Target.WorkspaceRoot, "/") || len(r.Target.WorkspaceRoot) > 4096) {
		return errors.New("invalid target")
	}
	if r.InputDigest != "" && !digestPattern.MatchString(r.InputDigest) {
		return errors.New("invalid input digest")
	}
	if r.Operation == InvalidRequest && r.Outcome != Failed && r.Outcome != Refused {
		return errors.New("invalid_request may only fail or be refused")
	}
	switch r.Outcome {
	case Succeeded, NoChange, RebuildRequired:
		if r.Target == nil || r.InputDigest == "" || r.Data == nil || r.Error != nil {
			return errors.New("successful result needs target, digest and data but no error")
		}
	case Failed, Cancelled, Partial, Refused, TimedOut:
		if r.Error == nil || r.Error.Code == "" || r.Error.Message == "" {
			return errors.New("failure result needs a structured error")
		}
		if !requestIDPattern.MatchString(r.Error.Code) || len(r.Error.Message) > 4096 || !validCategory(r.Error.Category) || !validNextAction(r.Error.NextAction) {
			return errors.New("invalid structured error")
		}
	default:
		return fmt.Errorf("unknown outcome %q", r.Outcome)
	}
	if r.Outcome == Partial {
		if r.CompletedEffects == nil || r.UnknownEffects == nil || *r.CompletedEffects == nil || *r.UnknownEffects == nil {
			return errors.New("partial result needs both effect lists")
		}
	} else if r.CompletedEffects != nil || r.UnknownEffects != nil {
		return errors.New("effects are only valid for partial results")
	}
	return nil
}

func validCategory(category string) bool {
	switch category {
	case "invalid_input", "denied", "missing_dependency", "incompatible", "timeout", "interrupted", "child_failure", "partial_failure", "internal":
		return true
	default:
		return false
	}
}

func validNextAction(action string) bool {
	switch action {
	case "none", "edit_configuration", "inspect_environment", "inspect_operation", "request_authorization", "install_dependency", "retry_read", "manual_recovery", "rebuild_environment", "restart_session", "read_guide":
		return true
	default:
		return false
	}
}

func validOperation(o Operation) bool {
	switch o {
	case BuildEnvironment, StartEnvironment, StopEnvironment, RemoveEnvironment,
		RebuildEnvironment, InspectEnvironment, CheckEnvironment, InspectWorkspace,
		CheckWorkspace, RefreshWorkspace, ReadLogs, InspectOperation, MigratePreview,
		MigrateApply, MigrateRollback, InvalidRequest:
		return true
	default:
		return false
	}
}

// ExitCode preserves the stable CLI error classes in chapter 18.
func (r Result) ExitCode() int {
	if r.Outcome == Succeeded || r.Outcome == NoChange {
		return 0
	}
	if r.Outcome == TimedOut || r.Error != nil && r.Error.Code == "operation_timeout" {
		return 124
	}
	if r.Outcome == Cancelled || r.Error != nil && r.Error.Code == "operation_cancelled" {
		return 130
	}
	if r.Error != nil {
		switch r.Error.Code {
		case "rebuild_required", "session_restart_required", "migration_conflict", "recovery_required":
			return 6
		}
		switch r.Error.Category {
		case "invalid_input":
			return 2
		case "denied":
			return 3
		case "missing_dependency", "incompatible":
			return 4
		}
	}
	if r.Outcome == RebuildRequired {
		return 6
	}
	return 5
}
