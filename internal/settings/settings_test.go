package settings

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// overlayFS is a fake FileSystem for the system layer, whose fixed path a
// test cannot write. Paths in files are served from memory (nil = absent);
// every other path is read from the real, temporary file tree.
type overlayFS struct{ files map[string][]byte }

func (o overlayFS) ReadFile(path string, limit int64) (string, []byte, error) {
	if contents, ok := o.files[path]; ok {
		if contents == nil {
			return path, nil, fs.ErrNotExist
		}
		return path, contents, nil
	}
	return OSFileSystem{}.ReadFile(path, limit)
}

// world is one isolated user: a temporary HOME, working directory, project
// and environment.
type world struct {
	t       *testing.T
	home    string
	work    string
	project string
	env     map[string]string
	fs      overlayFS
	request Request
}

func newWorld(t *testing.T) *world {
	t.Helper()
	base := t.TempDir()
	w := &world{t: t, home: filepath.Join(base, "home"), work: filepath.Join(base, "work"), project: filepath.Join(base, "project"),
		env: map[string]string{}, fs: overlayFS{files: map[string][]byte{SystemFile("linux"): nil}}}
	for _, dir := range []string{w.home, w.work, filepath.Join(w.project, ".devcontainer")} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	w.env["HOME"] = w.home
	w.env["PATH"] = "/usr/bin:/bin"
	w.writeProject(`{"image":"example.invalid/base:one"}`)
	w.request = Request{Mode: ModeCLI, GOOS: "linux", WorkDir: w.work}
	return w
}

func (w *world) write(path, contents string) string {
	w.t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		w.t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		w.t.Fatal(err)
	}
	return path
}

func (w *world) userFile() string { return filepath.Join(w.home, ".config", "aibox", "settings.json") }
func (w *world) writeUser(contents string) string {
	return w.write(w.userFile(), contents)
}
func (w *world) writeProject(contents string) {
	w.write(filepath.Join(w.project, ".devcontainer", "devcontainer.json"), contents)
}
func (w *world) setSystem(contents string) { w.fs.files[SystemFile("linux")] = []byte(contents) }

func (w *world) resolve() (*Resolved, *Rejection) {
	w.t.Helper()
	request := w.request
	for name, value := range w.env {
		request.Environ = append(request.Environ, name+"="+value)
	}
	return Load(w.fs, request).Resolve(ReadProject(w.project))
}

func (w *world) mustResolve() *Resolved {
	w.t.Helper()
	resolved, rejection := w.resolve()
	if rejection != nil {
		w.t.Fatalf("unexpected rejection: %+v", rejection)
	}
	return resolved
}

func (w *world) mustReject(reason Reason) *Rejection {
	w.t.Helper()
	resolved, rejection := w.resolve()
	if rejection == nil {
		w.t.Fatalf("expected %s, resolved %+v", reason, resolved.Settings)
	}
	if rejection.Reason != reason {
		w.t.Fatalf("expected %s, got %+v", reason, rejection)
	}
	return rejection
}

func TestDefaultsPerMode(t *testing.T) {
	// CFG-004, CFG-007: unset keys take the mode's default with layer built_in.
	for _, mode := range []Mode{ModeCLI, ModeMCP} {
		w := newWorld(t)
		w.request.Mode = mode
		resolved := w.mustResolve()
		if len(resolved.Settings) != 15 {
			t.Fatalf("%s: %d settings", mode, len(resolved.Settings))
		}
		for index, setting := range resolved.Settings {
			if setting.Key.Name != Keys[index].Name || setting.Layer != BuiltIn || setting.Source != "built-in" ||
				!reflect.DeepEqual(setting.Value, Keys[index].Default(mode)) || setting.Overridden != nil {
				t.Fatalf("%s: %+v", mode, setting)
			}
		}
		if resolved.Text(OutputFormat) != map[Mode]string{ModeCLI: "human", ModeMCP: "json"}[mode] ||
			resolved.Text(LoggingFormat) != map[Mode]string{ModeCLI: "plain", ModeMCP: "jsonl"}[mode] {
			t.Fatalf("%s: format defaults %v %v", mode, resolved.Text(OutputFormat), resolved.Text(LoggingFormat))
		}
		if resolved.TimeLimitSeconds(ReadClassLimitSeconds) != 30 || len(resolved.LayersSet) != 0 {
			t.Fatalf("%s: time limit or layers %v", mode, resolved.LayersSet)
		}
		want := []SourceState{{System, "/etc/aibox/settings.json", "absent"}, {User, w.userFile(), "absent"},
			{Project, filepath.Join(ResolvedPath(w.project), ".devcontainer", "devcontainer.json"), "loaded"}}
		if !reflect.DeepEqual(resolved.Sources, want) {
			t.Fatalf("%s: sources %+v", mode, resolved.Sources)
		}
	}
}

