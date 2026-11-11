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

// OperationRecordSchemaVersion identifies the durable operation record contract.
const OperationRecordSchemaVersion = "aibox.operation-record/v1"

const maxOperationRecordBytes = 1024 * 1024

var recordIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
var recordDigestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

// Resource records only a bounded native identity and observed state. It
// never stores a runtime credential, environment value or child transcript.
type Resource struct {
	Kind  string `json:"kind"`
	ID    string `json:"id"`
	Role  string `json:"role"`
	State string `json:"state"`
	Owned bool   `json:"owned"`
}

// OperationRecord is a non-secret, schema-shaped record of one operation's progress.
// It is not desired state and cannot authorize replay of an incomplete call.
type OperationRecord struct {
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

// Validate checks the closed record envelope before a record can be written
// or trusted after reading. Full semantic authorization remains the caller's
// responsibility and never follows merely from a valid record.
func (r OperationRecord) Validate() error {
	if r.SchemaVersion != OperationRecordSchemaVersion || !recordIDPattern.MatchString(r.OperationID) ||
		!recordDigestPattern.MatchString(r.RequestFingerprint) || !recordDigestPattern.MatchString(r.InputDigest) {
		return errors.New("invalid record identity or digest")
	}
	if r.Operation == contract.InvalidRequest || !validOperationRecordOperation(r.Operation) {
		return errors.New("invalid record operation")
	}
	if !filepath.IsAbs(r.ProjectRoot) || len(r.ProjectRoot) > 4096 || filepath.Clean(r.ProjectRoot) != r.ProjectRoot {
		return errors.New("invalid record project root")
	}
	if r.ConfigPath != "" && (!filepath.IsAbs(r.ConfigPath) || len(r.ConfigPath) > 4096) {
		return errors.New("invalid record config path")
	}
	if r.PolicyDigest != "" && !recordDigestPattern.MatchString(r.PolicyDigest) {
		return errors.New("invalid record policy digest")
	}
	if r.GrantID != "" && !recordIDPattern.MatchString(r.GrantID) {
		return errors.New("invalid record grant ID")
	}
	created, err := time.Parse(time.RFC3339Nano, r.CreatedAt)
	if err != nil {
		return errors.New("invalid record creation time")
	}
	updated, err := time.Parse(time.RFC3339Nano, r.UpdatedAt)
	if err != nil || updated.Before(created) {
		return errors.New("invalid record update time")
	}
	if r.Executor == "" || len(r.Executor) > 4096 || r.LastConfirmedStep == "" || len(r.LastConfirmedStep) > 4096 || !validOperationRecordState(r.State) {
		return errors.New("invalid record state or attribution")
	}
	if r.Resources == nil || r.CompletedEffects == nil || r.UnknownEffects == nil || len(r.Resources) > 64 || len(r.CompletedEffects) > 1000 || len(r.UnknownEffects) > 1000 {
		return errors.New("invalid record lists")
	}
	for _, resource := range r.Resources {
		if !oneOf(resource.Kind, "container", "image", "network", "volume", "file") ||
			!oneOf(resource.Role, "primary", "sidecar", "shared", "artifact") ||
			!oneOf(resource.State, "running", "stopped", "absent", "unknown", "built") ||
			resource.ID == "" || len(resource.ID) > 4096 {
			return errors.New("invalid record resource")
		}
	}
	for _, effect := range append(append([]contract.Effect{}, r.CompletedEffects...), r.UnknownEffects...) {
		if effect.Resource == "" || len(effect.Resource) > 4096 || effect.Action == "" || len(effect.Action) > 4096 || !oneOf(effect.Status, "confirmed", "unknown") {
			return errors.New("invalid record effect")
		}
	}
	if r.Result != nil {
		if err := r.Result.ValidateEnvelope(); err != nil {
			return fmt.Errorf("invalid record result: %w", err)
		}
		if r.Result.RequestID != r.OperationID || r.Result.Operation != r.Operation {
			return errors.New("record result identity mismatch")
		}
	}
	return nil
}

func validOperationRecordOperation(value contract.Operation) bool {
	return oneOf(string(value), "build_environment", "start_environment", "stop_environment", "remove_environment", "rebuild_environment", "inspect_environment", "check_environment", "inspect_workspace", "check_workspace", "refresh_workspace", "read_logs", "inspect_operation", "migrate_preview", "migrate_apply", "migrate_rollback")
}

func validOperationRecordState(value string) bool {
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

// marshalOperationRecord enforces a bounded, valid JSON record before any disk write.
func marshalOperationRecord(record OperationRecord) ([]byte, error) {
	if err := record.Validate(); err != nil {
		return nil, err
	}
	data, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	if len(data) > maxOperationRecordBytes {
		return nil, errors.New("record exceeds one MiB")
	}
	return data, nil
}
