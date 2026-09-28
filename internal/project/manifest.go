package project

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const maxManifestFiles = 512
const maxManifestFileBytes = 32 * 1024 * 1024
const maxManifestTotalBytes = 128 * 1024 * 1024

var manifestDigestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

// ManifestFile records one explicitly selected, non-secret control input.
// It does not claim to cover arbitrary build-context or shell dependencies.
type ManifestFile struct {
	Path   string `json:"path"`
	Digest string `json:"digest"`
	Size   int64  `json:"size"`
}

// InputManifest is the deterministic, non-secret material bound by an input
// digest. Callers discover the relevant native files through trusted rules or
// the pinned upstream CLI; this package never invents a parallel resolver.
type InputManifest struct {
	Version          string         `json:"version"`
	ProjectRoot      string         `json:"projectRoot"`
	Files            []ManifestFile `json:"files"`
	PolicyDigest     string         `json:"policyDigest,omitempty"`
	ExecutableDigest string         `json:"executableDigest,omitempty"`
}

// HashInputManifest hashes an explicitly selected file set and optional
// operator-owned policy/executable identities. It resolves every file inside
// the canonical root, checks the opened file identity, and returns stable
// JSON-hash material. The caller re-hashes after locking and before effects.
func HashInputManifest(root string, paths []string, policyDigest, executableDigest string) (InputManifest, string, error) {
	if len(paths) == 0 || len(paths) > maxManifestFiles {
		return InputManifest{}, "", errors.New("input manifest needs a bounded file set")
	}
	for _, digest := range []string{policyDigest, executableDigest} {
		if digest != "" && !manifestDigestPattern.MatchString(digest) {
			return InputManifest{}, "", errors.New("invalid operator input digest")
		}
	}
	canonicalRoot, err := CanonicalRoot(root)
	if err != nil {
		return InputManifest{}, "", err
	}
	manifest := InputManifest{Version: "aibox.input-manifest/v1", ProjectRoot: canonicalRoot,
		Files: []ManifestFile{}, PolicyDigest: policyDigest, ExecutableDigest: executableDigest}
	seen := make(map[string]bool, len(paths))
	var total int64
	for _, relative := range paths {
		if relative == "" || filepath.IsAbs(relative) || filepath.Clean(relative) != relative ||
			relative == "." || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return InputManifest{}, "", errors.New("manifest path must be a clean relative file")
		}
		normalized := filepath.ToSlash(relative)
		if seen[normalized] {
			return InputManifest{}, "", errors.New("duplicate manifest file")
		}
		seen[normalized] = true
		resolved, err := ConfinedRegularFile(canonicalRoot, filepath.Join(canonicalRoot, relative))
		if err != nil {
			return InputManifest{}, "", err
		}
		before, err := os.Stat(resolved)
		if err != nil {
			return InputManifest{}, "", err
		}
		if before.Size() > maxManifestFileBytes {
			return InputManifest{}, "", errors.New("manifest file exceeds byte limit")
		}
		file, err := os.Open(resolved)
		if err != nil {
			return InputManifest{}, "", err
		}
		opened, statErr := file.Stat()
		if statErr != nil || !os.SameFile(before, opened) {
			_ = file.Close()
			return InputManifest{}, "", errors.New("manifest file changed during inspection")
		}
		hash := sha256.New()
		size, readErr := io.Copy(hash, io.LimitReader(file, maxManifestFileBytes+1))
		closeErr := file.Close()
		if readErr != nil || closeErr != nil || size > maxManifestFileBytes {
			return InputManifest{}, "", errors.New("manifest file cannot be hashed within limit")
		}
		total += size
		if total > maxManifestTotalBytes {
			return InputManifest{}, "", errors.New("input manifest exceeds total byte limit")
		}
		manifest.Files = append(manifest.Files, ManifestFile{Path: normalized, Digest: "sha256:" + hex.EncodeToString(hash.Sum(nil)), Size: size})
	}
	sort.Slice(manifest.Files, func(i, j int) bool { return manifest.Files[i].Path < manifest.Files[j].Path })
	encoded, err := json.Marshal(manifest)
	if err != nil {
		return InputManifest{}, "", fmt.Errorf("encode input manifest: %w", err)
	}
	digest := sha256.Sum256(encoded)
	return manifest, "sha256:" + hex.EncodeToString(digest[:]), nil
}
