//go:build !linux && !darwin

// Package operation coordinates one operation per exact environment without
// storing a parallel desired-state model.
package operation

import (
	"context"
	"errors"
)

// LockKey identifies the project, config, and runtime endpoint to lock.
type LockKey struct{ ProjectRoot, ConfigPath, Endpoint string }

// Lock is unavailable until this platform has a reviewed advisory-lock port.
type Lock struct{}

// Acquire fails closed rather than using an unsafe process-local substitute.
func Acquire(context.Context, string, LockKey) (*Lock, error) {
	return nil, errors.New("operation locks are unsupported on this platform")
}

// Release is a no-op for an unavailable lock.
func (*Lock) Release() error { return nil }
