//go:build !linux && !darwin

package process

import "errors"

// ExecutableIdentity is unavailable until this platform has a reviewed port.
type ExecutableIdentity struct{ Path, Digest string }

// ResolveExecutable fails closed on unsupported platforms.
func ResolveExecutable(string, string) (ExecutableIdentity, error) {
	return ExecutableIdentity{}, errors.New("operator executable resolution is unsupported on this platform")
}
