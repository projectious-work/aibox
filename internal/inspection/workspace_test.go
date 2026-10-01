package inspection

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fixture creates only a native declaration; no container runtime is needed.
func fixture(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, ".devcontainer")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(dir, "devcontainer.json")
	if err := os.WriteFile(config, []byte(`{"image":"example.invalid/base:one"}`), 0600); err != nil {
		t.Fatal(err)
	}
	return root, config
}

func TestWorkspaceReadsOnlyDeclaredConfig(t *testing.T) {
	root, config := fixture(t)
	resolved, first, data, err := Workspace(root)
	if err != nil {
		t.Fatal(err)
	}
	if resolved != root || !strings.HasPrefix(first, "sha256:") || data.State != "unknown" || len(data.Resources) != 0 {
		t.Fatalf("unexpected read-only inspection: %q %q %+v", resolved, first, data)
	}
	if len(data.Configuration) != 1 || data.Configuration[0].Value != ".devcontainer/devcontainer.json" {
		t.Fatalf("missing sourced native declaration: %+v", data.Configuration)
	}
	_, again, _, err := Workspace(root)
	if err != nil || first != again {
		t.Fatalf("stable declaration changed digest: %q %q %v", first, again, err)
	}
	if err := os.WriteFile(config, []byte(`{"image":"example.invalid/base:two"}`), 0600); err != nil {
		t.Fatal(err)
	}
	_, changed, _, err := Workspace(root)
	if err != nil || changed == first {
		t.Fatalf("changed declaration retained digest: %q %q %v", first, changed, err)
	}
}

func TestWorkspaceRejectsMissingAndSymlinkedDeclaration(t *testing.T) {
	root, config := fixture(t)
	if err := os.Remove(config); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := Workspace(root); err == nil {
		t.Fatal("accepted missing native declaration")
	}
	outside := filepath.Join(t.TempDir(), "secret")
	if err := os.WriteFile(outside, []byte("secret"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, config); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := Workspace(root); err == nil {
		t.Fatal("accepted symlinked declaration")
	}
}

func TestWorkspaceReportsSelectedFeaturesWithoutOptionsOrSecrets(t *testing.T) {
	root, config := fixture(t)
	declaration := `{
		// A native project owns its Feature selections.
		"features": {
			"ghcr.io/example/hugo@sha256:abc": {"version": "0.165.0", "token": "private-value"},
			"ghcr.io/example/gh@sha256:def": {},
		},
		"containerEnv": {"PRIVATE_KEY": "another-private-value"},
	}`
	if err := os.WriteFile(config, []byte(declaration), 0600); err != nil {
		t.Fatal(err)
	}
	_, _, data, err := Workspace(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.Configuration) != 2 || data.Configuration[1].Key != "devcontainer.features" {
		t.Fatalf("missing sourced Feature selection: %+v", data.Configuration)
	}
	refs, ok := data.Configuration[1].Value.([]string)
	if !ok || len(refs) != 2 || refs[0] != "ghcr.io/example/gh@sha256:def" ||
		refs[1] != "ghcr.io/example/hugo@sha256:abc" {
		t.Fatalf("unexpected selected Feature refs: %#v", data.Configuration[1].Value)
	}
	if strings.Contains(strings.Join(refs, ","), "private-value") {
		t.Fatal("inspection exposed a Feature option")
	}
}
