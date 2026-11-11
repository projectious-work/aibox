//go:build !linux && !darwin

package process

import (
	"context"
	"errors"
	"os/exec"
	"time"
)

// runCommand fails closed on platforms without a reviewed process-group
// termination implementation; killing only the leader could leave helpers.
func runCommand(_ context.Context, _ *exec.Cmd, _ time.Duration) error {
	return errors.New("child process groups are unsupported on this platform")
}