// sampleValues gives two different valid values per key.
var sampleValues = map[string][2]any{
	OutputFormat: {"json", "yaml"}, OutputColor: {"always", "never"}, LoggingLevel: {"debug", "warn"},
	LoggingFormat: {"plain", "jsonl"}, LoggingFile: {"a.log", "b.log"}, LoggingRequired: {true, false},
	LoggingRotationMiB: {5, 6}, LoggingRetentionFiles: {3, 4}, LoggingRetentionDays: {2, 3},
	ExecutionTimeout: {5, 6}, ExecutionLockWait: {1, 2}, DiagnosticsOffline: {false, true},
	WorkspaceTheme: {"nord", "dracula"}, WorkspaceMode: {"dark", "light"}, WorkspaceLayout: {"dev", "ai"},
}

func settingsJSON(key Key, value any) string {
	group, member, _ := strings.Cut(key.Name, ".")
	encoded, _ := json.Marshal(value)
	return fmt.Sprintf(`{"schemaVersion":"1",%q:{%q:%s}}`, group, member, encoded)
}

// set makes layer give key the value, the way a user or operator would.
func (w *world) set(layer Layer, key Key, value any) string {
	w.t.Helper()
	text := fmt.Sprint(value)
	switch layer {
	case System:
		w.setSystem(settingsJSON(key, value))
		return SystemFile("linux")
	case User:
		return w.writeUser(settingsJSON(key, value))
	case Project:
		_, member, _ := strings.Cut(key.Name, ".")
		w.writeProject(fmt.Sprintf(`{"image":"x", // comment
			"customizations":{"aibox":{"workspace":{%q:%q}}}}`, member, text))
		return filepath.Join(ResolvedPath(w.project), ".devcontainer", "devcontainer.json")
	case EnvFile:
		path := w.write(filepath.Join(w.work, "settings.env"), key.Env+"="+text+"\n")
		w.request.EnvFiles = append(w.request.EnvFiles, "settings.env")
		return path
	case Environment:
		w.env[key.Env] = text
		return key.Env
	case Invocation:
		w.request.Flags = append(w.request.Flags, Flag{Name: key.Flag, Value: text})
		return key.Flag
	}
	w.t.Fatalf("no setter for %s", layer)
	return ""
}

func (w *world) expected(key Key, value any) any {
	if key.Name == LoggingFile {
		return filepath.Join(ResolvedPath(w.home), ".local", "state", "aibox", "logs", value.(string))
	}
	return value
}

func TestPrecedenceForEveryKeyAndLayerPair(t *testing.T) {
	// CFG-001, CFG-002, CFG-062: the higher allowed layer wins; the lower one
	// is listed in overridden; built_in is never listed.
	for _, key := range Keys {
		for low := 0; low < len(key.Layers); low++ {
			for high := low + 1; high < len(key.Layers); high++ {
				lower, higher := key.Layers[low], key.Layers[high]
				t.Run(fmt.Sprintf("%s/%s<%s", key.Name, lower, higher), func(t *testing.T) {
					w := newWorld(t)
					values := sampleValues[key.Name]
					lowSource := w.set(lower, key, values[0])
					highSource := w.set(higher, key, values[1])
					setting := w.mustResolve().Get(key.Name)
					if setting.Layer != higher || setting.Source != highSource || !reflect.DeepEqual(setting.Value, w.expected(key, values[1])) {
						t.Fatalf("winner %+v", setting)
					}
					want := []Assignment{{Key: key.Name, Layer: lower, Source: lowSource, Value: w.expected(key, values[0])}}
					if !reflect.DeepEqual(setting.Overridden, want) {
						t.Fatalf("overridden %+v, want %+v", setting.Overridden, want)
					}
				})
			}
		}
	}
}

