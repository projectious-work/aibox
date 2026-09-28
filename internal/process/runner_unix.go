//go:build linux || darwin

package process

import (
	"context"
	"errors"
	"os/exec"
	"syscall"
	"time"
)

// runCommand gives each invocation a separate process group so cancellation
// reaches helpers spawned by the delegated executable as well as its leader.
// The group is signalled only while Wait is outstanding, preventing a later
// invocation from receiving a signal for a reused process-group ID.
func runCommand(ctx context.Context, cmd *exec.Cmd, grace time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return err
	}
	waited := make(chan error, 1)
	go func() { waited <- cmd.Wait() }()
	select {
	case err := <-waited:
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	case <-ctx.Done():
		// ESRCH is harmless: the group may have exited between cancellation
		// and delivery. A different failure is still bounded by KILL below.
		_ = signalGroup(cmd.Process.Pid, syscall.SIGTERM)
	}
	timer := time.NewTimer(grace)
	defer timer.Stop()
	select {
	case <-waited:
		return ctx.Err()
	case <-timer.C:
		_ = signalGroup(cmd.Process.Pid, syscall.SIGKILL)
		<-waited
		return ctx.Err()
	}
}

func signalGroup(pid int, signal syscall.Signal) error {
	err := syscall.Kill(-pid, signal)
	if errors.Is(err, syscall.ESRCH) {
		return nil
	}
	return err
}
