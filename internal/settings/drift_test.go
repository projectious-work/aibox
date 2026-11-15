package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// The contract cannot be embedded (go:embed cannot reach spec/), so the
// catalog is compiled into Go. These drift tests read the frozen contract at
// test time and fail on any difference.
const contractDir = "../../spec/settings-and-logging/v1"

func readContract(t *testing.T, name string, target any) {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(contractDir, name))
	if err != nil {
		t.Fatalf("read contract %s: %v", name, err)
	}
	if err := json.Unmarshal(contents, target); err != nil {
		t.Fatalf("decode contract %s: %v", name, err)
	}
}

func TestRegisteredEnvironmentMatchesTheFixture(t *testing.T) {
	var fixture struct {
		Names []string `json:"names"`
	}
	readContract(t, "fixtures/registered-environment.json", &fixture)
	if !slices.Equal(fixture.Names, registeredEnvironment) {
		t.Fatalf("registered names drifted from the fixture:\nfixture %v\ncompiled %v", fixture.Names, registeredEnvironment)
	}
}

func TestCatalogMatchesInterfaces(t *testing.T) {
	var contract struct {
		Settings struct {
			Layers []Layer `json:"layers"`
			Keys   []struct {
				Key     string         `json:"key"`
				Type    Type           `json:"type"`
				Domain  string         `json:"domain"`
				Layers  []Layer        `json:"layers"`
				Env     string         `json:"env"`
				Flag    *string        `json:"flag"`
				Default map[string]any `json:"default"`
				Values  []string       `json:"values"`
				Minimum *int           `json:"minimum"`
				Maximum *int           `json:"maximum"`
			} `json:"keys"`
			RejectionReasons []Reason `json:"rejectionReasons"`
			EnvFile          struct {
				MaxFiles int `json:"maxFiles"`
			} `json:"envFile"`
			Files struct {
				System map[string]string `json:"system"`
			} `json:"files"`
		} `json:"settings"`
		Execution struct {
			ClassDefaults map[string]int `json:"classDefaults"`
		} `json:"execution"`
	}
	readContract(t, "interfaces.json", &contract)
	if !slices.Equal(contract.Settings.Layers, Layers) || !slices.Equal(contract.Settings.RejectionReasons, Reasons) {
		t.Fatalf("layers or reasons drifted")
	}
	if contract.Settings.EnvFile.MaxFiles != MaxEnvFiles || contract.Execution.ClassDefaults["read"] != ReadClassLimitSeconds ||
		contract.Execution.ClassDefaults["mutating"] != MutatingLimitSeconds {
		t.Fatal("limits drifted")
	}
	for goos, path := range contract.Settings.Files.System {
		if SystemFile(goos) != path {
			t.Fatalf("system file for %s drifted", goos)
		}
	}
	if len(contract.Settings.Keys) != len(Keys) {
		t.Fatalf("%d keys in the contract, %d compiled", len(contract.Settings.Keys), len(Keys))
	}
	for index, want := range contract.Settings.Keys {
		got := Keys[index]
		flag := ""
		if want.Flag != nil {
			flag = *want.Flag
		}
		if got.Name != want.Key || got.Type != want.Type || got.Domain != want.Domain || !slices.Equal(got.Layers, want.Layers) ||
			got.Env != want.Env || got.Flag != flag || !slices.Equal(got.Values, want.Values) {
			t.Fatalf("key %s drifted: %+v", want.Key, got)
		}
		if want.Minimum != nil && (got.Minimum != *want.Minimum || got.Maximum != *want.Maximum) {
			t.Fatalf("range of %s drifted", want.Key)
		}
		for mode, value := range want.Default {
			compiled := got.Default(Mode(mode))
			if number, ok := value.(float64); ok {
				value = int(number)
			}
			if !reflect.DeepEqual(compiled, value) {
				t.Fatalf("default of %s for %s drifted: %v != %v", want.Key, mode, compiled, value)
			}
		}
	}
}

func TestCatalogMatchesTheProcessSettingsSchema(t *testing.T) {
	type property struct {
		Type       string              `json:"type"`
		Enum       []string            `json:"enum"`
		Minimum    *int                `json:"minimum"`
		Maximum    *int                `json:"maximum"`
		MinLength  *int                `json:"minLength"`
		Const      string              `json:"const"`
		Properties map[string]property `json:"properties"`
	}
	var schema struct {
		Required   []string            `json:"required"`
		Properties map[string]property `json:"properties"`
	}
	readContract(t, "schemas/process-settings.schema.json", &schema)
	if !slices.Equal(schema.Required, []string{"schemaVersion"}) || schema.Properties["schemaVersion"].Const != "1" {
		t.Fatal("schemaVersion rule drifted")
	}
	count := 0
	for group, groupSchema := range schema.Properties {
		if group == "schemaVersion" {
			continue
		}
		if !isGroup(group) {
			t.Fatalf("schema group %s is not in the catalog", group)
		}
		for member, rule := range groupSchema.Properties {
			count++
			key, ok := Lookup(group + "." + member)
			if !ok {
				t.Fatalf("schema member %s.%s is not in the catalog", group, member)
			}
			switch {
			case rule.Enum != nil:
				if key.Type != Enum || !slices.Equal(rule.Enum, key.Values) {
					t.Fatalf("%s enum drifted", key.Name)
				}
			case rule.Type == "integer":
				if key.Type != Integer || *rule.Minimum != key.Minimum || *rule.Maximum != key.Maximum {
					t.Fatalf("%s range drifted", key.Name)
				}
			case rule.Type == "boolean":
				if key.Type != Boolean {
					t.Fatalf("%s type drifted", key.Name)
				}
			case rule.Type == "string":
				// workspace.theme is a string in the schema; its catalog
				// values still apply (CFG-021).
				if key.Type != Path && key.Name != WorkspaceTheme || rule.MinLength == nil || *rule.MinLength != 1 {
					t.Fatalf("%s string rule drifted", key.Name)
				}
			default:
				t.Fatalf("%s has an unexpected schema rule", key.Name)
			}
		}
	}
	if count != len(Keys) {
		t.Fatalf("schema has %d settings, catalog %d", count, len(Keys))
	}
	for _, key := range Keys {
		if !strings.Contains(key.Name, ".") {
			t.Fatalf("key %s is not dotted", key.Name)
		}
	}
}