func TestOverriddenListsEveryLowerLayerHighestFirst(t *testing.T) {
	// J-CFG-03 at component level.
	w := newWorld(t)
	key, _ := Lookup(OutputColor)
	w.set(System, key, "never")
	w.set(User, key, "always")
	w.set(EnvFile, key, "never")
	w.set(Environment, key, "always")
	w.set(Invocation, key, "never")
	setting := w.mustResolve().Get(OutputColor)
	var layers []Layer
	for _, item := range setting.Overridden {
		layers = append(layers, item.Layer)
	}
	if setting.Value != "never" || setting.Layer != Invocation || !reflect.DeepEqual(layers, []Layer{Environment, EnvFile, User, System}) {
		t.Fatalf("%+v", setting)
	}
}

func TestUnauthorizedLayersFailWithoutFallback(t *testing.T) {
	// CFG-003: a layer not listed for a key fails; no lower value is used.
	for _, key := range Keys {
		for _, layer := range []Layer{Project, EnvFile, Environment, Invocation} {
			if key.Allows(layer) || layer == Invocation && key.Flag == "" || layer == Project && key.Domain != "workspace" {
				continue
			}
			t.Run(key.Name+"/"+string(layer), func(t *testing.T) {
				w := newWorld(t)
				w.set(User, key, sampleValues[key.Name][0])
				source := w.set(layer, key, sampleValues[key.Name][1])
				rejection := w.mustReject(UnauthorizedSource)
				if rejection.Key != key.Name || rejection.Layer != layer || rejection.Source != source {
					t.Fatalf("%+v", rejection)
				}
			})
		}
	}
}

func TestProjectLayerReadsOnlyTheThreeWorkspaceMembers(t *testing.T) {
	// CFG-019, CFG-020, CFG-022: no settings file in the project or working
	// directory, and no other member of customizations.aibox, is read.
	w := newWorld(t)
	for _, dir := range []string{w.project, filepath.Join(w.project, ".devcontainer"), w.work} {
		w.write(filepath.Join(dir, "settings.json"), `{"schemaVersion":"1","logging":{"level":"trace"}}`)
	}
	w.writeProject(`{"image":"x","customizations":{"aibox":{"logging":{"level":"trace"},
		"workspace":{"theme":"nord","mode":"dark","layout":"dev","emphasis":"full"}}}}`)
	resolved := w.mustResolve()
	if setting := resolved.Get(LoggingLevel); setting.Value != "info" || setting.Layer != BuiltIn {
		t.Fatalf("%+v", setting)
	}
	for name, want := range map[string]string{WorkspaceTheme: "nord", WorkspaceMode: "dark", WorkspaceLayout: "dev"} {
		if setting := resolved.Get(name); setting.Value != want || setting.Layer != Project || setting.Key.Domain != "workspace" {
			t.Fatalf("%+v", setting)
		}
	}
	if !reflect.DeepEqual(resolved.LayersSet, []Layer{Project}) {
		t.Fatalf("layers %v", resolved.LayersSet)
	}
}

func TestWorkspaceValuesOutsideTheCatalogFail(t *testing.T) {
	// CFG-021 from the project, a settings file and the environment.
	w := newWorld(t)
	w.writeProject(`{"customizations":{"aibox":{"workspace":{"layout":"grid"}}}}`)
	if rejection := w.mustReject(InvalidValue); rejection.Layer != Project || rejection.Key != WorkspaceLayout {
		t.Fatalf("%+v", rejection)
	}
	w = newWorld(t)
	w.writeUser(`{"schemaVersion":"1","workspace":{"theme":"no-such-theme"}}`)
	w.mustReject(InvalidValue)
	w = newWorld(t)
	w.env["AIBOX_WORKSPACE__MODE"] = "sepia"
	w.mustReject(InvalidValue)
}

