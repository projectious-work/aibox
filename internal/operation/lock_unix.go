//go:build linux || darwin

// Package operation coordinates one operation per exact environment without
// storing a parallel desired-state model.
package operation

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/projectious-work/aibox/internal/project"
)

// LockKey identifies one environment by its canonical project, selected
// Dev Container configuration, and operator-owned runtime endpoint identity.
// Display names are deliberately not sufficient to serialize mutations.
type LockKey struct {
	ProjectRoot string
	ConfigPath  string
	Endpoint    string
}

// Lock holds a process-wide advisory lock until Release. Its file is retained
// after release because unlinking it can let a waiter and a new caller lock
// different inodes for the same environment.
type Lock struct {
	file *os.File
	once sync.Once
	err  error
}

// Acquire takes an exclusive, cancellable advisory lock in an existing
// operator-owned private directory with a trusted parent path. It checks the
// project's canonical paths before deriving the key; callers must revalidate
// inputs after acquisition. A nonprivate store is refused rather than
// silently weakened.
func Acquire(ctx context.Context, store string, key LockKey) (*Lock, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	root, err := project.CanonicalRoot(key.ProjectRoot)
	if err != nil || root != key.ProjectRoot {
		return nil, errors.New("lock project root must be canonical")
	}
	config, err := project.ConfinedRegularFile(root, key.ConfigPath)
	if err != nil || config != key.ConfigPath {
		return nil, errors.New("lock config path must be canonical and confined")
	}
	if key.Endpoint == "" || strings.ContainsRune(key.Endpoint, 0) {
		return nil, errors.New("lock endpoint identity is required")
	}
	if !filepath.IsAbs(store) || filepath.Clean(store) != store {
		return nil, errors.New("lock store must be an absolute clean path")
	}
	info, err := os.Lstat(store)
	if err != nil {
		return nil, err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok || !info.IsDir() || info.Mode().Perm()&0o077 != 0 || stat.Uid != uint32(os.Geteuid()) {
		return nil, errors.New("lock store must be a private directory owned by the operator")
	}
	// NUL separators prevent ambiguous concatenations of the three identities.
	digest := sha256.Sum256([]byte(root + "\x00" + config + "\x00" + key.Endpoint))
	name := filepath.Join(store, hex.EncodeToString(digest[:])+".lock")
	fd, err := syscall.Open(name, syscall.O_CREAT|syscall.O_RDWR|syscall.O_CLOEXEC|syscall.O_NOFOLLOW, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open operation lock: %w", err)
	}
	file := os.NewFile(uintptr(fd), name)
	if err := verifyLockFile(file); err != nil {
		_ = file.Close()
		return nil, err
	}
	for {
		err = syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			return &Lock{file: file}, nil
		}
		if !errors.Is(err, syscall.EWOULDBLOCK) && !errors.Is(err, syscall.EAGAIN) {
			_ = file.Close()
			return nil, fmt.Errorf("acquire operation lock: %w", err)
		}
		select {
		case <-ctx.Done():
			_ = file.Close()
			return nil, ctx.Err()
		case <-time.After(25 * time.Millisecond):
		}
	}
}

func verifyLockFile(file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok || !info.Mode().IsRegular() || info.Mode().Perm()&0o077 != 0 || stat.Uid != uint32(os.Geteuid()) {
		return errors.New("operation lock file must be private and operator-owned")
	}
	return nil
}

// Release unlocks and closes the descriptor. It is safe to call more than once.
func (l *Lock) Release() error {
	l.once.Do(func() {
		unlock := syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
		closeErr := l.file.Close()
		l.err = errors.Join(unlock, closeErr)
	})
	return l.err
}
