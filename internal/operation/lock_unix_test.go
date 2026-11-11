//go:build linux || darwin

package operation

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func lockFixture(t *testing.T) (string, LockKey) {
	t.Helper()
	root := t.TempDir()
	root, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(root, "devcontainer.json")
	if err := os.WriteFile(config, []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	store := t.TempDir()
	if err := os.Chmod(store, 0o700); err != nil {
		t.Fatal(err)
	}
	return store, LockKey{ProjectRoot: root, ConfigPath: config, Endpoint: "runtime-A"}
}

func TestAcquireSerializesExactEnvironment(t *testing.T) {
	store, key := lockFixture(t)
	first, err := Acquire(context.Background(), store, key)
	if err != nil {
		t.Fatal(err)
	}
	defer first.Release()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	if _, err := Acquire(ctx, store, key); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("same environment acquired twice: %v", err)
	}
	other := key
	other.Endpoint = "runtime-B"
	independent, err := Acquire(context.Background(), store, other)
	if err != nil {
		t.Fatalf("different endpoint was blocked: %v", err)
	}
	if err := independent.Release(); err != nil {
		t.Fatal(err)
	}
	if err := first.Release(); err != nil {
		t.Fatal(err)
	}
	second, err := Acquire(context.Background(), store, key)
	if err != nil {
		t.Fatalf("released lock was not reusable: %v", err)
	}
	if err := second.Release(); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(store)
	if err != nil || len(entries) != 2 {
		t.Fatalf("lock files should be retained: %v, %d", err, len(entries))
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil || info.Mode().Perm() != 0o600 {
			t.Fatalf("lock file is not private: %v, %v", info, err)
		}
	}
}

func TestAcquireRejectsUnsafeIdentityAndStore(t *testing.T) {
	store, key := lockFixture(t)
	bad := key
	bad.Endpoint = ""
	if _, err := Acquire(context.Background(), store, bad); err == nil {
		t.Fatal("accepted an empty endpoint identity")
	}
	bad = key
	bad.ConfigPath = filepath.Join(t.TempDir(), "outside.json")
	if _, err := Acquire(context.Background(), store, bad); err == nil {
		t.Fatal("accepted an out-of-root config")
	}
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(store, link); err != nil {
		t.Fatal(err)
	}
	if _, err := Acquire(context.Background(), link, key); err == nil {
		t.Fatal("accepted a symlinked lock store")
	}
	if err := os.Chmod(store, 0o777); err != nil {
		t.Fatal(err)
	}
	if _, err := Acquire(context.Background(), store, key); err == nil {
		t.Fatal("accepted a nonprivate lock store")
	}
}