func TestSettingsFileDefects(t *testing.T) {
	// CFG-013 to CFG-016 and settings.schemaDefects (J-CFG-07).
	cases := []struct {
		contents string
		reason   Reason
	}{
		{"not json", Malformed},
		{`["schemaVersion"]`, Malformed},
		{`{"schemaVersion":"1","operatorPolicy":{}}`, UnknownKey},
		{`{"schemaVersion":"1","logging":{"verbose":true}}`, UnknownKey},
		{`{"schemaVersion":"2"}`, UnsupportedVersion},
		{`{"schemaVersion":1}`, UnsupportedVersion},
		{`{"logging":{"level":"debug"}}`, UnsupportedVersion},
		{`{"schemaVersion":"1","logging":{"level":"verbose"}}`, InvalidValue},
		{`{"schemaVersion":"1","logging":{"rotationMiB":0}}`, InvalidValue},
		{`{"schemaVersion":"1","logging":{"rotationMiB":1.5}}`, InvalidValue},
		{`{"schemaVersion":"1","logging":{"required":"true"}}`, InvalidValue},
		{`{"schemaVersion":"1","logging":"debug"}`, InvalidValue},
		{`{"schemaVersion":"1","logging":{},"logging":{}}`, DuplicateKey},
		{`{"schemaVersion":"1","output":{"format":"json","format":"yaml"}}`, DuplicateKey},
		{`{"schemaVersion":"1","output":null}`, NullValue},
		{`{"schemaVersion":"1","output":{"color":null}}`, NullValue},
		{`{"schemaVersion":null}`, NullValue},
		{`{"schemaVersion":"1","logging":{"file":"/var/log/aibox.log"}}`, PathEscape},
		{`{"schemaVersion":"1","logging":{"file":"../x.log"}}`, PathEscape},
	}
	for _, test := range cases {
		w := newWorld(t)
		path := w.writeUser(test.contents)
		rejection := w.mustReject(test.reason)
		if rejection.Layer != User || rejection.Source != path {
			t.Fatalf("%s: %+v", test.contents, rejection)
		}
	}
	w := newWorld(t)
	if err := os.Chmod(w.writeUser(`{"schemaVersion":"1"}`), 0); err != nil {
		t.Fatal(err)
	}
	if os.Geteuid() != 0 {
		w.mustReject(Unreadable)
	}
	w = newWorld(t)
	if err := os.MkdirAll(w.userFile(), 0o700); err != nil {
		t.Fatal(err)
	}
	w.mustReject(Unreadable)
}

func TestDocumentOrderWithinAFile(t *testing.T) {
	// CFG-008: within a file the first defect in document order is reported;
	// a missing schemaVersion has no position and comes last.
	w := newWorld(t)
	w.writeUser(`{"logging":{"level":"loud"},"schemaVersion":"2","x":1}`)
	w.mustReject(InvalidValue)
	w = newWorld(t)
	w.writeUser(`{"x":1,"logging":{"level":"loud"}}`)
	w.mustReject(UnknownKey)
}

func TestSettingsFileSelection(t *testing.T) {
	// CFG-011, CFG-012, CFG-017, CFG-018.
	w := newWorld(t)
	w.writeUser(`{"schemaVersion":"1","output":{"color":"always"}}`)
	other := w.write(filepath.Join(w.work, "other.json"), `{"schemaVersion":"1","output":{"color":"never"}}`)
	w.request.SettingsFile, w.request.SettingsFileGiven = "other.json", true
	resolved := w.mustResolve()
	if setting := resolved.Get(OutputColor); setting.Value != "never" || setting.Source != other || setting.Overridden != nil {
		t.Fatalf("%+v", setting)
	}
	if resolved.Sources[1] != (SourceState{User, other, "loaded"}) {
		t.Fatalf("%+v", resolved.Sources)
	}
	w.request.SettingsFile = "missing.json"
	if rejection := w.mustReject(MissingSelectedFile); rejection.Layer != User {
		t.Fatalf("%+v", rejection)
	}

	w = newWorld(t)
	config := filepath.Join(t.TempDir(), "xdg")
	w.env["XDG_CONFIG_HOME"] = config
	path := w.write(filepath.Join(config, "aibox", "settings.json"), `{"schemaVersion":"1","logging":{"level":"debug"}}`)
	if setting := w.mustResolve().Get(LoggingLevel); setting.Value != "debug" || setting.Source != path {
		t.Fatalf("%+v", setting)
	}
	w = newWorld(t)
	w.env["XDG_CONFIG_HOME"] = "rel"
	w.write(filepath.Join(w.work, "rel", "aibox", "settings.json"), `{"schemaVersion":"1","logging":{"level":"debug"}}`)
	if setting := w.mustResolve().Get(LoggingLevel); setting.Value != "info" || setting.Layer != BuiltIn {
		t.Fatalf("relative XDG_CONFIG_HOME used: %+v", setting)
	}
}

