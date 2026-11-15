// Package settings resolves the aibox process settings of
// settings-and-logging@v1 from seven layers: built-in defaults, the system
// and user settings files, the project declaration, environment files, the
// process environment and the command line. It reads only; it never writes a
// file, changes the working directory or changes the process environment.
package settings

// Layer is one source of settings, listed from lowest to highest precedence.
type Layer string

const (
	// BuiltIn is the key's default for the running mode.
	BuiltIn Layer = "built_in"
	// System is the machine-wide settings file.
	System Layer = "system"
	// User is the user's settings file or the file named by --settings-file.
	User Layer = "user"
	// Project is customizations.aibox.workspace of the selected devcontainer.json.
	Project Layer = "project"
	// EnvFile is an environment file named by --env-file or AIBOX_ENV_FILE.
	EnvFile Layer = "env_file"
	// Environment is the process environment.
	Environment Layer = "environment"
	// Invocation is the command line of one invocation.
	Invocation Layer = "invocation"
)

// Layers lists every layer from lowest to highest precedence.
var Layers = []Layer{BuiltIn, System, User, Project, EnvFile, Environment, Invocation}

// Mode selects the per-mode defaults of the catalog.
type Mode string

const (
	// ModeCLI is a single CLI invocation.
	ModeCLI Mode = "cli"
	// ModeMCP is the local MCP server, mcp serve.
	ModeMCP Mode = "mcp"
)

// Type is the value type of a key; conversion rules depend on it.
type Type string

const (
	// Enum values are case-sensitive members of Key.Values.
	Enum Type = "enum"
	// Boolean values are exactly true or false.
	Boolean Type = "boolean"
	// Integer values are base-10 numbers within Key.Minimum and Key.Maximum.
	Integer Type = "integer"
	// Path values are non-empty strings.
	Path Type = "path"
)

// Key names of the catalog, for typed access by callers.
const (
	OutputFormat           = "output.format"
	OutputColor            = "output.color"
	LoggingLevel           = "logging.level"
	LoggingFormat          = "logging.format"
	LoggingFile            = "logging.file"
	LoggingRequired        = "logging.required"
	LoggingRotationMiB     = "logging.rotationMiB"
	LoggingRetentionFiles  = "logging.retentionFiles"
	LoggingRetentionDays   = "logging.retentionDays"
	ExecutionTimeout       = "execution.timeoutSeconds"
	ExecutionLockWait      = "execution.lockWaitSeconds"
	DiagnosticsOffline     = "diagnostics.offline"
	WorkspaceTheme         = "workspace.theme"
	WorkspaceMode          = "workspace.mode"
	WorkspaceLayout        = "workspace.layout"
	EnvFileVariable        = "AIBOX_ENV_FILE"
	BuiltInSource          = "built-in"
	MaxEnvFiles            = 16
	ReadClassLimitSeconds  = 30
	MutatingLimitSeconds   = 3600
	maxSettingsFileBytes   = 1024 * 1024
	maxEnvFileBytes        = 1024 * 1024
	maxProjectFileBytes    = 2 * 1024 * 1024
	environmentNamePrefix  = "AIBOX_"
	projectDeclarationPath = ".devcontainer/devcontainer.json"
)

// Key is one entry of settings.keys in interfaces.json.
type Key struct {
	// Name is the dotted key, for example output.format.
	Name string
	// Type selects the conversion rule.
	Type Type
	// Domain is the view domain: process or workspace.
	Domain string
	// Layers are the layers that may set the key; built_in is implied.
	Layers []Layer
	// Env is the process environment and env-file name of the key.
	Env string
	// Flag is the command-line flag, or empty when the key has none.
	Flag string
	// DefaultCLI and DefaultMCP are the defaults per mode; nil means unset.
	DefaultCLI, DefaultMCP any
	// Values are the members of an enum key.
	Values []string
	// Minimum and Maximum bound an integer key.
	Minimum, Maximum int
}

var processFileLayers = []Layer{System, User, EnvFile, Environment, Invocation}
var workspaceLayers = []Layer{System, User, Project, EnvFile, Environment}

