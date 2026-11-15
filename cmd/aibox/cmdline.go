package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/projectious-work/aibox/internal/settings"
)

// The commands of commands.inScope.
const (
	commandInspect  = "inspect"
	commandMCPServe = "mcp serve"
	commandVersion  = "version"
	commandHelp     = "help"
)

// flagSpec is one row of cli.flags or settings.selectors.
type flagSpec struct {
	name       string
	takesValue bool
	repeatable bool
	// setting is true when the flag sets a settings key at the invocation layer.
	setting  bool
	commands []string
}

var both = []string{commandInspect, commandMCPServe}
var inspectOnly = []string{commandInspect}

// flagTable drives the parser; a drift test compares it with interfaces.json.
var flagTable = []flagSpec{
	{name: "--format", takesValue: true, setting: true, commands: inspectOnly},
	{name: "--color", takesValue: true, setting: true, commands: inspectOnly},
	{name: "--log-level", takesValue: true, setting: true, commands: both},
	{name: "--log-format", takesValue: true, setting: true, commands: both},
	{name: "--log-file", takesValue: true, setting: true, commands: both},
	{name: "--timeout", takesValue: true, setting: true, commands: inspectOnly},
	{name: "--effective-config", commands: inspectOnly},
	{name: "--context", takesValue: true, commands: both},
	{name: "--settings-file", takesValue: true, commands: both},
	{name: "--env-file", takesValue: true, repeatable: true, commands: both},
	{name: "--project", takesValue: true, commands: inspectOnly},
}

// commandLine is a parsed command line. When usage is set the command line
// was rejected; the rejection hints are still filled in.
type commandLine struct {
	command           string
	flags             []settings.Flag
	settingsFile      string
	settingsFileGiven bool
	envFiles          []string
	project           string
	effectiveConfig   bool
	usage             *settings.Rejection

	// Hints for a rejected command line (logging.rejectionFormat, OUT-006).
	logFormatJSONL bool
	logLevelOff    bool
	envelopeFormat string
}

// parseCommandLine follows cli.syntax: --name value and --name=value; an
// unknown flag, a flag on a command that does not list it, a repeated
// non-repeatable flag, a value on --effective-config and any positional
// argument are usage_error (CFG-056). It checks no setting value: the
// resolver does that at the invocation layer, after the lower layers.
func parseCommandLine(args []string) commandLine {
	line := commandLine{project: "."}
	if len(args) == 0 {
		line.command = commandHelp
		return line
	}
	var rest []string
	switch args[0] {
	case "help", "--help", "-h":
		line.command, rest = commandHelp, args[1:]
	case "version", "--version":
		line.command, rest = commandVersion, args[1:]
	case "inspect":
		line.command, rest = commandInspect, args[1:]
	case "mcp":
		line.command = commandMCPServe
		if len(args) < 2 || args[1] != "serve" {
			line.usage = settings.Usage("the only mcp command is mcp serve")
			line.collectHints(args[1:])
			return line
		}
		rest = args[2:]
	default:
		line.command = ""
		line.usage = settings.Usage(fmt.Sprintf("%q is not available in this preview", args[0]))
		line.collectHints(args[1:])
		return line
	}
	if line.command == commandHelp || line.command == commandVersion {
		if len(rest) > 0 {
			line.usage = settings.Usage(fmt.Sprintf("%s takes no arguments", line.command))
		}
		return line
	}
	line.collectHints(rest)
	line.usage = line.parseFlags(rest)
	return line
}

func (line *commandLine) parseFlags(rest []string) *settings.Rejection {
	seen := map[string]bool{}
	for index := 0; index < len(rest); index++ {
		token := rest[index]
		if !strings.HasPrefix(token, "--") || token == "--" {
			return settings.Usage(fmt.Sprintf("positional argument %q is not allowed", token))
		}
		name, value, hasValue := strings.Cut(token, "=")
		spec, ok := lookupFlag(name, line.command)
		if !ok {
			return settings.Usage(fmt.Sprintf("%s is not a flag of %s", name, line.command))
		}
		if seen[name] && !spec.repeatable {
			return settings.Usage(fmt.Sprintf("%s is given more than once", name))
		}
		seen[name] = true
		switch {
		case !spec.takesValue && hasValue:
			return settings.Usage(fmt.Sprintf("%s takes no value", name))
		case spec.takesValue && !hasValue:
			if index+1 >= len(rest) {
				return settings.Usage(fmt.Sprintf("%s needs a value", name))
			}
			index++
			value = rest[index]
		}
		if rejection := line.apply(spec, value); rejection != nil {
			return rejection
		}
	}
	if len(line.envFiles) > settings.MaxEnvFiles {
		return settings.Usage(fmt.Sprintf("at most %d environment files may be named", settings.MaxEnvFiles))
	}
	return nil
}

func (line *commandLine) apply(spec flagSpec, value string) *settings.Rejection {
	switch {
	case spec.setting:
		line.flags = append(line.flags, settings.Flag{Name: spec.name, Value: value})
	case spec.name == "--effective-config":
		line.effectiveConfig = true
	case spec.name == "--context":
		// Unchanged from the operation core: optional, and only local.
		if value != "local" {
			return settings.Usage("only --context local is available in this preview")
		}
	case spec.name == "--settings-file":
		line.settingsFile, line.settingsFileGiven = value, true
	case spec.name == "--env-file":
		line.envFiles = append(line.envFiles, value)
	case spec.name == "--project":
		line.project = value
	}
	return nil
}

func lookupFlag(name, command string) (flagSpec, bool) {
	for _, spec := range flagTable {
		if spec.name == name && slices.Contains(spec.commands, command) {
			return spec, true
		}
	}
	return flagSpec{}, false
}

// collectHints reads --log-format jsonl, --log-level off and --format json
// or yaml from a command line that may be rejected, so the rejection can be
// written in the requested form (logging.rejectionFormat, OUT-006).
func (line *commandLine) collectHints(rest []string) {
	formatSeen := false
	for index := 0; index < len(rest); index++ {
		name, value, hasValue := strings.Cut(rest[index], "=")
		if !hasValue && index+1 < len(rest) && (name == "--log-format" || name == "--log-level" || name == "--format") {
			value = rest[index+1]
		}
		switch name {
		case "--log-format":
			line.logFormatJSONL = line.logFormatJSONL || value == "jsonl"
		case "--log-level":
			line.logLevelOff = line.logLevelOff || value == "off"
		case "--format":
			if !formatSeen && (value == "json" || value == "yaml") {
				line.envelopeFormat = value
			}
			formatSeen = true
		}
	}
}
