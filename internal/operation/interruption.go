package operation

import (
	"context"
	"errors"
	"time"

	"github.com/projectious-work/aibox/internal/contract"
)

// InterruptionInspectionLimit bounds fresh post-cancellation observation.
const InterruptionInspectionLimit = 10 * time.Second

// InspectEffect checks the actual native state of one candidate after the
// child has exited. A false result means absence was proved; an error means
// effects are unknown. It must honor its context and verify exact identity.
type InspectEffect func(context.Context, contract.Effect) (occurred bool, err error)

// InterruptionReport carries the conservative outcome and schema-required
// effect lists for a cancellation or deadline. It never assumes that killing
// the child undid effects already performed by the child or its descendants.
type InterruptionReport struct {
	Outcome          contract.Outcome
	Error            contract.Error
	CompletedEffects []contract.Effect
	UnknownEffects   []contract.Effect
}

// InspectInterruption runs independent, bounded post-cancellation inspection.
// A missing target, failed inspection or observed effect is partial. Only
// proven absence of every candidate allows cancelled or timed_out. Callers
// persist the report in a record before returning a mutating result.
func InspectInterruption(cause error, candidates []contract.Effect, inspect InspectEffect) InterruptionReport {
	report := InterruptionReport{
		CompletedEffects: []contract.Effect{}, UnknownEffects: []contract.Effect{},
	}
	timedOut := errors.Is(cause, context.DeadlineExceeded)
	if timedOut {
		report.Outcome = contract.TimedOut
		report.Error = contract.Error{Code: "operation_timeout", Category: "timeout", Message: "operation deadline passed; native state was inspected", Retryable: false, NextAction: "none"}
	} else {
		report.Outcome = contract.Cancelled
		report.Error = contract.Error{Code: "operation_cancelled", Category: "interrupted", Message: "operation was cancelled; native state was inspected", Retryable: false, NextAction: "none"}
	}
	if len(candidates) == 0 || inspect == nil {
		report.UnknownEffects = append(report.UnknownEffects, contract.Effect{Resource: "unresolved-target", Action: "inspect", Status: "unknown"})
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), InterruptionInspectionLimit)
		defer cancel()
		for _, candidate := range candidates {
			if candidate.Resource == "" || candidate.Action == "" {
				report.UnknownEffects = append(report.UnknownEffects, contract.Effect{Resource: "unresolved-target", Action: "inspect", Status: "unknown"})
				continue
			}
			if ctx.Err() != nil {
				report.UnknownEffects = append(report.UnknownEffects, contract.Effect{Resource: candidate.Resource, Action: candidate.Action, Status: "unknown"})
				continue
			}
			occurred, err := inspect(ctx, candidate)
			if err != nil || ctx.Err() != nil {
				report.UnknownEffects = append(report.UnknownEffects, contract.Effect{Resource: candidate.Resource, Action: candidate.Action, Status: "unknown"})
			} else if occurred {
				report.CompletedEffects = append(report.CompletedEffects, contract.Effect{Resource: candidate.Resource, Action: candidate.Action, Status: "confirmed"})
			}
		}
	}
	if len(report.CompletedEffects) != 0 || len(report.UnknownEffects) != 0 {
		report.Outcome = contract.Partial
		report.Error.Message = "operation was interrupted; native effects occurred or remain unknown"
		report.Error.NextAction = "inspect_operation"
	}
	return report
}
