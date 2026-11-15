package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/projectious-work/aibox/internal/contract"
	"github.com/projectious-work/aibox/internal/logging"
)

// Class is the operation class of execution.operationClass.
type Class string

const (
	// Read operations observe state; their default limit is 30 seconds.
	Read Class = "read"
	// Mutating operations change state; their default limit is 3600 seconds.
	Mutating Class = "mutating"
)

// Events is the port through which a use case emits the log events of its
// request. The logging scope of the request implements it.
type Events interface {
	Emit(logging.Event)
}

// SinkHealth is the port that reports whether the file sink has failed.
type SinkHealth interface {
	FileSinkFailed() bool
}

// Operation describes one execution of a use case.
type Operation struct {
	// Name is the operation of the result.
	Name contract.Operation
	// Class selects the required-sink rule.
	Class Class
	// TimeLimit is execution.timeoutSeconds or the class default.
	TimeLimit time.Duration
	// RequestID and Initiator go into every result.
	RequestID, Initiator string
	// LoggingRequired is logging.required.
	LoggingRequired bool
}

// Execute runs body under the operation's time limit and the cancellation
// of ctx. A mutating operation is refused with operation_record_unavailable
// before body runs when logging.required is true and the file sink has
// failed (LOG-031, D16). A body that outlives its deadline yields timed_out
// (EXE-002); a cancelled ctx yields cancelled (LOG-008).
func Execute(ctx context.Context, operation Operation, health SinkHealth, body func(context.Context) contract.Result) contract.Result {
	if operation.Class == Mutating && operation.LoggingRequired && health.FileSinkFailed() {
		result := contract.NewResult(operation.Name, operation.RequestID, operation.Initiator)
		result.Outcome = contract.Refused
		result.Error = contract.OperationRecordUnavailable("the required log file sink has failed; the operation was not started")
		return result
	}
	limited, cancel := context.WithTimeout(ctx, operation.TimeLimit)
	defer cancel()
	done := make(chan contract.Result, 1)
	go func() { done <- body(limited) }()
	select {
	case result := <-done:
		if limited.Err() == nil {
			return result
		}
	case <-limited.Done():
	}
	return interrupted(operation, limited.Err())
}

func interrupted(operation Operation, cause error) contract.Result {
	result := contract.NewResult(operation.Name, operation.RequestID, operation.Initiator)
	if errors.Is(cause, context.DeadlineExceeded) {
		result.Outcome = contract.TimedOut
		result.Error = contract.OperationTimeout("the operation exceeded its time limit")
		return result
	}
	result.Outcome = contract.Cancelled
	result.Error = contract.OperationCancelled("the operation was interrupted")
	return result
}