func TestSystemFileByOperatingSystem(t *testing.T) {
	// CFG-010; macOS paths are specified and unit-tested only (targets note).
	if SystemFile("linux") != "/etc/aibox/settings.json" || SystemFile("darwin") != "/Library/Application Support/aibox/settings.json" {
		t.Fatal(SystemFile("linux"), SystemFile("darwin"))
	}
	w := newWorld(t)
	w.setSystem(`{"schemaVersion":"1","logging":{"level":"warn"}}`)
	if setting := w.mustResolve().Get(LoggingLevel); setting.Value != "warn" || setting.Layer != System || setting.Source != "/etc/aibox/settings.json" {
		t.Fatalf("%+v", setting)
	}
}

func TestEnvFileGrammar(t *testing.T) {
	// CFG-033 to CFG-038 (J-CFG-12, J-CFG-13).
	cases := []struct {
		contents string
		reason   Reason
	}{
		{"export AIBOX_LOGGING__LEVEL=debug\n", Malformed},
		{" AIBOX_LOGGING__LEVEL=debug\n", Malformed},
		{"aibox_logging__level=debug\n", Malformed},
		{"AIBOX_LOGGING__LEVEL\n", Malformed},
		{"AIBOX_LOGGING__LEVEL=debug\nAIBOX_LOGGING__LEVEL=warn\n", DuplicateKey},
		{"AIBOX_NOPE=1\n", UnknownKey},
		{"PATH=/x\n", UnknownKey},
		{"AIBOX_ENV_FILE=/x\n", UnknownKey},
		{"AIBOX_WORKSPACE__THEME=\"nord\"\n", InvalidValue},
		{"AIBOX_LOGGING__LEVEL=debug \n", InvalidValue},
		{"AIBOX_LOGGING__REQUIRED=true\n", UnauthorizedSource},
		{"AIBOX_LOGGING__FILE=\n", InvalidValue},
		{"\xff=1\n", Malformed},
	}
	for _, test := range cases {
		w := newWorld(t)
		w.write(filepath.Join(w.work, "a.env"), test.contents)
		w.request.EnvFiles = []string{"a.env"}
		if rejection := w.mustReject(test.reason); rejection.Layer != EnvFile {
			t.Fatalf("%q: %+v", test.contents, rejection)
		}
	}
	w := newWorld(t)
	path := w.write(filepath.Join(w.work, "a.env"), "# comment\r\n\r\n \t\r\nAIBOX_WORKSPACE__THEME=nord\r\n")
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}
	w.request.EnvFiles = []string{"a.env"}
	if setting := w.mustResolve().Get(WorkspaceTheme); setting.Value != "nord" || setting.Layer != EnvFile || setting.Source != path {
		t.Fatalf("%+v", setting)
	}
}

