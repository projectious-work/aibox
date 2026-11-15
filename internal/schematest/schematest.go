// Package schematest validates JSON documents against the frozen contract
// schemas in tests. It implements only the JSON Schema 2020-12 keywords that
// the settings-and-logging@v1 schemas use and fails on any other keyword, so
// that a schema change cannot pass unnoticed. Product code does not use it:
// documents are validated in tests, not at run time.
package schematest

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"
)

// Schema is one loaded schema document.
type Schema struct {
	root map[string]any
}

// ContractPath returns the path of a file of the frozen contract.
func ContractPath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "spec", "settings-and-logging", "v1", name)
}

// Load reads a schema of the frozen contract, for example
// "schemas/log-event.schema.json".
func Load(name string) (*Schema, error) {
	contents, err := os.ReadFile(ContractPath(name))
	if err != nil {
		return nil, err
	}
	var root map[string]any
	if err := json.Unmarshal(contents, &root); err != nil {
		return nil, err
	}
	return &Schema{root: root}, nil
}

// Validate checks one JSON document.
func (s *Schema) Validate(document []byte) error {
	var value any
	decoder := json.NewDecoder(strings.NewReader(string(document)))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return err
	}
	return s.check(s.root, value, "$")
}

var ignored = map[string]bool{"$schema": true, "$id": true, "title": true, "description": true, "$defs": true, "$comment": true}

func (s *Schema) check(schema map[string]any, value any, at string) error {
	for keyword := range schema {
		if ignored[keyword] {
			continue
		}
		if err := s.keyword(schema, keyword, value, at); err != nil {
			return err
		}
	}
	return nil
}

func (s *Schema) keyword(schema map[string]any, keyword string, value any, at string) error {
	rule := schema[keyword]
	switch keyword {
	case "type":
		return checkType(rule, value, at)
	case "enum":
		for _, allowed := range rule.([]any) {
			if equal(allowed, value) {
				return nil
			}
		}
		return fmt.Errorf("%s: %v is not in the enum", at, value)
	case "const":
		if !equal(rule, value) {
			return fmt.Errorf("%s: %v is not %v", at, value, rule)
		}
	case "required":
		object, ok := value.(map[string]any)
		if !ok {
			return nil
		}
		for _, name := range rule.([]any) {
			if _, present := object[name.(string)]; !present {
				return fmt.Errorf("%s: %s is required", at, name)
			}
		}
	case "properties", "additionalProperties":
		return s.checkProperties(schema, value, at)
	case "items":
		if array, ok := value.([]any); ok {
			for index, item := range array {
				if err := s.check(rule.(map[string]any), item, fmt.Sprintf("%s[%d]", at, index)); err != nil {
					return err
				}
			}
		}
	case "maxItems", "minItems":
		if array, ok := value.([]any); ok && !within(keyword, float64(len(array)), rule) {
			return fmt.Errorf("%s: %d items violate %s", at, len(array), keyword)
		}
	case "minLength", "maxLength":
		if text, ok := value.(string); ok && !within(keyword, float64(utf8.RuneCountInString(text)), rule) {
			return fmt.Errorf("%s: length violates %s", at, keyword)
		}
	case "minimum", "maximum":
		if number, ok := value.(json.Number); ok {
			parsed, _ := number.Float64()
			if !within(keyword, parsed, rule) {
				return fmt.Errorf("%s: %v violates %s", at, number, keyword)
			}
		}
	case "pattern":
		if text, ok := value.(string); ok && !regexp.MustCompile(rule.(string)).MatchString(text) {
			return fmt.Errorf("%s: %q does not match %s", at, text, rule)
		}
	case "format":
		if text, ok := value.(string); ok && rule == "date-time" {
			if _, err := time.Parse(time.RFC3339Nano, text); err != nil {
				return fmt.Errorf("%s: %q is not a date-time", at, text)
			}
		}
	case "allOf":
		for _, sub := range rule.([]any) {
			if err := s.check(sub.(map[string]any), value, at); err != nil {
				return err
			}
		}
	case "if":
		if s.check(rule.(map[string]any), value, at) == nil {
			if then, ok := schema["then"].(map[string]any); ok {
				return s.check(then, value, at)
			}
		} else if otherwise, ok := schema["else"].(map[string]any); ok {
			return s.check(otherwise, value, at)
		}
	case "then", "else":
	case "not":
		if s.check(rule.(map[string]any), value, at) == nil {
			return fmt.Errorf("%s: matches a forbidden schema", at)
		}
	case "$ref":
		return s.check(s.resolve(rule.(string)), value, at)
	default:
		return fmt.Errorf("schematest: unsupported keyword %q at %s", keyword, at)
	}
	return nil
}

func (s *Schema) checkProperties(schema map[string]any, value any, at string) error {
	object, ok := value.(map[string]any)
	if !ok {
		return nil
	}
	properties, _ := schema["properties"].(map[string]any)
	for name, member := range object {
		if sub, declared := properties[name]; declared {
			if err := s.check(sub.(map[string]any), member, at+"."+name); err != nil {
				return err
			}
			continue
		}
		switch additional := schema["additionalProperties"].(type) {
		case bool:
			if !additional {
				return fmt.Errorf("%s: unexpected member %q", at, name)
			}
		case map[string]any:
			if err := s.check(additional, member, at+"."+name); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Schema) resolve(reference string) map[string]any {
	node := any(s.root)
	for _, part := range strings.Split(strings.TrimPrefix(reference, "#/"), "/") {
		node = node.(map[string]any)[part]
	}
	return node.(map[string]any)
}

func checkType(rule, value any, at string) error {
	types, ok := rule.([]any)
	if !ok {
		types = []any{rule}
	}
	for _, name := range types {
		if hasType(name.(string), value) {
			return nil
		}
	}
	return fmt.Errorf("%s: %v is not of type %v", at, value, rule)
}

func hasType(name string, value any) bool {
	switch name {
	case "object":
		_, ok := value.(map[string]any)
		return ok
	case "array":
		_, ok := value.([]any)
		return ok
	case "string":
		_, ok := value.(string)
		return ok
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "null":
		return value == nil
	case "number":
		_, ok := value.(json.Number)
		return ok
	case "integer":
		number, ok := value.(json.Number)
		if !ok {
			return false
		}
		parsed, err := number.Float64()
		return err == nil && parsed == math.Trunc(parsed)
	}
	return false
}

func within(keyword string, actual float64, rule any) bool {
	limit, _ := rule.(float64)
	if strings.HasPrefix(keyword, "min") {
		return actual >= limit
	}
	return actual <= limit
}

func equal(expected, actual any) bool {
	if number, ok := actual.(json.Number); ok {
		parsed, _ := number.Float64()
		return expected == parsed
	}
	left, _ := json.Marshal(expected)
	right, _ := json.Marshal(actual)
	return string(left) == string(right)
}
