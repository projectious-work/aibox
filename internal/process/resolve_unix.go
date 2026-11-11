//go:build linux || darwin

package process

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/projectious-work/aibox/internal/project"
)

const maxExecutableBytes = 128 * 1024 * 1024

// ExecutableIdentity binds one resolved operator-configured program to the
// bytes inspected before execution. Re-resolve it after acquiring the lock;
// a policy grant cannot rely on a stale identity returned here.
type ExecutableIdentity struct {
	Path   string
	Digest string
}

// ResolveExecutable accepts only an operator-supplied absolute path outside
// the project. It resolves links, checks ownership and writable ancestors,
// verifies executable permission and hashes the opened regular file. The
// caller, not project configuration, supplies configuredPath.
func ResolveExecutable(configuredPath, projectRoot string) (ExecutableIdentity, error) {
	if !filepath.IsAbs(configuredPath) {
		return ExecutableIdentity{}, errors.New("operator executable path must be absolute")
	}
	root, err := project.CanonicalRoot(projectRoot)
	if err != nil {
		return ExecutableIdentity{}, err
	}
	resolved, err := filepath.EvalSymlinks(configuredPath)
	if err != nil {
		return ExecutableIdentity{}, fmt.Errorf("resolve operator executable: %w", err)
	}
	resolved = filepath.Clean(resolved)
	if relative, err := filepath.Rel(root, resolved); err == nil && relative != ".." && relative != "." && !hasParentPrefix(relative) {
		return ExecutableIdentity{}, errors.New("operator executable is inside the project")
	}
	if err := trustedExecutableAncestors(resolved); err != nil {
		return ExecutableIdentity{}, err
	}
	file, err := os.Open(resolved)
	if err != nil {
		return ExecutableIdentity{}, fmt.Errorf("open operator executable: %w", err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return ExecutableIdentity{}, err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok || !info.Mode().IsRegular() || info.Mode().Perm()&0o111 == 0 || info.Mode().Perm()&0o022 != 0 ||
		(stat.Uid != 0 && stat.Uid != uint32(os.Geteuid())) {
		return ExecutableIdentity{}, errors.New("operator executable is not a trusted executable file")
	}
	if info.Size() > maxExecutableBytes {
		return ExecutableIdentity{}, errors.New("operator executable exceeds identity limit")
	}
	current, err := os.Stat(resolved)
	if err != nil || !os.SameFile(info, current) {
		return ExecutableIdentity{}, errors.New("operator executable changed during resolution")
	}
	hash := sha256.New()
	read, err := io.Copy(hash, io.LimitReader(file, maxExecutableBytes+1))
	if err != nil || read > maxExecutableBytes {
		return ExecutableIdentity{}, errors.New("operator executable cannot be hashed within limit")
	}
	return ExecutableIdentity{Path: resolved, Digest: "sha256:" + hex.EncodeToString(hash.Sum(nil))}, nil
}

func hasParentPrefix(relative string) bool {
	return len(relative) >= 3 && relative[:3] == ".."+string(filepath.Separator)
}

func trustedExecutableAncestors(path string) error {
	for directory := filepath.Dir(path); ; directory = filepath.Dir(directory) {
		info, err := os.Stat(directory)
		if err != nil {
			return err
		}
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok || !info.IsDir() {
			return errors.New("operator executable parent is not a directory")
		}
		if info.Mode().Perm()&0o022 != 0 && !(info.Mode()&os.ModeSticky != 0 && stat.Uid == 0) {
			return errors.New("operator executable parent is writable by others")
		}
		if directory == filepath.Dir(directory) {
			return nil
		}
	}
}