// Keys is the catalog in the order of settings.keys. The view lists one
// entry per key in this order.
var Keys = []Key{
	{Name: OutputFormat, Type: Enum, Domain: "process", Layers: processFileLayers, Env: "AIBOX_OUTPUT__FORMAT", Flag: "--format",
		DefaultCLI: "human", DefaultMCP: "json", Values: []string{"human", "json", "yaml"}},
	{Name: OutputColor, Type: Enum, Domain: "process", Layers: processFileLayers, Env: "AIBOX_OUTPUT__COLOR", Flag: "--color",
		DefaultCLI: "auto", DefaultMCP: "auto", Values: []string{"auto", "always", "never"}},
	{Name: LoggingLevel, Type: Enum, Domain: "process", Layers: processFileLayers, Env: "AIBOX_LOGGING__LEVEL", Flag: "--log-level",
		DefaultCLI: "info", DefaultMCP: "info", Values: []string{"trace", "debug", "info", "warn", "error", "off"}},
	{Name: LoggingFormat, Type: Enum, Domain: "process", Layers: processFileLayers, Env: "AIBOX_LOGGING__FORMAT", Flag: "--log-format",
		DefaultCLI: "plain", DefaultMCP: "jsonl", Values: []string{"plain", "jsonl"}},
	{Name: LoggingFile, Type: Path, Domain: "process", Layers: []Layer{System, User, Invocation}, Env: "AIBOX_LOGGING__FILE", Flag: "--log-file"},
	{Name: LoggingRequired, Type: Boolean, Domain: "process", Layers: []Layer{System, User}, Env: "AIBOX_LOGGING__REQUIRED",
		DefaultCLI: false, DefaultMCP: false},
	{Name: LoggingRotationMiB, Type: Integer, Domain: "process", Layers: []Layer{System, User}, Env: "AIBOX_LOGGING__ROTATION_MIB",
		DefaultCLI: 10, DefaultMCP: 10, Minimum: 1, Maximum: 1024},
	{Name: LoggingRetentionFiles, Type: Integer, Domain: "process", Layers: []Layer{System, User}, Env: "AIBOX_LOGGING__RETENTION_FILES",
		DefaultCLI: 7, DefaultMCP: 7, Minimum: 1, Maximum: 100},
	{Name: LoggingRetentionDays, Type: Integer, Domain: "process", Layers: []Layer{System, User}, Env: "AIBOX_LOGGING__RETENTION_DAYS",
		DefaultCLI: 7, DefaultMCP: 7, Minimum: 1, Maximum: 365},
	{Name: ExecutionTimeout, Type: Integer, Domain: "process", Layers: processFileLayers, Env: "AIBOX_EXECUTION__TIMEOUT_SECONDS", Flag: "--timeout",
		Minimum: 1, Maximum: 86400},
	{Name: ExecutionLockWait, Type: Integer, Domain: "process", Layers: []Layer{System, User, EnvFile, Environment}, Env: "AIBOX_EXECUTION__LOCK_WAIT_SECONDS",
		DefaultCLI: 30, DefaultMCP: 30, Minimum: 0, Maximum: 600},
	{Name: DiagnosticsOffline, Type: Boolean, Domain: "process", Layers: []Layer{System, User, EnvFile, Environment}, Env: "AIBOX_DIAGNOSTICS__OFFLINE",
		DefaultCLI: true, DefaultMCP: true},
	{Name: WorkspaceTheme, Type: Enum, Domain: "workspace", Layers: workspaceLayers, Env: "AIBOX_WORKSPACE__THEME", Values: []string{
		"andromeeda", "aurora-x", "ayu", "borland", "catppuccin", "dracula", "everforest", "github", "gruvbox", "houston",
		"kanagawa", "laserwave", "material", "min", "monokai", "mono", "moonlight", "night-owl", "nord", "norton",
		"one-dark", "plastic", "poimandres", "phosphor", "projectious", "contrast", "contrast-mono", "red", "rose-pine",
		"slack", "snazzy", "solarized", "synthwave84", "tokyo-night", "vesper", "vitesse", "vs-code"}},
	{Name: WorkspaceMode, Type: Enum, Domain: "workspace", Layers: workspaceLayers, Env: "AIBOX_WORKSPACE__MODE", Values: []string{"dark", "light", "auto"}},
	{Name: WorkspaceLayout, Type: Enum, Domain: "workspace", Layers: workspaceLayers, Env: "AIBOX_WORKSPACE__LAYOUT", Values: []string{"ai", "dev", "focus", "cowork"}},
}

// Lookup returns the catalog key with the given dotted name.
func Lookup(name string) (Key, bool) {
	for _, key := range Keys {
		if key.Name == name {
			return key, true
		}
	}
	return Key{}, false
}

// lookupEnv returns the catalog key whose env name is name.
func lookupEnv(name string) (Key, bool) {
	for _, key := range Keys {
		if key.Env == name {
			return key, true
		}
	}
	return Key{}, false
}

// lookupFlag returns the catalog key set by a command-line flag.
func lookupFlag(flag string) (Key, bool) {
	for _, key := range Keys {
		if key.Flag != "" && key.Flag == flag {
			return key, true
		}
	}
	return Key{}, false
}

// Allows reports whether layer may set the key (CFG-002).
func (k Key) Allows(layer Layer) bool {
	for _, allowed := range k.Layers {
		if allowed == layer {
			return true
		}
	}
	return false
}

// Default returns the key's default for the mode; nil means unset.
func (k Key) Default(mode Mode) any {
	if mode == ModeMCP {
		return k.DefaultMCP
	}
	return k.DefaultCLI
}
