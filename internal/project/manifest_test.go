package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInputManifestStableAndBoundToFileBytes(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".devcontainer"), 0700); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(root, ".devcontainer", "devcontainer.json")
	dockerfile := filepath.Join(root, ".devcontainer", "Dockerfile")
	for path, content := range map[string]string{config: `{"build":{"dockerfile":"Dockerfile"}}`, dockerfile: "FROM scratch\n"} {
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	paths := []string{".devcontainer/Dockerfile", ".devcontainer/devcontainer.json"}
	manifest, first, err := HashInputManifest(root, paths, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if manifest.Files[0].Path != ".devcontainer/Dockerfile" || !strings.HasPrefix(first, "sha256:") {
		t.Fatalf("unexpected manifest: %+v %s", manifest, first)
	}
	_, reordered, err := HashInputManifest(root, []string{paths[1], paths[0]}, "", "")
	if err != nil || reordered != first {
		t.Fatalf("file order changed digest: %s %s %v", first, reordered, err)
	}
	if err := os.WriteFile(dockerfile, []byte("FROM busybox\n"), 0600); err != nil {
		t.Fatal(err)
	}
	_, changed, err := HashInputManifest(root, paths, "", "")
	if err != nil || changed == first {
		t.Fatalf("changed control file retained digest: %s %s %v", first, changed, err)
	}
}

func TestInputManifestRejectsTraversalAndEscapedSymlinks(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "devcontainer.json")
	if err := os.WriteFile(file, []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, paths := range [][]string{{}, {"../outside"}, {"/etc/passwd"}, {"devcontainer.json", "devcontainer.json"}} {
		if _, _, err := HashInputManifest(root, paths, "", ""); err == nil {
			t.Fatalf("accepted unsafe manifest paths: %v", paths)
		}
	}
	outside := filepath.Join(t.TempDir(), "secret")
	if err := os.WriteFile(outside, []byte("secret"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "escape")); err != nil {
		t.Fatal(err)
	}
	if _, _, err := HashInputManifest(root, []string{"escape"}, "", ""); err == nil {
		t.Fatal("accepted manifest symlink outside the project")
	}
}
