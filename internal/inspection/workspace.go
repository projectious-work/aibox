// Package inspection implements read-only workspace discovery. It deliberately
// does not create a runtime client or infer container state from configuration.
package inspection

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/projectious-work/aibox/internal/project"
)

const maxConfigBytes = 1024 * 1024

// Configuration is a non-secret, sourced value in the inspect-workspace result.
type Configuration struct {
	Key            string `json:"key"`
	Value          any    `json:"value"`
	Redacted       bool   `json:"redacted"`
	Layer          string `json:"layer"`
	Source         string `json:"source"`
	Domain         string `json:"domain"`
	InvocationOnly bool   `json:"invocationOnly"`
}

// WorkspaceData matches the closed inspect_workspace payload. Unknown means
// no runtime was consulted, not that a declared configuration is unhealthy.
type WorkspaceData struct {
	State         string          `json:"state"`
	Resources     []any           `json:"resources"`
	Configuration []Configuration `json:"configuration"`
	ObservedAt    string          `json:"observedAt"`
}

// Workspace reads one native project declaration and returns its canonical
// root, input-manifest digest and schema-shaped local inspection payload.
// The digest binds only the declared config bytes, not an arbitrary build
// context or runtime state. It must not authorize a later mutation.
func Workspace(root string) (string, string, WorkspaceData, error) {
	canonicalRoot, err := project.CanonicalRoot(root)
	if err != nil {
		return "", "", WorkspaceData{}, err
	}
	configDir := filepath.Join(canonicalRoot, ".devcontainer")
	configPath := filepath.Join(configDir, "devcontainer.json")
	// Reject symlinked declaration paths before opening. The opened file is
	// compared with the lstat identity to close the check/open race.
	for _, candidate := range []string{configDir, configPath} {
		info, statErr := os.Lstat(candidate)
		if statErr != nil {
			return "", "", WorkspaceData{}, fmt.Errorf("inspect native configuration: %w", statErr)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", "", WorkspaceData{}, errors.New("native configuration path must not be a symlink")
		}
	}
	info, err := os.Lstat(configPath)
	if err != nil || !info.Mode().IsRegular() {
		return "", "", WorkspaceData{}, errors.New("native configuration must be a regular file")
	}
	if info.Size() > maxConfigBytes {
		return "", "", WorkspaceData{}, errors.New("native configuration exceeds the inspection limit")
	}
	file, err := os.Open(configPath)
	if err != nil {
		return "", "", WorkspaceData{}, fmt.Errorf("open native configuration: %w", err)
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return "", "", WorkspaceData{}, errors.New("native configuration changed during inspection")
	}
	content, err := io.ReadAll(io.LimitReader(file, maxConfigBytes+1))
	if err != nil || len(content) > maxConfigBytes {
		return "", "", WorkspaceData{}, errors.New("native configuration cannot be read within the inspection limit")
	}
	// A stable manifest names the root, relative declaration and byte hash.
	// Go's JSON encoder sorts object keys; no secret-bearing environment is read.
	configHash := sha256.Sum256(content)
	manifest, err := json.Marshal(map[string]string{
		"workspaceRoot": canonicalRoot,
		"configPath":    ".devcontainer/devcontainer.json",
		"configSha256":  hex.EncodeToString(configHash[:]),
	})
	if err != nil {
		return "", "", WorkspaceData{}, err
	}
	digest := sha256.Sum256(manifest)
	data := WorkspaceData{
		State:     "unknown",
		Resources: []any{},
		Configuration: []Configuration{{
			Key: "devcontainer.config", Value: ".devcontainer/devcontainer.json",
			Redacted: false, Layer: "project", Source: ".devcontainer/devcontainer.json",
			Domain: "native", InvocationOnly: false,
		}},
		ObservedAt: time.Now().UTC().Format(time.RFC3339Nano),
	}
	return canonicalRoot, "sha256:" + hex.EncodeToString(digest[:]), data, nil
}
