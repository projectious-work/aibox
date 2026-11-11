package operation

import (
	"context"
	"errors"
	"testing"

	"github.com/projectious-work/aibox/internal/contract"
)

func TestInspectInterruptionClassifiesObservedAndUnknownEffects(t *testing.T) {
	candidates := []contract.Effect{{Resource: "container:one", Action: "start"}, {Resource: "container:two", Action: "start"}}
	report := InspectInterruption(context.Canceled, candidates, func(_ context.Context, candidate contract.Effect) (bool, error) {
		if candidate.Resource == "container:one" {
			return true, nil
		}
		return false, errors.New("runtime unavailable")
	})
	if report.Outcome != contract.Partial || len(report.CompletedEffects) != 1 || len(report.UnknownEffects) != 1 || report.Error.NextAction != "inspect_operation" {
		t.Fatalf("unsafe interruption classification: %+v", report)
	}
}

func TestInspectInterruptionNeedsProofForNoEffect(t *testing.T) {
	candidate := []contract.Effect{{Resource: "container:one", Action: "start"}}
	absent := func(context.Context, contract.Effect) (bool, error) { return false, nil }
	for _, tc := range []struct {
		cause error
		want  contract.Outcome
	}{
		{context.Canceled, contract.Cancelled},
		{context.DeadlineExceeded, contract.TimedOut},
	} {
		report := InspectInterruption(tc.cause, candidate, absent)
		if report.Outcome != tc.want || len(report.CompletedEffects) != 0 || len(report.UnknownEffects) != 0 {
			t.Fatalf("proven absence: %+v", report)
		}
	}
	for _, report := range []InterruptionReport{
		InspectInterruption(context.Canceled, nil, absent),
		InspectInterruption(context.Canceled, candidate, nil),
		InspectInterruption(context.Canceled, candidate, func(context.Context, contract.Effect) (bool, error) { return false, errors.New("unknown") }),
	} {
		if report.Outcome != contract.Partial || len(report.UnknownEffects) == 0 {
			t.Fatalf("unknown effects reported harmless: %+v", report)
		}
	}
}
