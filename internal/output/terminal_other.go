//go:build !linux && !darwin

package output

// IsTerminal reports false on targets without a terminal check; colour then
// needs output.color always.
func IsTerminal(uintptr) bool { return false }
