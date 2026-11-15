// Package logging emits the semantic, redacted log events of
// settings-and-logging@v1 to stderr and to an optional rotating file. Events
// pass a level filter and redaction before they reach a queue or a sink; a
// bounded queue and time limits keep a slow sink from blocking an operation.
package logging

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// SchemaVersion is the const of log-event.schema.json.
const SchemaVersion = "aibox.log-event/v1"

// Executor is logging.members.executor.
const Executor = "aibox"

// MaxLineBytes is logging.fileSink.maxLineBytes.
const MaxLineBytes = 65536

// Severity is an event severity, lowest first.
type Severity int

// The severities of logging.severities.
const (
	Trace Severity = iota
	Debug
	Info
	Warn
	Error
	off
)

var severityNames = []string{"trace", "debug", "info", "warn", "error"}

// String returns the severity as the schema spells it.
func (s Severity) String() string {
	if s < Trace || s > Error {
		return "off"
	}
	return severityNames[s]
}

// threshold maps logging.level to the lowest emitted severity; off emits none.
func threshold(level string) Severity {
	for index, name := range severityNames {
		if name == level {
			return Severity(index)
		}
	}
	return off
}

// Format is logging.format.
type Format string

const (
	// Plain is one human-readable line per event (logging.plainLine).
	Plain Format = "plain"
	// JSONL is one JSON object per line.
	JSONL Format = "jsonl"
)

// Event is one log event. Empty optional members are omitted.
type Event struct {
	Timestamp     time.Time
	Severity      Severity
	Name          string
	Component     string
	Operation     string
	RequestID     string
	Sequence      int
	Initiator     string
	Executor      string
	Outcome       string
	ErrorCategory string
	// Fields are string values (logging.fieldValues).
	Fields map[string]string
}

// member is one optional member in the schema's order.
type member struct{ name, value string }

func (e Event) members() []member {
	return []member{{"operation", e.Operation}, {"initiator", e.Initiator}, {"executor", e.Executor},
		{"outcome", e.Outcome}, {"errorCategory", e.ErrorCategory}}
}

func timestamp(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

func quote(text string) string {
	encoded, _ := json.Marshal(text)
	return string(encoded)
}

func (e Event) fieldNames() []string {
	names := make([]string, 0, len(e.Fields))
	for name := range e.Fields {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// encodeJSON writes one JSON object with a stable member order.
func (e Event) encodeJSON() []byte {
	var b strings.Builder
	fmt.Fprintf(&b, `{"schemaVersion":%s,"timestamp":%s,"severity":%s,"event":%s,"component":%s`,
		quote(SchemaVersion), quote(timestamp(e.Timestamp)), quote(e.Severity.String()), quote(e.Name), quote(e.Component))
	if e.Operation != "" {
		fmt.Fprintf(&b, `,"operation":%s`, quote(e.Operation))
	}
	fmt.Fprintf(&b, `,"requestId":%s,"sequence":%d`, quote(e.RequestID), e.Sequence)
	for _, m := range e.members()[1:] {
		if m.value != "" {
			fmt.Fprintf(&b, `,%s:%s`, quote(m.name), quote(m.value))
		}
	}
	if len(e.Fields) > 0 {
		b.WriteString(`,"fields":{`)
		for index, name := range e.fieldNames() {
			if index > 0 {
				b.WriteByte(',')
			}
			fmt.Fprintf(&b, `%s:%s`, quote(name), quote(e.Fields[name]))
		}
		b.WriteByte('}')
	}
	b.WriteString("}\n")
	return []byte(b.String())
}

// encodePlain writes logging.plainLine: timestamp, severity, event name,
// then member=value pairs and fields as name=value. A value with a space, a
// control character or a quote, or an empty value, is JSON-quoted.
func (e Event) encodePlain() []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s %s", timestamp(e.Timestamp), strings.ToUpper(e.Severity.String()), e.Name)
	pairs := []member{{"component", e.Component}, {"operation", e.Operation}, {"requestId", e.RequestID},
		{"sequence", fmt.Sprint(e.Sequence)}}
	pairs = append(pairs, e.members()[1:]...)
	for _, m := range pairs {
		if m.value != "" {
			fmt.Fprintf(&b, " %s=%s", m.name, plainValue(m.value))
		}
	}
	for _, name := range e.fieldNames() {
		fmt.Fprintf(&b, " %s=%s", name, plainValue(e.Fields[name]))
	}
	b.WriteByte('\n')
	return []byte(b.String())
}

func plainValue(value string) string {
	if value == "" || strings.ContainsFunc(value, func(r rune) bool { return r <= ' ' || r == '"' || r == 0x7f }) {
		return quote(value)
	}
	return value
}

// fit removes fields, last name first, until the JSON line fits
// MaxLineBytes, and then sets fields.truncated (logging.fileSink.longLine).
func (e Event) fit() Event {
	if len(e.encodeJSON()) <= MaxLineBytes {
		return e
	}
	fields := make(map[string]string, len(e.Fields)+1)
	for name, value := range e.Fields {
		fields[name] = value
	}
	e.Fields = fields
	names := e.fieldNames()
	for len(names) > 0 {
		delete(e.Fields, names[len(names)-1])
		names = names[:len(names)-1]
		e.Fields["truncated"] = "true"
		if len(e.encodeJSON()) <= MaxLineBytes {
			break
		}
		delete(e.Fields, "truncated")
	}
	e.Fields["truncated"] = "true"
	return e
}
