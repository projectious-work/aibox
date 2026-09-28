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
