package settings

import (
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// Assignment is one value that a layer gives a key.
type Assignment struct {
	// Key is the dotted key name.
	Key string
	// Layer is the layer that set the value.
	Layer Layer
	// Source is the unredacted source as effectiveConfig.source defines.
	Source string
	// Value is the converted value; a logging.file value is the confined,
	// symlink-resolved absolute path.
	Value any
}

// SourceState is one element of settingsSources, with an unredacted path.
type SourceState struct {
	// Layer is system, user, project or env_file.
	Layer Layer
	// Path is the absolute, symlink-resolved path when the file exists.
	Path string
	// State is loaded or absent.
	State string
}

// layerResult is what one source contributed: its assignments in document
// order, or the first defect it holds.
type layerResult struct {
	layer       Layer
	source      *SourceState
	assignments []Assignment
	rejection   *Rejection
}

// valueChecker applies the checks that every layer shares after conversion.
type valueChecker struct {
	mode  Mode
	paths Paths
}

// check applies, for one setting, the shape checks before the authority
// check (settings.failureOrder), then confines a log file path (LOG-014).
func (c valueChecker) check(key Key, layer Layer, source string, value any, ok bool) (Assignment, *Rejection) {
	if !ok || !modeAdmits(c.mode, key, value) {
		return Assignment{}, reject(InvalidValue, layer, source, key.Name, fmt.Sprintf("%s has a value that is not allowed", key.Name))
	}
	if !key.Allows(layer) {
		return Assignment{}, reject(UnauthorizedSource, layer, source, key.Name, fmt.Sprintf("%s cannot be set by the %s layer", key.Name, layer))
	}
	if key.Name == LoggingFile {
		confined, err := c.paths.ConfineLogFile(c.paths.AnchorLogFile(value.(string)))
		if err != nil {
			return Assignment{}, reject(PathEscape, layer, source, key.Name, "logging.file is outside the log directory")
		}
		value = confined
	}
	return Assignment{Key: key.Name, Layer: layer, Source: source, Value: value}, nil
}

// loadSettingsFile reads a system or user settings file (CFG-010 to CFG-018).
func loadSettingsFile(fsys FileSystem, check valueChecker, layer Layer, path string, selected bool) layerResult {
	result := layerResult{layer: layer, source: &SourceState{Layer: layer, Path: path, State: "absent"}}
	resolved, contents, err := fsys.ReadFile(path, maxSettingsFileBytes)
	switch {
	case errors.Is(err, fs.ErrNotExist) && selected:
		result.rejection = reject(MissingSelectedFile, layer, path, "", "the selected settings file does not exist")
		return result
	case errors.Is(err, fs.ErrNotExist):
		return result
	case err != nil:
		result.rejection = reject(Unreadable, layer, path, "", "the settings file cannot be read")
		return result
	}
	result.source = &SourceState{Layer: layer, Path: resolved, State: "loaded"}
	document, err := parseJSON(contents)
	if err != nil || document.kind != jsonObject {
		result.rejection = reject(Malformed, layer, resolved, "", "the settings file is not a JSON object")
		return result
	}
	result.assignments, result.rejection = walkSettingsDocument(check, layer, resolved, document)
	return result
}

// walkSettingsDocument checks the document against the process-settings
// schema in document order and maps each defect to settings.schemaDefects.
func walkSettingsDocument(check valueChecker, layer Layer, source string, document *jsonNode) ([]Assignment, *Rejection) {
	var assignments []Assignment
	seen := map[string]bool{}
	sawVersion := false
	for _, member := range document.members {
		if seen[member.name] {
			return nil, reject(DuplicateKey, layer, source, "", fmt.Sprintf("member %q occurs twice", member.name))
		}
		seen[member.name] = true
		if member.name == "schemaVersion" {
			if member.value.kind == jsonNull {
				return nil, reject(NullValue, layer, source, "", "schemaVersion is null")
			}
			if member.value.kind != jsonString || member.value.text != "1" {
				return nil, reject(UnsupportedVersion, layer, source, "", "schemaVersion is not \"1\"")
			}
			sawVersion = true
			continue
		}
		if !isGroup(member.name) {
			return nil, reject(UnknownKey, layer, source, "", fmt.Sprintf("member %q is not a setting", member.name))
		}
		groupAssignments, rejection := walkSettingsGroup(check, layer, source, member)
		if rejection != nil {
			return nil, rejection
		}
		assignments = append(assignments, groupAssignments...)
	}
	if !sawVersion {
		return nil, reject(UnsupportedVersion, layer, source, "", "schemaVersion is missing")
	}
	return assignments, nil
}

func walkSettingsGroup(check valueChecker, layer Layer, source string, group jsonMember) ([]Assignment, *Rejection) {
	if group.value.kind == jsonNull {
		return nil, reject(NullValue, layer, source, "", fmt.Sprintf("member %q is null", group.name))
	}
	if group.value.kind != jsonObject {
		return nil, reject(InvalidValue, layer, source, "", fmt.Sprintf("member %q is not an object", group.name))
	}
	var assignments []Assignment
	seen := map[string]bool{}
	for _, member := range group.value.members {
		name := group.name + "." + member.name
		if seen[member.name] {
			return nil, reject(DuplicateKey, layer, source, name, fmt.Sprintf("%s occurs twice", name))
		}
		seen[member.name] = true
		key, known := Lookup(name)
		if !known {
			return nil, reject(UnknownKey, layer, source, "", fmt.Sprintf("%q is not a setting", name))
		}
		if member.value.kind == jsonNull {
			return nil, reject(NullValue, layer, source, name, fmt.Sprintf("%s is null", name))
		}
		value, ok := convertJSON(key, member.value)
		assignment, rejection := check.check(key, layer, source, value, ok)
		if rejection != nil {
			return nil, rejection
		}
		assignments = append(assignments, assignment)
	}
	return assignments, nil
}

func isGroup(name string) bool {
	for _, key := range Keys {
		if strings.HasPrefix(key.Name, name+".") {
			return true
		}
	}
	return false
}

var envKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// loadEnvFile reads one environment file by settings.envFile.grammar.
func loadEnvFile(fsys FileSystem, check valueChecker, path string) layerResult {
	result := layerResult{layer: EnvFile, source: &SourceState{Layer: EnvFile, Path: path, State: "absent"}}
	resolved, contents, err := fsys.ReadFile(path, maxEnvFileBytes)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		result.rejection = reject(MissingSelectedFile, EnvFile, path, "", "the environment file does not exist")
		return result
	case errors.Is(err, ErrNotRegular):
		result.rejection = reject(NotRegularFile, EnvFile, resolved, "", "the environment file is not a regular file")
		return result
	case err != nil:
		result.rejection = reject(Unreadable, EnvFile, path, "", "the environment file cannot be read")
		return result
	}
	result.source = &SourceState{Layer: EnvFile, Path: resolved, State: "loaded"}
	if !utf8.Valid(contents) {
		result.rejection = reject(Malformed, EnvFile, resolved, "", "the environment file is not UTF-8 text")
		return result
	}
	seen := map[string]bool{}
	for number, line := range strings.Split(string(contents), "\n") {
		line = strings.TrimSuffix(line, "\r")
		if strings.Trim(line, " \t") == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name, value, found := strings.Cut(line, "=")
		if !found || !envKeyPattern.MatchString(name) {
			result.rejection = reject(Malformed, EnvFile, resolved, "", fmt.Sprintf("line %d is not KEY=VALUE", number+1))
			return result
		}
		if seen[name] {
			result.rejection = reject(DuplicateKey, EnvFile, resolved, "", fmt.Sprintf("%s is set twice", name))
			return result
		}
		seen[name] = true
		key, known := lookupEnv(name)
		if !known {
			result.rejection = reject(UnknownKey, EnvFile, resolved, "", fmt.Sprintf("%s is not a setting", name))
			return result
		}
		converted, ok := convertText(key, value)
		assignment, rejection := check.check(key, EnvFile, resolved, converted, ok)
		if rejection != nil {
			result.rejection = rejection
			return result
		}
		result.assignments = append(result.assignments, assignment)
	}
	return result
}

