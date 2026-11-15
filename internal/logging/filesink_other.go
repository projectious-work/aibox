//go:build !linux && !darwin

package logging

import "errors"

// FileSink is unavailable on targets without flock; the sink then counts as
// failed and the operation continues with log_sink_degraded.
type FileSink struct{}

// OpenFileSink always fails on this target.
func OpenFileSink(FileSinkConfig) (*FileSink, string) { return nil, OpenFailed }

// Append always fails on this target.
func (*FileSink) Append([]byte) error { return errors.New("file sink unsupported") }

// Close has nothing to release on this target.
func (*FileSink) Close() error { return nil }
