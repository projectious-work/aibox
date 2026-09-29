package operation

import (
	"errors"
	"time"
)

// validateOperationRecordTransition prevents a saved operation ID from being reused
// for a different request or moved backward after an observed effect. The
// caller holds the exact-environment lock across read, validate and Save.
func validateOperationRecordTransition(previous, next OperationRecord) error {
	if previous.OperationID != next.OperationID || previous.RequestFingerprint != next.RequestFingerprint ||
		previous.Operation != next.Operation || previous.ProjectRoot != next.ProjectRoot ||
		previous.ConfigPath != next.ConfigPath || previous.InputDigest != next.InputDigest ||
		previous.PolicyDigest != next.PolicyDigest || previous.GrantID != next.GrantID ||
		previous.CreatedAt != next.CreatedAt || previous.Executor != next.Executor {
		return errors.New("record immutable request identity changed")
	}
	previousTime, _ := time.Parse(time.RFC3339Nano, previous.UpdatedAt)
	nextTime, _ := time.Parse(time.RFC3339Nano, next.UpdatedAt)
	if nextTime.Before(previousTime) {
		return errors.New("record update time moved backward")
	}
	if previous.State == next.State {
		if terminalOperationRecordState(previous.State) {
			return errors.New("terminal record is immutable")
		}
		return nil
	}
	if !allowedOperationRecordTransition(previous.State, next.State) {
		return errors.New("invalid record state transition")
	}
	return nil
}

func terminalOperationRecordState(state string) bool {
	return oneOf(state, "succeeded", "no_change", "failed", "cancelled", "partial", "timed_out", "refused", "rebuild_required")
}

func allowedOperationRecordTransition(from, to string) bool {
	switch from {
	case "validated":
		return oneOf(to, "authorized", "refused", "failed")
	case "authorized":
		return oneOf(to, "executing", "inspected", "no_change", "refused", "failed")
	case "executing":
		return oneOf(to, "inspected", "partial", "failed", "cancelled", "timed_out")
	case "inspected":
		return terminalOperationRecordState(to)
	default:
		return false
	}
}
