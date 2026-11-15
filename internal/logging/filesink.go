package logging

import (
	"errors"
	"time"
)

// ErrInsecureMode reports a log file with a group or other permission bit
// (LOG-020).
var ErrInsecureMode = errors.New("the log file has group or other permissions")

// errLockTimeout reports a lock that was not free within the wait limit.
var errLockTimeout = errors.New("the log file lock was not free in time")

// The open failure reasons of logging.sink_degraded.
const (
	OpenFailed   = "open_failed"
	InsecureMode = "insecure_mode"
)

// FileSinkConfig is the resolved configuration of the file sink.
type FileSinkConfig struct {
	// Path is the confined, symlink-resolved absolute log file.
	Path string
	// RotationMiB is logging.rotationMiB.
	RotationMiB int
	// RetentionFiles is logging.retentionFiles, including the active file.
	RetentionFiles int
	// RetentionDays is logging.retentionDays.
	RetentionDays int
	// LockWait bounds the wait for the lock; past it the append fails.
	LockWait time.Duration
	// Now is the clock used for retention by age.
	Now func() time.Time
}

// maxRotated is the highest n of the rotation set <file>.<n>.
const maxRotated = 99
