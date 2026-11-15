package redact

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const canary = "canary-0123456789"

func TestKnownSecretsFromTheEnvironment(t *testing.T) {
	// RED-001, D10: values of secret-named variables of at least 8 bytes.
	redactor := New([]string{"EXAMPLE_API_TOKEN=" + canary, "MY_PASSWORD=short", "PLAIN=" + "not-a-secret-value",
		"GitHub_ApiKey=abcdefgh12", "DB_PASSWD=abcdefgh", "PRIVATEKEY=0123456789ab"}, "")
	got := redactor.Secrets("a " + canary + " short not-a-secret-value abcdefgh12 abcdefgh 0123456789ab")
	if got != "a [redacted] short not-a-secret-value [redacted] [redacted] [redacted]" {
		t.Fatalf("got %q", got)
	}
}

func TestSecretShapes(t *testing.T) {
	// RED-002 for each shape of redaction.secretShapes.
	samples := []string{
		"ghp_" + strings.Repeat("a", 36), "ghs_" + strings.Repeat("Z", 40), "github_pat_" + strings.Repeat("x_", 11),
		"sk-" + strings.Repeat("A1_-", 5), "AKIA" + strings.Repeat("Q", 16), "-----BEGIN RSA PRIVATE KEY-----",
		"-----BEGIN PRIVATE KEY-----", "Bearer " + strings.Repeat("t.", 8) + "==", "bearer " + strings.Repeat("k", 16),
	}
	var redactor *Redactor
	for _, sample := range samples {
		if got := redactor.Secrets("x " + sample + " y"); got != "x [redacted] y" {
			t.Fatalf("%q -> %q", sample, got)
		}
	}
	for _, harmless := range []string{"ghp_short", "sk-abc", "AKIA123", "bearer short"} {
		if got := redactor.Secrets(harmless); got != harmless {
			t.Fatalf("%q -> %q", harmless, got)
		}
	}
}

func TestHomePrefix(t *testing.T) {
	// RED-003: only a path prefix, never a longer sibling name.
	home := t.TempDir()
	redactor := New(nil, home)
	cases := map[string]string{
		home:                             "~",
		home + "/a/b":                    "~/a/b",
		"path=" + home + "/x":            "path=~/x",
		`"` + home + `"`:                 `"~"`,
		home + "x/y":                     home + "x/y",
		"/other" + home + "/y":           "/other" + home + "/y",
		"from " + home + "/a to " + home: "from ~/a to ~",
	}
	for input, want := range cases {
		if got := redactor.Home(input); got != want {
			t.Fatalf("%q -> %q, want %q", input, got, want)
		}
	}
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(home, link); err != nil {
		t.Fatal(err)
	}
	if got := New(nil, link).Home(home + "/a"); got != "~/a" {
		t.Fatalf("resolved home not replaced: %q", got)
	}
	if got := New(nil, "/").Home("/etc/x"); got != "/etc/x" {
		t.Fatalf("root home replaced: %q", got)
	}
}

func TestAllRedactsSecretsInsideHomePaths(t *testing.T) {
	home := t.TempDir()
	redactor := New([]string{"EXAMPLE_API_TOKEN=" + canary}, home)
	if got := redactor.All(home + "/" + canary + ".env"); got != "~/[redacted].env" {
		t.Fatalf("got %q", got)
	}
}

func TestWriterRedactsValuesSplitOverWrites(t *testing.T) {
	// RED-004: a secret split over several writes of child output.
	var out bytes.Buffer
	writer := New([]string{"SERVICE_TOKEN=" + canary}, "").NewWriter(&out)
	shape := "ghp_" + strings.Repeat("b", 36)
	stream := "one " + canary + " two\nthree " + shape + "\nlast " + canary
	for _, fragment := range splitEvery(stream, 3) {
		if _, err := writer.Write([]byte(fragment)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if got := out.String(); got != "one [redacted] two\nthree [redacted]\nlast [redacted]" {
		t.Fatalf("got %q", got)
	}
}

func splitEvery(text string, size int) []string {
	var parts []string
	for len(text) > size {
		parts = append(parts, text[:size])
		text = text[size:]
	}
	return append(parts, text)
}

func TestShapesMatchTheContract(t *testing.T) {
	// Drift test against interfaces.json, which go:embed cannot reach.
	contents, err := os.ReadFile("../../spec/settings-and-logging/v1/interfaces.json")
	if err != nil {
		t.Fatal(err)
	}
	var contract struct {
		Redaction struct {
			Replacement       string   `json:"replacement"`
			SecretNamePattern string   `json:"secretNamePattern"`
			SecretShapes      []string `json:"secretShapes"`
		} `json:"redaction"`
	}
	if err := json.Unmarshal(contents, &contract); err != nil {
		t.Fatal(err)
	}
	if contract.Redaction.Replacement != Replacement || contract.Redaction.SecretNamePattern != SecretNamePattern ||
		!slices.Equal(contract.Redaction.SecretShapes, SecretShapes) {
		t.Fatalf("redaction drifted: %+v", contract.Redaction)
	}
}
