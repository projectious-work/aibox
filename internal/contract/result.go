// Package contract contains wire-facing operation values shared by future CLI
// and MCP adapters. It has no dependency on runtime or presentation packages.
package contract

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const ResultSchemaVersion = "aibox.operation-result/v1"

type Operation string

const (
	BuildEnvironment   Operation = "build_environment"
	StartEnvironment   Operation = "start_environment"
	StopEnvironment    Operation = "stop_environment"
	RemoveEnvironment  Operation = "remove_environment"
	RebuildEnvironment Operation = "rebuild_environment"
	InspectEnvironment Operation = "inspect_environment"
	CheckEnvironment   Operation = "check_environment"
	InspectWorkspace   Operation = "inspect_workspace"
	CheckWorkspace     Operation = "check_workspace"
	RefreshWorkspace   Operation = "refresh_workspace"
	ReadLogs           Operation = "read_logs"
	InspectOperation   Operation = "inspect_operation"
	MigratePreview     Operation = "migrate_preview"
	MigrateApply       Operation = "migrate_apply"
	MigrateRollback    Operation = "migrate_rollback"
	InvalidRequest     Operation = "invalid_request"
)

type Outcome string

const (
	Succeeded       Outcome = "succeeded"
	NoChange        Outcome = "no_change"
	RebuildRequired Outcome = "rebuild_required"
	Failed          Outcome = "failed"
	Cancelled       Outcome = "cancelled"
	Partial         Outcome = "partial"
	Refused         Outcome = "refused"
	TimedOut        Outcome = "timed_out"
)

type Actors struct {
	Initiator string `json:"initiator"`
	Executor  string `json:"executor"`
}

type Target struct {
	Scope          string `json:"scope"`
	WorkspaceRoot  string `json:"workspaceRoot"`
	RuntimeContext string `json:"runtimeContext,omitempty"`
	ResourceID     string `json:"resourceId,omitempty"`
}

type Error struct {
	Code       string `json:"code"`
	Category   string `json:"category"`
	Message    string `json:"message"`
	Retryable  bool   `json:"retryable"`
	NextAction string `json:"nextAction"`
}

type Effect struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Status   string `json:"status"`
}

// Result follows operation-result.schema.json. Data is deliberately opaque at
// this layer; operation-specific constructors and schema checks are still due.
type Result struct {
	SchemaVersion    string    `json:"schemaVersion"`
	Operation        Operation `json:"operation"`
	RequestID        string    `json:"requestId"`
	Actors           Actors    `json:"actors"`
	Target           *Target   `json:"target,omitempty"`
	InputDigest      string    `json:"inputDigest,omitempty"`
	Outcome          Outcome   `json:"outcome"`
	ChangedResources []string  `json:"changedResources"`
	Warnings         []string  `json:"warnings"`
	Evidence         []string  `json:"evidence"`
	Data             any       `json:"data,omitempty"`
	Error            *Error    `json:"error,omitempty"`
	CompletedEffects *[]Effect `json:"completedEffects,omitempty"`
	UnknownEffects   *[]Effect `json:"unknownEffects,omitempty"`
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
