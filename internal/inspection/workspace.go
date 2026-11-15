// Package inspection implements read-only workspace discovery. It deliberately
// does not create a runtime client or infer container state from configuration.
package inspection

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/projectious-work/aibox/internal/contract"
	"github.com/projectious-work/aibox/internal/project"
	"github.com/tailscale/hujson"
)

const maxNativeDeclarationBytes = 2 * 1024 * 1024

// Configuration is a non-secret, sourced value in the inspect-workspace result.
type Configuration = contract.ConfigurationEntry

// WorkspaceData matches the closed inspect_workspace payload. Unknown means
// no runtime was consulted, not that a declared configuration is unhealthy.
type WorkspaceData = contract.WorkspaceData

// Workspace reads one native project declaration and returns its canonical
// root, input-manifest digest and schema-shaped local inspection payload.
// The digest binds only the declared config bytes, not an arbitrary build
// context or runtime state. It must not authorize a later mutation. It stops
// between its reads when ctx is done.
func Workspace(ctx context.Context, root string) (string, string, WorkspaceData, error) {
	if err := ctx.Err(); err != nil {
		return "", "", WorkspaceData{}, err
	}
	canonicalRoot, err := project.CanonicalRoot(root)
	if err != nil {
		return "", "", WorkspaceData{}, err
	}
	// This preview explicitly selects only the native entry point. Later
	// lifecycle callers add upstream-discovered control files and operator
	// identities before binding any mutation to a manifest digest.
	_, digest, err := project.HashInputManifest(canonicalRoot, []string{".devcontainer/devcontainer.json"}, "", "")
	if err != nil {
		return "", "", WorkspaceData{}, err
	}
	if err := ctx.Err(); err != nil {
		return "", "", WorkspaceData{}, err
	}
	features, err := selectedFeatureRefs(canonicalRoot)
	if err != nil {
		return "", "", WorkspaceData{}, err
	}
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
	if len(features) > 0 {
		data.Configuration = append(data.Configuration, Configuration{
			Key: "devcontainer.features", Value: features,
			Redacted: false, Layer: "project", Source: ".devcontainer/devcontainer.json",
			Domain: "native", InvocationOnly: false,
		})
	}
	return canonicalRoot, digest, data, nil
}

// selectedFeatureRefs reports only native Feature identities. It intentionally
// leaves the full option objects and all other native fields to the upstream
// Dev Container CLI, and never echoes values that might contain credentials.
func selectedFeatureRefs(root string) ([]string, error) {
	filePath, err := project.ConfinedRegularFile(root, filepath.Join(root, ".devcontainer/devcontainer.json"))
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || info.Size() > maxNativeDeclarationBytes {
		return nil, errors.New("native declaration exceeds inspection limit")
	}
	contents, err := io.ReadAll(io.LimitReader(file, maxNativeDeclarationBytes+1))
	if err != nil || len(contents) > maxNativeDeclarationBytes {
		return nil, errors.New("native declaration cannot be read within inspection limit")
	}
	standard, err := hujson.Standardize(contents)
	if err != nil {
		return nil, err
	}
	var native struct {
		Features map[string]json.RawMessage `json:"features"`
	}
	if err := json.Unmarshal(standard, &native); err != nil {
		return nil, err
	}
	refs := make([]string, 0, len(native.Features))
	for ref := range native.Features {
		refs = append(refs, ref)
	}
	sort.Strings(refs)
	return refs, nil
}
