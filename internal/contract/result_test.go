package contract

import (
	"encoding/json"
	"strings"
	"testing"
)

func validResult() Result {
	return Result{
		SchemaVersion:    ResultSchemaVersion,
		Operation:        InspectWorkspace,
		RequestID:        "request-1",
		Actors:           Actors{Initiator: "test", Executor: "test"},
		Target:           &Target{Scope: "local", WorkspaceRoot: "/workspace"},
		InputDigest:      "sha256:" + strings.Repeat("a", 64),
		Outcome:          Succeeded,
		ChangedResources: []string{},
		Warnings:         []string{},
		Evidence:         []string{},
		Data:             map[string]any{"state": "absent"},
	}
}

func TestResultEnvelope(t *testing.T) {
	good := validResult()
	if err := good.ValidateEnvelope(); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(good)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{`"schemaVersion"`, `"operation"`, `"requestId"`, `"changedResources":[]`, `"warnings":[]`, `"evidence":[]`} {
		if !strings.Contains(string(encoded), field) {
			t.Fatalf("missing %s in %s", field, encoded)
		}
	}
	bad := good
	bad.RequestID = "not allowed"
	if bad.ValidateEnvelope() == nil {
		t.Fatal("accepted malformed request ID")
	}
	bad = good
	bad.InputDigest = "sha256:short"
	if bad.ValidateEnvelope() == nil {
		t.Fatal("accepted malformed digest")
	}
	bad = good
	bad.Error = &Error{Code: "internal_error", Message: "should not be present"}
	if bad.ValidateEnvelope() == nil {
		t.Fatal("accepted error on successful result")
	}
	bad = good
	bad.Outcome = Failed
	bad.Data = nil
	if bad.ValidateEnvelope() == nil {
		t.Fatal("accepted error outcome without error")
	}
	bad.Error = &Error{Code: "child_failed", Category: "child_failure", Message: "failed", NextAction: "inspect_operation"}
	if err := bad.ValidateEnvelope(); err != nil {
		t.Fatal(err)
	}
	bad.Error.NextAction = "run_arbitrary_command"
	if bad.ValidateEnvelope() == nil {
		t.Fatal("accepted unknown recovery action")
	}
	bad.Error.NextAction = "inspect_operation"
	bad.Outcome = Partial
	if bad.ValidateEnvelope() == nil {
		t.Fatal("accepted partial result without effect lists")
	}
	empty := []Effect{}
	bad.CompletedEffects, bad.UnknownEffects = &empty, &empty
	if err := bad.ValidateEnvelope(); err != nil {
		t.Fatal(err)
	}
	encoded, err = json.Marshal(bad)
	if err != nil || !strings.Contains(string(encoded), `"completedEffects":[]`) || !strings.Contains(string(encoded), `"unknownEffects":[]`) {
		t.Fatalf("missing required empty partial lists: %s, %v", encoded, err)
	}
}

func TestExitCode(t *testing.T) {
	for _, tc := range []struct {
		outcome  Outcome
		category string
		code     string
		want     int
	}{
		{Succeeded, "", "", 0},
		{RebuildRequired, "", "", 6},
		{Refused, "denied", "not_authorized", 3},
		{Failed, "invalid_input", "invalid_input", 2},
		{Failed, "missing_dependency", "dependency_missing", 4},
		{Failed, "child_failure", "child_failed", 5},
		{Partial, "timeout", "operation_timeout", 124},
		{Partial, "interrupted", "operation_cancelled", 130},
		{Failed, "internal", "recovery_required", 6},
	} {
		r := Result{Outcome: tc.outcome, Error: &Error{Category: tc.category, Code: tc.code}}
		if got := r.ExitCode(); got != tc.want {
			t.Errorf("%s/%s: got %d, want %d", tc.outcome, tc.code, got, tc.want)
		}
	}
}