func TestEnvFileSelection(t *testing.T) {
	// CFG-030 to CFG-032, CFG-036 (J-CFG-11, J-CFG-12).
	w := newWorld(t)
	w.write(filepath.Join(w.work, ".env"), "AIBOX_LOGGING__LEVEL=trace\n")
	w.write(filepath.Join(w.project, ".env"), "AIBOX_LOGGING__LEVEL=trace\n")
	if setting := w.mustResolve().Get(LoggingLevel); setting.Layer != BuiltIn {
		t.Fatalf(".env loaded implicitly: %+v", setting)
	}
	a := w.write(filepath.Join(w.work, "a.env"), "AIBOX_LOGGING__LEVEL=debug\n")
	b := w.write(filepath.Join(w.work, "b.env"), "AIBOX_LOGGING__LEVEL=warn\n")
	w.request.EnvFiles = []string{"a.env", "b.env"}
	resolved := w.mustResolve()
	setting := resolved.Get(LoggingLevel)
	if setting.Value != "warn" || setting.Source != b || len(setting.Overridden) != 1 || setting.Overridden[0].Source != a {
		t.Fatalf("%+v", setting)
	}
	if got := resolved.Sources[3:]; !reflect.DeepEqual(got, []SourceState{{EnvFile, a, "loaded"}, {EnvFile, b, "loaded"}}) {
		t.Fatalf("%+v", got)
	}
	w.request.EnvFiles = nil
	w.env[EnvFileVariable] = a
	if setting := w.mustResolve().Get(LoggingLevel); setting.Value != "debug" || setting.Source != a {
		t.Fatalf("%+v", setting)
	}
	w.request.EnvFiles = []string{"b.env"}
	resolved = w.mustResolve()
	if setting := resolved.Get(LoggingLevel); setting.Value != "warn" || len(resolved.Sources) != 4 {
		t.Fatalf("AIBOX_ENV_FILE used beside --env-file: %+v %+v", setting, resolved.Sources)
	}
	w.request.EnvFiles = nil
	w.env[EnvFileVariable] = "rel.env"
	w.mustReject(RelativeEnvFile)
	delete(w.env, EnvFileVariable)
	w.request.EnvFiles = []string{"missing.env"}
	w.mustReject(MissingSelectedFile)
	w.request.EnvFiles = []string{"."}
	w.mustReject(NotRegularFile)
	if os.Geteuid() != 0 {
		if err := os.Chmod(a, 0); err != nil {
			t.Fatal(err)
		}
		w.request.EnvFiles = []string{"a.env"}
		w.mustReject(Unreadable)
	}
}

func TestProcessEnvironmentMapping(t *testing.T) {
	// CFG-039 to CFG-042 (J-CFG-14).
	w := newWorld(t)
	w.env["AIBOX_LOGGING__LEVEL"] = "debug"
	if setting := w.mustResolve().Get(LoggingLevel); setting.Value != "debug" || setting.Layer != Environment || setting.Source != "AIBOX_LOGGING__LEVEL" {
		t.Fatalf("%+v", setting)
	}
	w = newWorld(t)
	w.env["AIBOX_LOGING__LEVEL"] = "debug"
	if rejection := w.mustReject(UnknownEnvironmentName); rejection.Suggestion != "AIBOX_LOGGING__LEVEL" {
		t.Fatalf("%+v", rejection)
	}
	w = newWorld(t)
	for name, value := range map[string]string{"AIBOX_TMUX_SESSION": "x", "AIBOX_TMUX_ANYTHING": "y", "AIBOX_VERSION": "0.35.4", "AIBOX_LOG_LEVEL": "trace"} {
		w.env[name] = value
	}
	if resolved := w.mustResolve(); len(resolved.LayersSet) != 0 {
		t.Fatalf("registered names took part: %v", resolved.LayersSet)
	}
	for name, value := range map[string]string{"AIBOX_DIAGNOSTICS__OFFLINE": "TRUE", "AIBOX_EXECUTION__TIMEOUT_SECONDS": "+5",
		"AIBOX_EXECUTION__TIMEOUT_SECONDS ": "5", "AIBOX_OUTPUT__COLOR": ""} {
		w = newWorld(t)
		w.env[name] = value
		if strings.HasSuffix(name, " ") {
			w.mustReject(UnknownEnvironmentName)
			continue
		}
		w.mustReject(InvalidValue)
	}
	for _, value := range []string{"0", "86401", " 5", "5.0"} {
		w = newWorld(t)
		w.env["AIBOX_EXECUTION__TIMEOUT_SECONDS"] = value
		w.mustReject(InvalidValue)
	}
	w = newWorld(t)
	w.env["AIBOX_EXECUTION__TIMEOUT_SECONDS"] = "-0"
	w.mustReject(InvalidValue)
	w = newWorld(t)
	w.env["AIBOX_LOGGING__REDACT"] = "false"
	w.mustReject(UnknownEnvironmentName)
}

func TestNearestNameTiesGoToTheEarlierKey(t *testing.T) {
	if got := nearestName("AIBOX_WORKSPACE__X"); got != "AIBOX_WORKSPACE__MODE" {
		t.Fatalf("got %s", got)
	}
	if levenshtein("kitten", "sitting") != 3 || levenshtein("", "abc") != 3 {
		t.Fatal("levenshtein")
	}
}

