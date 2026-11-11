//go:build linux || darwin

package process

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveExecutableBindsOperatorFile(t *testing.T) {
	projectRoot := t.TempDir()
	operatorRoot := t.TempDir()
	path := filepath.Join(operatorRoot, "fake-devcontainer")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	first, err := ResolveExecutable(path, projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	if first.Path != path || !strings.HasPrefix(first.Digest, "sha256:") {
		t.Fatalf("unexpected executable identity: %+v", first)
	}
	if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 1\n"), 0700); err != nil {
		t.Fatal(err)
	}
	changed, err := ResolveExecutable(path, projectRoot)
	if err != nil || changed.Digest == first.Digest {
		t.Fatalf("changed executable retained identity: %+v %v", changed, err)
	}
}

func TestResolveExecutableRejectsProjectAndWritablePaths(t *testing.T) {
	projectRoot := t.TempDir()
	inside := filepath.Join(projectRoot, "fake-devcontainer")
	if err := os.WriteFile(inside, []byte("#!/bin/sh\n"), 0700); err != nil {
		t.Fatal(err)
	}
	if _, err := ResolveExecutable(inside, projectRoot); err == nil {
		t.Fatal("accepted repository-owned executable")
	}
	operatorRoot := t.TempDir()
	outside := filepath.Join(operatorRoot, "fake-devcontainer")
	if err := os.WriteFile(outside, []byte("#!/bin/sh\n"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(outside, 0722); err != nil {
		t.Fatal(err)
	}
	if _, err := ResolveExecutable(outside, projectRoot); err == nil {
		t.Fatal("accepted writable executable")
	}
	if err := os.Chmod(outside, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(operatorRoot, 0777); err != nil {
		t.Fatal(err)
	}
	if _, err := ResolveExecutable(outside, projectRoot); err == nil {
		t.Fatal("accepted executable below writable directory")
	}
}
