package usecase

import (
	"strings"

	"github.com/projectious-work/aibox/internal/contract"
	"github.com/projectious-work/aibox/internal/redact"
	"github.com/projectious-work/aibox/internal/settings"
)

// EffectiveConfiguration builds the settings entries and settingsSources of
// the view (CFG-060 to CFG-064). Values, sources and paths are redacted,
// with the home prefix written as ~; an entry is marked redacted only when
// its value holds the replacement. An unset execution.timeoutSeconds is
// reported as the class limit with layer built_in (D5, EXE-003). The view
// is built from the resolved settings only: no hook, container or engine
// is involved (CFG-066).
func EffectiveConfiguration(resolved *settings.Resolved, redactor *redact.Redactor, classLimit int) ([]contract.ConfigurationEntry, []contract.SettingsSource) {
	entries := make([]contract.ConfigurationEntry, 0, len(resolved.Settings))
	for _, setting := range resolved.Settings {
		value := setting.Value
		if setting.Key.Name == settings.ExecutionTimeout && value == nil {
			value = classLimit
		}
		value, redacted := redactValue(redactor, value)
		entry := contract.ConfigurationEntry{
			Key: setting.Key.Name, Value: value, Redacted: redacted, Layer: string(setting.Layer),
			Source: redactor.All(setting.Source), Domain: setting.Key.Domain, InvocationOnly: setting.Layer == settings.Invocation,
		}
		for _, lower := range setting.Overridden {
			lowerValue, lowerRedacted := redactValue(redactor, lower.Value)
			entry.Overridden = append(entry.Overridden, contract.OverriddenValue{
				Layer: string(lower.Layer), Source: redactor.All(lower.Source), Value: lowerValue, Redacted: lowerRedacted,
			})
		}
		entries = append(entries, entry)
	}
	sources := make([]contract.SettingsSource, 0, len(resolved.Sources))
	for _, source := range resolved.Sources {
		sources = append(sources, contract.SettingsSource{Layer: string(source.Layer), Path: redactor.All(source.Path), State: source.State})
	}
	return entries, sources
}

func redactValue(redactor *redact.Redactor, value any) (any, bool) {
	text, ok := value.(string)
	if !ok {
		return value, false
	}
	text = redactor.All(text)
	return text, strings.Contains(text, redact.Replacement)
}

// redactNative applies RED-001 and RED-002 to the native entries, which keep
// their project-relative form.
func redactNative(redactor *redact.Redactor, entries []contract.ConfigurationEntry) []contract.ConfigurationEntry {
	for index, entry := range entries {
		switch value := entry.Value.(type) {
		case string:
			entries[index].Value = redactor.Secrets(value)
		case []string:
			redacted := make([]string, len(value))
			for i, item := range value {
				redacted[i] = redactor.Secrets(item)
			}
			entries[index].Value = redacted
		}
		entries[index].Redacted = strings.Contains(strings.Join(stringsOf(entries[index].Value), "\n"), redact.Replacement)
	}
	return entries
}

func stringsOf(value any) []string {
	switch typed := value.(type) {
	case string:
		return []string{typed}
	case []string:
		return typed
	}
	return nil
}
