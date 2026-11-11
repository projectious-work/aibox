package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCanonicalRootAndConfinedFile(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "work")
	outside := filepath.Join(parent, "workspace-sibling")
	for _, dir := range []string{root, outside} {
		if err := os.Mkdir(dir, 0700); err != nil {
			t.Fatal(err)
		}
	}
	insideFile := filepath.Join(root, "devcontainer.json")
	outsideFile := filepath.Join(outside, "devcontainer.json")
	for _, file := range []string{insideFile, outsideFile} {
		if err := os.WriteFile(file, []byte("{}"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	rootAlias := filepath.Join(parent, "alias")
	if err := os.Symlink(root, rootAlias); err != nil {
		t.Fatal(err)
	}
	canonical, err := CanonicalRoot(rootAlias)
	if err != nil || canonical != root {
		t.Fatalf("root alias: %q, %v", canonical, err)
	}
	resolved, err := ConfinedRegularFile(canonical, filepath.Join(rootAlias, "devcontainer.json"))
	if err != nil || resolved != insideFile {
		t.Fatalf("inside file: %q, %v", resolved, err)
	}
	for _, candidate := range []string{
		outsideFile,
		filepath.Join(root, "..", "workspace-sibling", "devcontainer.json"),
		root,
	} {
		if _, err := ConfinedRegularFile(canonical, candidate); err == nil {
			t.Errorf("accepted outside or non-file path %q", candidate)
		}
	}
	escape := filepath.Join(root, "escape.json")
	if err := os.Symlink(outsideFile, escape); err != nil {
		t.Fatal(err)
	}
	if _, err := ConfinedRegularFile(canonical, escape); err == nil {
		t.Fatal("accepted symlink escape")
	}
}

func TestRejectRelativePaths(t *testing.T) {
	if _, err := CanonicalRoot("."); err == nil {
		t.Fatal("accepted relative root")
	}
	if _, err := ConfinedRegularFile(t.TempDir(), "devcontainer.json"); err == nil {
		t.Fatal("accepted relative file")
	}
}
