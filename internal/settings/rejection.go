package settings

import "fmt"

// Reason is a rejection reason of settings.rejectionReasons. Every rejection
// stops the command before its operation starts with error code
// invalid_input and exit 2.
type Reason string

// The closed list of rejection reasons.
const (
	Malformed              Reason = "malformed"
	Unreadable             Reason = "unreadable"
	MissingSelectedFile    Reason = "missing_selected_file"
	NotRegularFile         Reason = "not_regular_file"
	DuplicateKey           Reason = "duplicate_key"
	UnknownKey             Reason = "unknown_key"
	UnsupportedVersion     Reason = "unsupported_version"
	NullValue              Reason = "null_value"
	InvalidValue           Reason = "invalid_value"
	UnauthorizedSource     Reason = "unauthorized_source"
	UnknownEnvironmentName Reason = "unknown_environment_name"
	RelativeEnvFile        Reason = "relative_env_file"
	PathEscape             Reason = "path_escape"
	NoHome                 Reason = "no_home"
	UsageError             Reason = "usage_error"
)

// Reasons lists every rejection reason in the order of interfaces.json.
var Reasons = []Reason{Malformed, Unreadable, MissingSelectedFile, NotRegularFile, DuplicateKey, UnknownKey,
	UnsupportedVersion, NullValue, InvalidValue, UnauthorizedSource, UnknownEnvironmentName, RelativeEnvFile,
	PathEscape, NoHome, UsageError}

// Rejection is the first defect found in the order of settings.failureOrder.
// Source is unredacted; callers redact it before it leaves the process.
type Rejection struct {
	// Reason is the rejection reason.
	Reason Reason
	// Key is the dotted settings key, when one is known.
	Key string
	// Layer is the layer that holds the defect, when one is known.
	Layer Layer
	// Source is the file path, environment name or flag, when known.
	Source string
	// Suggestion is the nearest env name, only for unknown_environment_name.
	Suggestion string
	// Message is a short explanation for the bootstrap diagnostic.
	Message string
}

// Error makes a Rejection usable as an error value.
func (r *Rejection) Error() string {
	return fmt.Sprintf("%s: %s", r.Reason, r.Message)
}

// Usage returns a usage_error rejection for a defective command line.
func Usage(message string) *Rejection {
	return &Rejection{Reason: UsageError, Message: message}
}

func reject(reason Reason, layer Layer, source, key, message string) *Rejection {
	return &Rejection{Reason: reason, Key: key, Layer: layer, Source: source, Message: message}
}
