// Package project resolves native project paths without changing them.
// Callers must revalidate paths after acquiring an operation lock and before
// each effect; a successful lookup here is not a durable authorization grant.
package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// CanonicalRoot returns the existing, symlink-resolved project directory.
// It rejects relative paths and non-directories but does not authorize a root.
func CanonicalRoot(root string) (string, error) {
	if !filepath.IsAbs(root) {
		return "", errors.New("project root must be absolute")
	}
	real, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", fmt.Errorf("resolve project root: %w", err)
	}
	info, err := os.Stat(real)
	if err != nil {
		return "", fmt.Errorf("stat project root: %w", err)
	}
	if !info.IsDir() {
		return "", errors.New("project root is not a directory")
	}
	return filepath.Clean(real), nil
}

// ConfinedRegularFile resolves an existing absolute file path and refuses
// symlinks or traversal that escape the root. It does not open or lock the
// file, so callers must revalidate identity when operation inputs are frozen.
func ConfinedRegularFile(canonicalRoot, candidate string) (string, error) {
	if !filepath.IsAbs(canonicalRoot) || !filepath.IsAbs(candidate) {
		return "", errors.New("root and file path must be absolute")
	}
	root, err := CanonicalRoot(canonicalRoot)
	if err != nil {
		return "", err
	}
	real, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return "", fmt.Errorf("resolve project file: %w", err)
	}
	relative, err := filepath.Rel(root, real)
	if err != nil || relative == "." || relative == ".." || relativeHasParent(relative) {
		return "", errors.New("project file escapes root")
	}
	info, err := os.Stat(real)
	if err != nil {
		return "", fmt.Errorf("stat project file: %w", err)
	}
	if !info.Mode().IsRegular() {
		return "", errors.New("project file is not regular")
	}
	return filepath.Clean(real), nil
}

// relativeHasParent detects a parent segment in filepath.Rel output instead
// of relying on a string prefix, which could confuse a sibling's basename.
func relativeHasParent(relative string) bool {
	for relative != "." && relative != string(filepath.Separator) {
		part := filepath.Base(relative)
		if part == ".." {
			return true
		}
		parent := filepath.Dir(relative)
		if parent == relative {
			break
		}
		relative = parent
	}
	return false
}