func TestFailureOrder(t *testing.T) {
	// CFG-008: no_home, then system, user, project, env files in load order,
	// environment, invocation; a shape defect before unauthorized_source.
	w := newWorld(t)
	w.writeUser("not json")
	w.request.Flags = []Flag{{"--log-level", "loud"}}
	w.env["AIBOX_NOPE"] = "1"
	w.mustReject(Malformed)

	w = newWorld(t)
	w.setSystem(`{"schemaVersion":"2"}`)
	w.writeUser("not json")
	w.mustReject(UnsupportedVersion)

	w = newWorld(t)
	w.writeProject(`{"customizations":{"aibox":{"workspace":{"layout":"grid"}}}}`)
	w.write(filepath.Join(w.work, "a.env"), "AIBOX_NOPE=1\n")
	w.request.EnvFiles = []string{"a.env"}
	w.mustReject(InvalidValue)

	w = newWorld(t)
	w.write(filepath.Join(w.work, "a.env"), "AIBOX_LOGGING__LEVEL=debug\n")
	w.write(filepath.Join(w.work, "b.env"), "AIBOX_NOPE=1\n")
	w.request.EnvFiles = []string{"missing.env", "b.env"}
	w.mustReject(MissingSelectedFile)
	w.request.EnvFiles = []string{"b.env", "missing.env"}
	w.mustReject(UnknownKey)

	w = newWorld(t)
	w.env["AIBOX_ZZZ"] = "1"
	w.env["AIBOX_LOGGING__FILE"] = "x.log"
	if rejection := w.mustReject(UnauthorizedSource); rejection.Key != LoggingFile {
		t.Fatalf("%+v", rejection)
	}
	w = newWorld(t)
	w.env["AIBOX_LOGGING__FILE"] = ""
	w.mustReject(InvalidValue)

	w = newWorld(t)
	w.writeUser("not json")
	delete(w.env, "HOME")
	w.mustReject(NoHome)
}

func TestNoHome(t *testing.T) {
	// CFG-024: HOME only, no user database; unset, empty or relative fails.
	for _, home := range []string{"", "relative/home"} {
		w := newWorld(t)
		w.env["HOME"] = home
		w.mustReject(NoHome)
	}
}

func TestLogFileAnchoringAndConfinement(t *testing.T) {
	// LOG-013, LOG-014 (J-LOG-06); D8: only inside <state>/aibox/logs.
	w := newWorld(t)
	logs := filepath.Join(w.home, ".local", "state", "aibox", "logs")
	w.request.Flags = []Flag{{"--log-file", "run.log"}}
	if got := w.mustResolve().Text(LoggingFile); got != filepath.Join(ResolvedPath(w.home), ".local/state/aibox/logs/run.log") {
		t.Fatalf("got %s", got)
	}
	state := filepath.Join(t.TempDir(), "state")
	w.env["XDG_STATE_HOME"] = state
	w.request.Flags = []Flag{{"--log-file", "sub/run.log"}}
	if got := w.mustResolve().Text(LoggingFile); got != filepath.Join(ResolvedPath(filepath.Dir(state)), "state/aibox/logs/sub/run.log") {
		t.Fatalf("got %s", got)
	}
	delete(w.env, "XDG_STATE_HOME")
	for _, value := range []string{"/tmp/x.log", "../x.log", logs, filepath.Join(w.home, ".local/state/aibox/operations/keep.json")} {
		w.request.Flags = []Flag{{"--log-file", value}}
		if rejection := w.mustReject(PathEscape); rejection.Layer != Invocation || rejection.Source != "--log-file" {
			t.Fatalf("%s: %+v", value, rejection)
		}
	}
	outside := w.write(filepath.Join(t.TempDir(), "outside.log"), "")
	if err := os.MkdirAll(logs, 0o700); err != nil {
		t.Fatal(err)
	}
	for name, target := range map[string]string{"link.log": outside, "dangling.log": filepath.Join(t.TempDir(), "none", "x.log")} {
		if err := os.Symlink(target, filepath.Join(logs, name)); err != nil {
			t.Fatal(err)
		}
		w.request.Flags = []Flag{{"--log-file", name}}
		w.mustReject(PathEscape)
	}
	inside := w.write(filepath.Join(logs, "real.log"), "")
	if err := os.Symlink(inside, filepath.Join(logs, "alias.log")); err != nil {
		t.Fatal(err)
	}
	w.request.Flags = []Flag{{"--log-file", "alias.log"}}
	if got := w.mustResolve().Text(LoggingFile); got != ResolvedPath(inside) {
		t.Fatalf("got %s", got)
	}
}

