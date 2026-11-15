package settings

import (
	"regexp"
	"slices"
	"strconv"
)

var integerText = regexp.MustCompile(`^-?[0-9]+$`)

// convertText converts an environment, env-file or flag value by
// settings.conversion. It returns false for every value that the key's type,
// values or range do not admit, including the empty value.
func convertText(key Key, text string) (any, bool) {
	if text == "" {
		return nil, false
	}
	switch key.Type {
	case Boolean:
		switch text {
		case "true":
			return true, true
		case "false":
			return false, true
		}
		return nil, false
	case Integer:
		if !integerText.MatchString(text) {
			return nil, false
		}
		number, err := strconv.Atoi(text)
		if err != nil || number < key.Minimum || number > key.Maximum {
			return nil, false
		}
		return number, true
	case Enum:
		return text, slices.Contains(key.Values, text)
	case Path:
		return text, true
	}
	return nil, false
}

// convertJSON converts a settings-file value. Null is checked by the caller.
func convertJSON(key Key, node *jsonNode) (any, bool) {
	switch key.Type {
	case Boolean:
		return node.boolean, node.kind == jsonBool
	case Integer:
		if node.kind != jsonNumber {
			return nil, false
		}
		number, err := strconv.ParseInt(string(node.number), 10, 64)
		if err != nil {
			value, floatErr := strconv.ParseFloat(string(node.number), 64)
			if floatErr != nil || value != float64(int64(value)) {
				return nil, false
			}
			number = int64(value)
		}
		if number < int64(key.Minimum) || number > int64(key.Maximum) {
			return nil, false
		}
		return int(number), true
	case Enum, Path:
		if node.kind != jsonString {
			return nil, false
		}
		return convertText(key, node.text)
	}
	return nil, false
}

// modeAdmits applies the per-mode value rule: mcp serve writes events only
// as JSON Lines, so logging.format plain from any layer stops it (CFG-054).
func modeAdmits(mode Mode, key Key, value any) bool {
	return !(mode == ModeMCP && key.Name == LoggingFormat && value == "plain")
}