// loadEnvironment maps AIBOX_* variables by their env names, in byte order
// of the names (CFG-039 to CFG-042).
func loadEnvironment(check valueChecker, env map[string]string) layerResult {
	result := layerResult{layer: Environment}
	names := make([]string, 0, len(env))
	for name := range env {
		if strings.HasPrefix(name, environmentNamePrefix) {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		if name == EnvFileVariable {
			continue
		}
		key, known := lookupEnv(name)
		if !known {
			if registered(name) {
				continue
			}
			result.rejection = reject(UnknownEnvironmentName, Environment, name, "", fmt.Sprintf("%s is not a known environment name", name))
			result.rejection.Suggestion = nearestName(name)
			return result
		}
		converted, ok := convertText(key, env[name])
		assignment, rejection := check.check(key, Environment, name, converted, ok)
		if rejection != nil {
			result.rejection = rejection
			return result
		}
		result.assignments = append(result.assignments, assignment)
	}
	return result
}

// Flag is one settings flag of the command line, in command-line order.
type Flag struct {
	// Name is the flag, for example --log-level.
	Name string
	// Value is the raw value as given.
	Value string
}

// loadInvocation converts the settings flags of the command line.
func loadInvocation(check valueChecker, flags []Flag) layerResult {
	result := layerResult{layer: Invocation}
	for _, flag := range flags {
		key, known := lookupFlag(flag.Name)
		if !known {
			result.rejection = Usage(fmt.Sprintf("%s is not a settings flag", flag.Name))
			return result
		}
		converted, ok := convertText(key, flag.Value)
		assignment, rejection := check.check(key, Invocation, flag.Name, converted, ok)
		if rejection != nil {
			result.rejection = rejection
			return result
		}
		result.assignments = append(result.assignments, assignment)
	}
	return result
}

// registered reports whether name is a registered runtime identifier.
func registered(name string) bool {
	for _, entry := range registeredEnvironment {
		if name == entry || strings.HasSuffix(entry, "_") && strings.HasPrefix(name, entry) {
			return true
		}
	}
	return false
}

// nearestName returns the env name with the smallest Levenshtein distance;
// ties go to the earlier key (settings.nearestName).
func nearestName(name string) string {
	best, bestDistance := "", -1
	for _, key := range Keys {
		if distance := levenshtein(name, key.Env); bestDistance < 0 || distance < bestDistance {
			best, bestDistance = key.Env, distance
		}
	}
	return best
}

func levenshtein(a, b string) int {
	previous := make([]int, len(b)+1)
	for j := range previous {
		previous[j] = j
	}
	for i := 1; i <= len(a); i++ {
		current := make([]int, len(b)+1)
		current[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			current[j] = min(previous[j]+1, current[j-1]+1, previous[j-1]+cost)
		}
		previous = current
	}
	return previous[len(b)]
}
