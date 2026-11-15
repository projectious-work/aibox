// Package output writes an operation result to stdout as JSON, YAML or a
// human summary (settings-and-logging@v1, output). It writes nothing else to
// stdout; diagnostics go to stderr through the logging package.
package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/projectious-work/aibox/internal/contract"
	"github.com/projectious-work/aibox/internal/redact"
)

// The result formats of output.formats.
const (
	Human = "human"
	JSON  = "json"
	YAML  = "yaml"
)

// Options selects the format of one result.
type Options struct {
	// Format is output.format.
	Format string
	// Color enables ANSI colour; see ColorEnabled.
	Color bool
	// Redactor redacts the human summary, including the workspace root.
	Redactor *redact.Redactor
}

// Write writes the result in the selected format, followed by a newline.
func Write(out io.Writer, result contract.Result, options Options) error {
	encoded, err := json.Marshal(result)
	if err != nil {
		return err
	}
	switch options.Format {
	case JSON:
		_, err = out.Write(append(encoded, '\n'))
	case YAML:
		var text []byte
		text, err = ToYAML(encoded)
		if err == nil {
			_, err = out.Write(text)
		}
	default:
		_, err = io.WriteString(out, options.Redactor.Secrets(humanSummary(result, encoded, options.Color)))
	}
	return err
}

// ColorEnabled is output.colorRule: colour only in human output, only when
// NO_COLOR is unset or empty, and only when color is always, or auto with a
// terminal on stdout (OUT-005).
func ColorEnabled(format, color, noColor string, terminal bool) bool {
	return format == Human && noColor == "" && (color == "always" || color == "auto" && terminal)
}

// BootstrapDiagnostic is the one-line form of output.bootstrapDiagnostic:
// "aibox: invalid_input: <reason>: <text>", with the text redacted and the
// home prefix written as ~ (OUT-007).
func BootstrapDiagnostic(reason, message string, redactor *redact.Redactor) string {
	message = strings.Join(strings.Fields(redactor.All(message)), " ")
	if message == "" {
		message = "the request was rejected"
	}
	return fmt.Sprintf("aibox: invalid_input: %s: %s\n", reason, message)
}

const (
	green  = "\x1b[32m"
	red    = "\x1b[31m"
	yellow = "\x1b[33m"
	reset  = "\x1b[0m"
)

func paint(text, color string, enabled bool) string {
	if !enabled {
		return text
	}
	return color + text + reset
}

// humanSummary is text for a person; its wording is not a contract (D6).
// It always names the operation, the outcome and, on failure, the error code.
func humanSummary(result contract.Result, encoded []byte, color bool) string {
	var b strings.Builder
	outcomeColor := red
	if result.Outcome == contract.Succeeded || result.Outcome == contract.NoChange {
		outcomeColor = green
	}
	fmt.Fprintf(&b, "%s %s\n", result.Operation, paint(string(result.Outcome), outcomeColor, color))
	fmt.Fprintf(&b, "  request: %s\n", result.RequestID)
	if result.Target != nil {
		fmt.Fprintf(&b, "  workspace: %s\n", result.Target.WorkspaceRoot)
	}
	if result.Error != nil {
		fmt.Fprintf(&b, "  error: %s (%s)\n  message: %s\n  next action: %s\n", paint(result.Error.Code, red, color),
			result.Error.Category, result.Error.Message, result.Error.NextAction)
	}
	var envelope struct {
		Data *contract.WorkspaceData `json:"data"`
	}
	if json.Unmarshal(encoded, &envelope) == nil && envelope.Data != nil {
		writeWorkspace(&b, *envelope.Data)
	}
	if len(result.Warnings) > 0 {
		b.WriteString("  warnings:\n")
		for _, warning := range result.Warnings {
			fmt.Fprintf(&b, "    %s\n", paint(warning, yellow, color))
		}
	}
	return b.String()
}

func writeWorkspace(b *strings.Builder, data contract.WorkspaceData) {
	if data.State != "" {
		fmt.Fprintf(b, "  state: %s\n", data.State)
	}
	if len(data.Configuration) > 0 {
		b.WriteString("  configuration:\n")
		for _, entry := range data.Configuration {
			fmt.Fprintf(b, "    %s = %s (%s: %s)\n", entry.Key, humanValue(entry.Value), entry.Layer, entry.Source)
		}
	}
	if len(data.SettingsSources) > 0 {
		b.WriteString("  settings sources:\n")
		for _, source := range data.SettingsSources {
			fmt.Fprintf(b, "    %s: %s (%s)\n", source.Layer, source.Path, source.State)
		}
	}
}

func humanValue(value any) string {
	if value == nil {
		return "unset"
	}
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(value)
	return strings.TrimSuffix(b.String(), "\n")
}
