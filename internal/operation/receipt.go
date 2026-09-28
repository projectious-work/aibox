package operation

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/projectious-work/aibox/internal/contract"
)

// ReceiptSchemaVersion identifies the durable operation record contract.
const ReceiptSchemaVersion = "aibox.operation-receipt/v1"

const maxReceiptBytes = 1024 * 1024

var receiptIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
var receiptDigestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

// Resource records only a bounded native identity and observed state. It
// never stores a runtime credential, environment value or child transcript.
type Resource struct {
	Kind  string `json:"kind"`
	ID    string `json:"id"`
	Role  string `json:"role"`
	State string `json:"state"`
	Owned bool   `json:"owned"`
}

// Receipt is a non-secret, schema-shaped record of one operation's progress.
// It is not desired state and cannot authorize replay of an incomplete call.
type Receipt struct {
	SchemaVersion      string             `json:"schemaVersion"`
	OperationID        string             `json:"operationId"`
	RequestFingerprint string             `json:"requestFingerprint"`
	Operation          contract.Operation `json:"operation"`
	ProjectRoot        string             `json:"projectRoot"`
	ConfigPath         string             `json:"configPath,omitempty"`
	InputDigest        string             `json:"inputDigest"`
	PolicyDigest       string             `json:"policyDigest,omitempty"`
	GrantID            string             `json:"grantId,omitempty"`
	CreatedAt          string             `json:"createdAt"`
	UpdatedAt          string             `json:"updatedAt"`
	Executor           string             `json:"executor"`
	State              string             `json:"state"`
	LastConfirmedStep  string             `json:"lastConfirmedStep"`
	Resources          []Resource         `json:"resources"`
	CompletedEffects   []contract.Effect  `json:"completedEffects"`
	UnknownEffects     []contract.Effect  `json:"unknownEffects"`
	Result             *contract.Result   `json:"result,omitempty"`
}

// Validate checks the closed receipt envelope before a record can be written
// or trusted after reading. Full semantic authorization remains the caller's
// responsibility and never follows merely from a valid receipt.
func (r Receipt) Validate() error {
	if r.SchemaVersion != ReceiptSchemaVersion || !receiptIDPattern.MatchString(r.OperationID) ||
		!receiptDigestPattern.MatchString(r.RequestFingerprint) || !receiptDigestPattern.MatchString(r.InputDigest) {
		return errors.New("invalid receipt identity or digest")
	}
	if r.Operation == contract.InvalidRequest || !validReceiptOperation(r.Operation) {
		return errors.New("invalid receipt operation")
	}
	if !filepath.IsAbs(r.ProjectRoot) || len(r.ProjectRoot) > 4096 || filepath.Clean(r.ProjectRoot) != r.ProjectRoot {
		return errors.New("invalid receipt project root")
	}
	if r.ConfigPath != "" && (!filepath.IsAbs(r.ConfigPath) || len(r.ConfigPath) > 4096) {
		return errors.New("invalid receipt config path")
	}
	if r.PolicyDigest != "" && !receiptDigestPattern.MatchString(r.PolicyDigest) {
		return errors.New("invalid receipt policy digest")
	}
	if r.GrantID != "" && !receiptIDPattern.MatchString(r.GrantID) {
		return errors.New("invalid receipt grant ID")
	}
	created, err := time.Parse(time.RFC3339Nano, r.CreatedAt)
	if err != nil {
		return errors.New("invalid receipt creation time")
	}
	updated, err := time.Parse(time.RFC3339Nano, r.UpdatedAt)
	if err != nil || updated.Before(created) {
		return errors.New("invalid receipt update time")
	}
	if r.Executor == "" || len(r.Executor) > 4096 || r.LastConfirmedStep == "" || len(r.LastConfirmedStep) > 4096 || !validReceiptState(r.State) {
		return errors.New("invalid receipt state or attribution")
	}
	if r.Resources == nil || r.CompletedEffects == nil || r.UnknownEffects == nil || len(r.Resources) > 64 || len(r.CompletedEffects) > 1000 || len(r.UnknownEffects) > 1000 {
		return errors.New("invalid receipt lists")
	}
	for _, resource := range r.Resources {
		if !oneOf(resource.Kind, "container", "image", "network", "volume", "file") ||
			!oneOf(resource.Role, "primary", "sidecar", "shared", "artifact") ||
			!oneOf(resource.State, "running", "stopped", "absent", "unknown", "built") ||
			resource.ID == "" || len(resource.ID) > 4096 {
			return errors.New("invalid receipt resource")
		}
	}
	for _, effect := range append(append([]contract.Effect{}, r.CompletedEffects...), r.UnknownEffects...) {
		if effect.Resource == "" || len(effect.Resource) > 4096 || effect.Action == "" || len(effect.Action) > 4096 || !oneOf(effect.Status, "confirmed", "unknown") {
			return errors.New("invalid receipt effect")
		}
	}
	if r.Result != nil {
		if err := r.Result.ValidateEnvelope(); err != nil {
			return fmt.Errorf("invalid receipt result: %w", err)
		}
		if r.Result.RequestID != r.OperationID || r.Result.Operation != r.Operation {
			return errors.New("receipt result identity mismatch")
		}
	}
	return nil
}

func validReceiptOperation(value contract.Operation) bool {
	return oneOf(string(value), "build_environment", "start_environment", "stop_environment", "remove_environment", "rebuild_environment", "inspect_environment", "check_environment", "inspect_workspace", "check_workspace", "refresh_workspace", "read_logs", "inspect_operation", "migrate_preview", "migrate_apply", "migrate_rollback")
}

func validReceiptState(value string) bool {
	return oneOf(value, "validated", "authorized", "executing", "inspected", "succeeded", "no_change", "failed", "cancelled", "partial", "timed_out", "refused", "rebuild_required")
}

func oneOf(value string, choices ...string) bool {
	for _, choice := range choices {
		if value == choice {
			return true
		}
	}
	return false
}

// marshalReceipt enforces a bounded, valid JSON record before any disk write.
func marshalReceipt(receipt Receipt) ([]byte, error) {
	if err := receipt.Validate(); err != nil {
		return nil, err
	}
	data, err := json.Marshal(receipt)
	if err != nil {
		return nil, err
	}
	if len(data) > maxReceiptBytes {
		return nil, errors.New("receipt exceeds one MiB")
	}
	return data, nil
}
