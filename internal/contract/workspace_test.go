package contract

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/projectious-work/aibox/internal/schematest"
)

func TestInvalidRequestEnvelope(t *testing.T) {
	// OUT-006, D19: invalid_request, failed, invalid_input, no target or digest.
	result := InvalidRequestResult("req-1", "cli", "the settings file is not JSON")
	if err := result.ValidateEnvelope(); err != nil {
		t.Fatal(err)
	}
	encoded, _ := json.Marshal(result)
	schema, err := schematest.Load("schemas/operation-result.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := schema.Validate(encoded); err != nil {
		t.Fatalf("%v: %s", err, encoded)
	}
	if strings.Contains(string(encoded), "target") || strings.Contains(string(encoded), "inputDigest") ||
		result.Error.NextAction != "edit_configuration" || result.ExitCode() != 2 {
		t.Fatalf("%s", encoded)
	}
}

func TestErrorsMatchTheContract(t *testing.T) {
	contents, err := os.ReadFile(schematest.ContractPath("interfaces.json"))
	if err != nil {
		t.Fatal(err)
	}
	var contract struct {
		Errors []struct {
			Code       string `json:"code"`
			Category   string `json:"category"`
			Exit       int    `json:"exit"`
			Retryable  bool   `json:"retryable"`
			NextAction string `json:"nextAction"`
		} `json:"errors"`
		Warnings []string `json:"warnings"`
	}
	if err := json.Unmarshal(contents, &contract); err != nil {
		t.Fatal(err)
	}
	built := map[string]*Error{"invalid_input": InvalidInput("m"), "operation_timeout": OperationTimeout("m"),
		"operation_cancelled": OperationCancelled("m"), "operation_record_unavailable": OperationRecordUnavailable("m")}
	outcomes := map[string]Outcome{"invalid_input": Failed, "operation_timeout": TimedOut, "operation_cancelled": Cancelled,
		"operation_record_unavailable": Refused}
	for _, want := range contract.Errors {
		got := built[want.Code]
		if got == nil || got.Category != want.Category || got.Retryable != want.Retryable || got.NextAction != want.NextAction {
			t.Fatalf("%s drifted: %+v", want.Code, got)
		}
		result := NewResult(InspectWorkspace, "r", "cli")
		result.Outcome, result.Error = outcomes[want.Code], got
		if err := result.ValidateEnvelope(); err != nil || result.ExitCode() != want.Exit {
			t.Fatalf("%s: exit %d, %v", want.Code, result.ExitCode(), err)
		}
	}
	if len(contract.Warnings) != 1 || contract.Warnings[0] != WarningLogSinkDegraded {
		t.Fatal("warnings drifted")
	}
}

func TestAddWarningOnce(t *testing.T) {
	result := NewResult(InspectWorkspace, "r", "cli")
	result.AddWarning(WarningLogSinkDegraded)
	result.AddWarning(WarningLogSinkDegraded)
	if len(result.Warnings) != 1 {
		t.Fatal(result.Warnings)
	}
}
