package settings

import (
	"path/filepath"
	"strings"
)

// Request holds everything the resolver reads apart from files. The
// environment is passed in, never read from or written to the process, so a
// value from an environment file can never reach a child process (CFG-043).
type Request struct {
	// Mode selects the per-mode defaults and value rules.
	Mode Mode
	// GOOS selects settings.files.system.
	GOOS string
	// Environ is the process environment as KEY=VALUE strings.
	Environ []string
	// WorkDir is the working directory at process start, the anchor of
	// --settings-file and --env-file. The resolver never changes directory.
	WorkDir string
	// SettingsFile is the --settings-file value; it is used when
	// SettingsFileGiven is true.
	SettingsFile      string
	SettingsFileGiven bool
	// EnvFiles are the --env-file values, left to right.
	EnvFiles []string
	// Flags are the settings flags of the command line, in order.
	Flags []Flag
}

// Loaded holds every layer except the project layer. mcp serve keeps one
// Loaded for its whole life (CFG-055); each tool call adds the project
// layer of its project.
type Loaded struct {
	mode    Mode
	paths   Paths
	noHome  bool
	system  layerResult
	user    layerResult
	envFile []layerResult
	env     layerResult
	flags   layerResult
}

// Load reads the system and user files, the environment files, the
// environment and the flags. It never fails: defects are kept per layer and
// Resolve reports the first in the order of settings.failureOrder.
func Load(fsys FileSystem, request Request) *Loaded {
	env := environMap(request.Environ)
	paths, err := NewPaths(env)
	loaded := &Loaded{mode: request.Mode, paths: paths}
	if err != nil {
		loaded.noHome = true
		return loaded
	}
	check := valueChecker{mode: request.Mode, paths: paths}
	loaded.system = loadSettingsFile(fsys, check, System, SystemFile(request.GOOS), false)
	if request.SettingsFileGiven {
		loaded.user = loadSettingsFile(fsys, check, User, anchor(request.WorkDir, request.SettingsFile), true)
	} else {
		loaded.user = loadSettingsFile(fsys, check, User, paths.UserFile(), false)
	}
	loaded.envFile = loadEnvFiles(fsys, check, request, env)
	loaded.env = loadEnvironment(check, env)
	loaded.flags = loadInvocation(check, request.Flags)
	return loaded
}

// loadEnvFiles selects the environment files: every --env-file, left to
// right, or else AIBOX_ENV_FILE, which must be absolute (CFG-030 to CFG-032).
func loadEnvFiles(fsys FileSystem, check valueChecker, request Request, env map[string]string) []layerResult {
	var results []layerResult
	if len(request.EnvFiles) > 0 {
		for _, path := range request.EnvFiles {
			results = append(results, loadEnvFile(fsys, check, anchor(request.WorkDir, path)))
		}
		return results
	}
	if path, set := env[EnvFileVariable]; set {
		if !filepath.IsAbs(path) {
			return []layerResult{{layer: EnvFile, rejection: reject(RelativeEnvFile, EnvFile, EnvFileVariable, "", "AIBOX_ENV_FILE is not an absolute path")}}
		}
		results = append(results, loadEnvFile(fsys, check, filepath.Clean(path)))
	}
	return results
}

// Paths returns the per-user directories; they are empty when the home
// directory is unknown.
func (l *Loaded) Paths() Paths { return l.paths }

// Resolve combines the loaded layers with an optional project layer. It
// returns the first defect in the order no_home, system, user, project,
// env_file (load order), environment, invocation (CFG-008).
func (l *Loaded) Resolve(project *ProjectLayer) (*Resolved, *Rejection) {
	if l.noHome {
		return nil, &Rejection{Reason: NoHome, Message: "the home directory of the running user cannot be determined"}
	}
	ordered := []layerResult{l.system, l.user}
	if project != nil {
		ordered = append(ordered, project.result)
	}
	ordered = append(ordered, l.envFile...)
	ordered = append(ordered, l.env, l.flags)
	for _, result := range ordered {
		if result.rejection != nil {
			return nil, result.rejection
		}
	}
	return combine(l.mode, l.paths, ordered), nil
}

// Setting is the resolved value of one key with its provenance.
type Setting struct {
	// Key is the catalog entry.
	Key Key
	// Layer is the winning layer; built_in when no layer set the key.
	Layer Layer
	// Source is the unredacted winning source.
	Source string
	// Value is the winning value; nil means unset.
	Value any
	// Overridden lists lower values, highest layer first, without built_in.
	Overridden []Assignment
}

// Resolved is the effective configuration of one invocation or server.
type Resolved struct {
	// Mode is the running mode.
	Mode Mode
	// Paths are the per-user directories.
	Paths Paths
	// Settings has one element per key, in catalog order.
	Settings []Setting
	// Sources are the settingsSources elements in their defined order.
	Sources []SourceState
	// LayersSet lists the layers that set at least one key, lowest first.
	LayersSet []Layer
}

func combine(mode Mode, paths Paths, ordered []layerResult) *Resolved {
	resolved := &Resolved{Mode: mode, Paths: paths}
	set := map[Layer]bool{}
	var all []Assignment
	for _, result := range ordered {
		if result.source != nil {
			resolved.Sources = append(resolved.Sources, *result.source)
		}
		for _, assignment := range result.assignments {
			set[assignment.Layer] = true
			all = append(all, assignment)
		}
	}
	for _, layer := range Layers {
		if set[layer] {
			resolved.LayersSet = append(resolved.LayersSet, layer)
		}
	}
	for _, key := range Keys {
		setting := Setting{Key: key, Layer: BuiltIn, Source: BuiltInSource, Value: key.Default(mode)}
		var values []Assignment
		for _, assignment := range all {
			if assignment.Key == key.Name {
				values = append(values, assignment)
			}
		}
		if count := len(values); count > 0 {
			winner := values[count-1]
			setting.Layer, setting.Source, setting.Value = winner.Layer, winner.Source, winner.Value
			for index := count - 2; index >= 0; index-- {
				setting.Overridden = append(setting.Overridden, values[index])
			}
		}
		resolved.Settings = append(resolved.Settings, setting)
	}
	return resolved
}

// Get returns the resolved setting of a key of the catalog.
func (r *Resolved) Get(name string) Setting {
	for _, setting := range r.Settings {
		if setting.Key.Name == name {
			return setting
		}
	}
	panic("settings: unknown key " + name)
}

// Text returns a string-valued setting, or "" when it is unset.
func (r *Resolved) Text(name string) string {
	value, _ := r.Get(name).Value.(string)
	return value
}

// Integer returns an integer-valued setting, or 0 when it is unset.
func (r *Resolved) Integer(name string) int {
	value, _ := r.Get(name).Value.(int)
	return value
}

// Flag returns a boolean setting.
func (r *Resolved) Flag(name string) bool {
	value, _ := r.Get(name).Value.(bool)
	return value
}

// TimeLimitSeconds returns execution.timeoutSeconds, or the class default
// when no layer sets it (D5, EXE-001).
func (r *Resolved) TimeLimitSeconds(classDefault int) int {
	if limit := r.Integer(ExecutionTimeout); limit > 0 {
		return limit
	}
	return classDefault
}

func environMap(environ []string) map[string]string {
	env := make(map[string]string, len(environ))
	for _, entry := range environ {
		name, value, found := strings.Cut(entry, "=")
		if !found || name == "" {
			continue
		}
		if _, duplicate := env[name]; !duplicate {
			env[name] = value
		}
	}
	return env
}

func anchor(workDir, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(workDir, path)
}
