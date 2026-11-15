package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// ToYAML converts one JSON document to block-style YAML 1.2 by output.yaml:
// every string and mapping key is a double-quoted scalar with JSON escapes;
// numbers, true, false and null are plain; empty collections are [] and {}.
// Member order follows the JSON document. No YAML library is used (D7).
func ToYAML(document []byte) ([]byte, error) {
	decoder := json.NewDecoder(bytes.NewReader(document))
	decoder.UseNumber()
	value, err := readValue(decoder)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	switch node := value.(type) {
	case object:
		if len(node) == 0 {
			b.WriteString("{}\n")
		}
		writeMapping(&b, node, 0)
	case []any:
		if len(node) == 0 {
			b.WriteString("[]\n")
		}
		writeSequence(&b, node, 0)
	default:
		b.WriteString(scalar(node) + "\n")
	}
	return []byte(b.String()), nil
}

// object keeps the members of a JSON object in document order.
type object []pair

type pair struct {
	name  string
	value any
}

func readValue(decoder *json.Decoder) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return token, nil
	}
	switch delim {
	case '{':
		members := object{}
		for decoder.More() {
			name, err := decoder.Token()
			if err != nil {
				return nil, err
			}
			value, err := readValue(decoder)
			if err != nil {
				return nil, err
			}
			members = append(members, pair{name.(string), value})
		}
		_, err = decoder.Token()
		return members, err
	case '[':
		items := []any{}
		for decoder.More() {
			value, err := readValue(decoder)
			if err != nil {
				return nil, err
			}
			items = append(items, value)
		}
		_, err = decoder.Token()
		return items, err
	}
	return nil, errors.New("unexpected JSON delimiter")
}

func indent(level int) string { return strings.Repeat("  ", level) }

// inline returns the text of a value that fits after "key: " or "- ".
func inline(value any) (string, bool) {
	switch node := value.(type) {
	case object:
		return "{}", len(node) == 0
	case []any:
		return "[]", len(node) == 0
	default:
		return scalar(node), true
	}
}

func writeMapping(b *strings.Builder, members object, level int) {
	for _, member := range members {
		writeMember(b, member, indent(level), level)
	}
}

func writeMember(b *strings.Builder, member pair, prefix string, level int) {
	b.WriteString(prefix + quoteString(member.name) + ":")
	if text, ok := inline(member.value); ok {
		b.WriteString(" " + text + "\n")
		return
	}
	b.WriteString("\n")
	writeBlock(b, member.value, level+1)
}

func writeBlock(b *strings.Builder, value any, level int) {
	if members, ok := value.(object); ok {
		writeMapping(b, members, level)
		return
	}
	writeSequence(b, value.([]any), level)
}

// writeSequence writes "- item" lines. A mapping item starts on the dash
// line; a nested sequence item starts on the next line, indented.
func writeSequence(b *strings.Builder, items []any, level int) {
	for _, item := range items {
		if text, ok := inline(item); ok {
			b.WriteString(indent(level) + "- " + text + "\n")
			continue
		}
		if members, ok := item.(object); ok {
			for index, member := range members {
				prefix := indent(level + 1)
				if index == 0 {
					prefix = indent(level) + "- "
				}
				writeMember(b, member, prefix, level+1)
			}
			continue
		}
		b.WriteString(indent(level) + "-\n")
		writeSequence(b, item.([]any), level+1)
	}
}

func scalar(value any) string {
	switch node := value.(type) {
	case nil:
		return "null"
	case bool:
		if node {
			return "true"
		}
		return "false"
	case json.Number:
		return node.String()
	case string:
		return quoteString(node)
	}
	return quoteString(fmt.Sprint(value))
}

// quoteString writes a YAML double-quoted scalar using only escapes that
// JSON and YAML 1.2 share. Characters that YAML does not allow unescaped
// (C0 and C1 controls, DEL, line and paragraph separators, BOM) use \uXXXX.
func quoteString(text string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range text {
		switch {
		case r == '"':
			b.WriteString(`\"`)
		case r == '\\':
			b.WriteString(`\\`)
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
			b.WriteString(`\t`)
		case r < 0x20 || r == 0x7f || r >= 0x80 && r <= 0x9f || r == 0x2028 || r == 0x2029 || r == 0xfeff || r == utf8.RuneError:
			fmt.Fprintf(&b, `\u%04x`, r)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}
