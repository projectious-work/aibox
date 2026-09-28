// Package inspection implements read-only workspace discovery. It deliberately
// does not create a runtime client or infer container state from configuration.
package inspection

import (
	"time"

	"github.com/projectious-work/aibox/internal/project"
)

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
	// This preview explicitly selects only the native entry point. Later
	// lifecycle callers add upstream-discovered control files and operator
	// identities before binding any mutation to a manifest digest.
	_, digest, err := project.HashInputManifest(canonicalRoot, []string{".devcontainer/devcontainer.json"}, "", "")
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
	return canonicalRoot, digest, data, nil
}