func TestMCPModeRejectsPlainEventsFromAnyLayer(t *testing.T) {
	// CFG-054.
	w := newWorld(t)
	w.request.Mode = ModeMCP
	w.writeUser(`{"schemaVersion":"1","logging":{"format":"plain"}}`)
	if rejection := w.mustReject(InvalidValue); rejection.Layer != User {
		t.Fatalf("%+v", rejection)
	}
	w = newWorld(t)
	w.request.Mode = ModeMCP
	w.env["AIBOX_LOGGING__FORMAT"] = "plain"
	w.mustReject(InvalidValue)
	w = newWorld(t)
	w.request.Mode = ModeMCP
	w.request.Flags = []Flag{{"--log-format", "plain"}}
	w.mustReject(InvalidValue)
}

func TestFlagValues(t *testing.T) {
	// CFG-050, CFG-051.
	w := newWorld(t)
	w.request.Flags = []Flag{{"--log-level", "debug"}, {"--timeout", "5"}}
	resolved := w.mustResolve()
	if setting := resolved.Get(LoggingLevel); setting.Value != "debug" || setting.Layer != Invocation || setting.Source != "--log-level" {
		t.Fatalf("%+v", setting)
	}
	if resolved.TimeLimitSeconds(ReadClassLimitSeconds) != 5 || !slices.Equal(resolved.LayersSet, []Layer{Invocation}) {
		t.Fatalf("%v %v", resolved.TimeLimitSeconds(30), resolved.LayersSet)
	}
	for _, flag := range []Flag{{"--log-level", "loud"}, {"--timeout", "0"}, {"--format", "xml"}, {"--format", ""}} {
		w.request.Flags = []Flag{flag}
		w.mustReject(InvalidValue)
	}
}

func TestResolutionLeavesProcessStateAlone(t *testing.T) {
	// CFG-005, CFG-043, CFG-053: resolution writes no file, never changes the
	// working directory or the process environment, and keeps the caller's
	// environment slice, which a process-runner request would be built from.
	w := newWorld(t)
	w.write(filepath.Join(w.work, "a.env"), "AIBOX_LOGGING__LEVEL=debug\n")
	w.write(filepath.Join(w.work, "other.json"), `{"schemaVersion":"1","logging":{"file":"x/run.log"}}`)
	w.request.EnvFiles = []string{"a.env"}
	w.request.SettingsFile, w.request.SettingsFileGiven = "other.json", true
	before := snapshot(t, filepath.Dir(w.home))
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	environ := slices.Clone(os.Environ())
	request := w.request
	request.Environ = []string{"HOME=" + w.home}
	callerEnviron := slices.Clone(request.Environ)
	resolved, rejection := Load(w.fs, request).Resolve(ReadProject(w.project))
	if rejection != nil || resolved.Text(LoggingLevel) != "debug" {
		t.Fatalf("%+v %+v", rejection, resolved)
	}
	after, _ := os.Getwd()
	if after != cwd || !slices.Equal(os.Environ(), environ) || os.Getenv("AIBOX_LOGGING__LEVEL") != "" ||
		!slices.Equal(request.Environ, callerEnviron) {
		t.Fatal("resolution changed the working directory or an environment")
	}
	if got := snapshot(t, filepath.Dir(w.home)); !reflect.DeepEqual(before, got) {
		t.Fatalf("file tree changed:\n%v\n%v", before, got)
	}
}

func snapshot(t *testing.T, root string) []string {
	t.Helper()
	var entries []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		entries = append(entries, fmt.Sprintf("%s %v %d %d", path, info.Mode(), info.Size(), info.ModTime().UnixNano()))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return entries
}
