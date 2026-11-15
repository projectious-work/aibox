package settings

import (
	"errors"
	"os"
	"path/filepath"
)

// Paths holds the per-user directories that settings and the file sink use.
// It is the only place where the XDG rules and the log-file anchoring are
// written; the logging package receives their results.
type Paths struct {
	// Home is the running user's home directory, taken from HOME only.
	Home string
	// ConfigHome is XDG_CONFIG_HOME when absolute, otherwise $HOME/.config.
	ConfigHome string
	// StateHome is XDG_STATE_HOME when absolute, otherwise $HOME/.local/state.
	StateHome string
}

var errNoHome = errors.New("HOME is not set to an absolute path")
var errPathEscape = errors.New("the log file is outside the permitted root")

// NewPaths derives the directories from the environment. The home directory
// comes from HOME only, never from the user database, so the binary stays
// static (CFG-024). A HOME that is empty or relative cannot be used.
func NewPaths(env map[string]string) (Paths, error) {
	home := env["HOME"]
	if home == "" || !filepath.IsAbs(home) {
		return Paths{}, errNoHome
	}
	home = filepath.Clean(home)
	paths := Paths{Home: home, ConfigHome: filepath.Join(home, ".config"), StateHome: filepath.Join(home, ".local", "state")}
	if value := env["XDG_CONFIG_HOME"]; filepath.IsAbs(value) {
		paths.ConfigHome = filepath.Clean(value)
	}
	if value := env["XDG_STATE_HOME"]; filepath.IsAbs(value) {
		paths.StateHome = filepath.Clean(value)
	}
	return paths, nil
}

// SystemFile returns settings.files.system for the operating system.
func SystemFile(goos string) string {
	if goos == "darwin" {
		return "/Library/Application Support/aibox/settings.json"
	}
	return "/etc/aibox/settings.json"
}

// UserFile returns the default user settings file (CFG-011).
func (p Paths) UserFile() string {
	return filepath.Join(p.ConfigHome, "aibox", "settings.json")
}

// LogDirectory is logging.fileSink.stateDirectory and permittedRoot.
func (p Paths) LogDirectory() string {
	return filepath.Join(p.StateHome, "aibox", "logs")
}

// OperationRecordDirectory is logging.fileSink.operationRecordDirectory. The
// file sink never touches it.
func (p Paths) OperationRecordDirectory() string {
	return filepath.Join(p.StateHome, "aibox", "operations")
}

// AnchorLogFile makes a logging.file value absolute: a relative value is
// anchored to the log directory (LOG-013), whatever layer set it.
func (p Paths) AnchorLogFile(value string) string {
	if filepath.IsAbs(value) {
		return filepath.Clean(value)
	}
	return filepath.Join(p.LogDirectory(), value)
}

// ConfineLogFile returns the symlink-resolved form of an anchored log file,
// or errPathEscape when that form is not inside the symlink-resolved log
// directory (LOG-014). Components that do not exist yet are kept as written.
func (p Paths) ConfineLogFile(anchored string) (string, error) {
	root, err := resolveExisting(p.LogDirectory())
	if err != nil {
		return "", errPathEscape
	}
	candidate, err := resolveExisting(filepath.Clean(anchored))
	if err != nil {
		return "", errPathEscape
	}
	relative, err := filepath.Rel(root, candidate)
	if err != nil || relative == "." || relative == ".." || len(relative) > 2 && relative[:3] == ".."+string(filepath.Separator) {
		return "", errPathEscape
	}
	return candidate, nil
}

// resolveExisting resolves symlinks in the longest existing prefix of path
// and appends the rest unchanged. A dangling symlink is an error, because
// creating the file would follow it to an unchecked place.
func resolveExisting(path string) (string, error) {
	current, rest := path, ""
	for {
		if _, err := os.Lstat(current); err == nil {
			resolved, err := filepath.EvalSymlinks(current)
			if err != nil {
				return "", err
			}
			return filepath.Join(resolved, rest), nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			return path, nil
		}
		rest = filepath.Join(filepath.Base(current), rest)
		current = parent
	}
}

// ResolvedPath returns the symlink-resolved form of an existing path, or the
// path itself when it cannot be resolved.
func ResolvedPath(path string) string {
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		return resolved
	}
	return path
}
