package contract

// WarningLogSinkDegraded is the result warning when the file sink failed
// during the operation (LOG-030).
const WarningLogSinkDegraded = "log_sink_degraded"

// ConfigurationEntry is one element of data.configuration of
// inspect_workspace: a native entry or a settings entry with its provenance.
type ConfigurationEntry struct {
	// Key is a native key such as devcontainer.config or a settings key.
	Key string `json:"key"`
	// Value is the effective, redacted value; null means unset.
	Value any `json:"value"`
	// Redacted is true when Value holds the redaction replacement.
	Redacted bool `json:"redacted"`
	// Layer is the winning layer.
	Layer string `json:"layer"`
	// Source is the winning source, redacted.
	Source string `json:"source"`
	// Domain is native, process or workspace.
	Domain string `json:"domain"`
	// InvocationOnly is true only when the winning layer is invocation.
	InvocationOnly bool `json:"invocationOnly"`
	// Overridden lists lower values, highest layer first (CFG-062).
	Overridden []OverriddenValue `json:"overridden,omitempty"`
}

// OverriddenValue is a lower layer's value that lost to a higher layer.
type OverriddenValue struct {
	Layer    string `json:"layer"`
	Source   string `json:"source"`
	Value    any    `json:"value"`
	Redacted bool   `json:"redacted"`
}

// SettingsSource is one element of data.settingsSources (CFG-063).
type SettingsSource struct {
	// Layer is system, user, project or env_file.
	Layer string `json:"layer"`
	// Path is absolute, redacted, with the home prefix as ~.
	Path string `json:"path"`
	// State is loaded or absent.
	State string `json:"state"`
}

// WorkspaceData is the closed inspect_workspace payload. SettingsSources is
// present only in the effective-configuration view (CFG-067).
type WorkspaceData struct {
	State           string               `json:"state"`
	Resources       []any                `json:"resources"`
	Configuration   []ConfigurationEntry `json:"configuration"`
	ObservedAt      string               `json:"observedAt"`
	SettingsSources []SettingsSource     `json:"settingsSources,omitempty"`
}

// The structured errors of this version (interfaces.json, errors).

// InvalidInput is a defect in the request, the settings or the project.
func InvalidInput(message string) *Error {
	return &Error{Code: "invalid_input", Category: "invalid_input", Message: message, Retryable: false, NextAction: "edit_configuration"}
}

// OperationTimeout reports an operation that exceeded its time limit.
func OperationTimeout(message string) *Error {
	return &Error{Code: "operation_timeout", Category: "timeout", Message: message, Retryable: true, NextAction: "retry_read"}
}

// OperationCancelled reports an operation interrupted by SIGINT or SIGTERM.
func OperationCancelled(message string) *Error {
	return &Error{Code: "operation_cancelled", Category: "interrupted", Message: message, Retryable: true, NextAction: "retry_read"}
}

// OperationRecordUnavailable refuses a mutating operation while a required
// file sink has failed (LOG-031, D16).
func OperationRecordUnavailable(message string) *Error {
	return &Error{Code: "operation_record_unavailable", Category: "internal", Message: message, Retryable: false, NextAction: "manual_recovery"}
}

// NewResult returns an envelope with the common lists present and empty.
func NewResult(operation Operation, requestID, initiator string) Result {
	return Result{
		SchemaVersion: ResultSchemaVersion, Operation: operation, RequestID: requestID,
		Actors:           Actors{Initiator: initiator, Executor: "aibox"},
		ChangedResources: []string{}, Warnings: []string{}, Evidence: []string{},
	}
}

// InvalidRequestResult is the envelope of a failure before the operation
// starts: operation invalid_request, no target and no input digest (OUT-006).
func InvalidRequestResult(requestID, initiator, message string) Result {
	result := NewResult(InvalidRequest, requestID, initiator)
	result.Outcome = Failed
	result.Error = InvalidInput(message)
	return result
}

// AddWarning appends a warning once.
func (r *Result) AddWarning(warning string) {
	for _, existing := range r.Warnings {
		if existing == warning {
			return
		}
	}
	r.Warnings = append(r.Warnings, warning)
}
