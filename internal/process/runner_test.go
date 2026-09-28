package process

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestHelperProcess(t *testing.T) {
	if os.Getenv("AIBOX_PROCESS_TEST_HELPER") != "1" {
		return
	}
	switch os.Getenv("AIBOX_PROCESS_TEST_MODE") {
	case "echo":
		_, _ = os.Stdout.WriteString(strings.Join(os.Args, "|") + "\n")
		_, _ = os.Stderr.WriteString("secret-on-stderr\n")
	case "exit":
		os.Exit(17)
	case "wait":
		if ready := os.Getenv("AIBOX_PROCESS_TEST_READY"); ready != "" {
			_ = os.WriteFile(ready, []byte("ready"), 0o600)
		}
		if marker := os.Getenv("AIBOX_PROCESS_TEST_SURVIVED"); marker != "" {
			time.Sleep(600 * time.Millisecond)
			_ = os.WriteFile(marker, []byte("survived"), 0o600)
			break
		}
		time.Sleep(10 * time.Second)
	case "spawn":
		child := exec.Command(os.Args[0], "-test.run=TestHelperProcess")
		child.Env = []string{
			"AIBOX_PROCESS_TEST_HELPER=1",
			"AIBOX_PROCESS_TEST_MODE=wait",
			"AIBOX_PROCESS_TEST_READY=" + os.Getenv("AIBOX_PROCESS_TEST_READY"),
			"AIBOX_PROCESS_TEST_SURVIVED=" + os.Getenv("AIBOX_PROCESS_TEST_SURVIVED"),
		}
		child.Stdout, child.Stderr = os.Stdout, os.Stderr
		if err := child.Start(); err != nil {
			os.Exit(19)
		}
		_, _ = child.Process.Wait()
	case "ignore-term":
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM)
		if ready := os.Getenv("AIBOX_PROCESS_TEST_READY"); ready != "" {
			_ = os.WriteFile(ready, []byte("ready"), 0o600)
		}
		time.Sleep(10 * time.Second)
	}
	os.Exit(0)
}

func TestRunCancellationEscalatesToKill(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("process-group cancellation is supported on Linux and macOS")
	}
	ready := filepath.Join(t.TempDir(), "ready")
	req := helperRequest(t, "ignore-term")
	req.Environment = append(req.Environment, "AIBOX_PROCESS_TEST_READY="+ready)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.Command(req.Executable, req.Arguments...)
	cmd.Dir, cmd.Env = req.Directory, req.Environment
	done := make(chan error, 1)
	go func() { done <- runCommand(ctx, cmd, 50*time.Millisecond) }()
	deadline := time.After(2 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		select {
		case <-deadline:
			t.Fatal("TERM-resistant child did not become ready")
		case <-time.After(10 * time.Millisecond):
		}
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cancellation attribution: %v", err)
		}
		if cmd.ProcessState == nil || cmd.ProcessState.ExitCode() != -1 {
			t.Fatalf("TERM-resistant child was not killed: %v", cmd.ProcessState)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("TERM-resistant child survived escalation")
	}
}

func TestRunCancellationTerminatesDescendants(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("process-group cancellation is supported on Linux and macOS")
	}
	root := t.TempDir()
	ready := filepath.Join(root, "ready")
	survived := filepath.Join(root, "survived")
	req := helperRequest(t, "spawn")
	req.Environment = append(req.Environment,
		"AIBOX_PROCESS_TEST_READY="+ready,
		"AIBOX_PROCESS_TEST_SURVIVED="+survived,
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { _, err := Run(ctx, req); done <- err }()
	deadline := time.After(2 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		select {
		case <-deadline:
			t.Fatal("descendant did not become ready")
		case <-time.After(10 * time.Millisecond):
		}
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cancellation attribution: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("process group did not terminate promptly")
	}
	time.Sleep(700 * time.Millisecond)
	if _, err := os.Stat(survived); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("descendant survived cancellation: %v", err)
	}
}

func helperRequest(t *testing.T, mode string) Request {
	t.Helper()
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	return Request{
		Executable:  self,
		Arguments:   []string{"-test.run=TestHelperProcess", "literal$(echo unsafe)"},
		Directory:   t.TempDir(),
		Environment: []string{"AIBOX_PROCESS_TEST_HELPER=1", "AIBOX_PROCESS_TEST_MODE=" + mode},
		Secrets:     []string{"secret-on-stderr"},
	}
}

func TestRunArgumentVectorAndRedaction(t *testing.T) {
	req := helperRequest(t, "echo")
	out, err := Run(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.Stdout, "literal$(echo unsafe)") || strings.Contains(out.Stderr, "secret-on-stderr") {
		t.Fatalf("argument or redaction mismatch: %+v", out)
	}
	if out.Stderr != strings.Repeat("*", len("secret-on-stderr"))+"\n" {
		t.Fatalf("unexpected redacted diagnostic: %q", out.Stderr)
	}
}

func TestRunChildFailureAndCancellation(t *testing.T) {
	out, err := Run(context.Background(), helperRequest(t, "exit"))
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || out.ExitCode != 17 {
		t.Fatalf("lost child exit attribution: %+v, %v", out, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, err = Run(ctx, helperRequest(t, "wait"))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("cancellation not attributed to context: %v", err)
	}
}

func TestTailBoundsAndBoundaryRedaction(t *testing.T) {
	tail := newTail(DiagnosticLimit + len("secret") - 1)
	_, _ = tail.Write([]byte(strings.Repeat("x", DiagnosticLimit)))
	_, _ = tail.Write([]byte("sec"))
	_, _ = tail.Write([]byte("ret-end"))
	redacted := redactTail(tail, []string{"secret"}, len("secret")-1)
	if len(redacted) > DiagnosticLimit || strings.Contains(redacted, "secret") || !strings.HasSuffix(redacted, "******-end") {
		t.Fatalf("tail bound or cross-chunk redaction failed: len=%d suffix=%q", len(redacted), redacted[len(redacted)-20:])
	}
	if tail.total <= DiagnosticLimit {
		t.Fatal("test did not exceed output bound")
	}
}

func TestTailRedactsSecretAtDiscardBoundary(t *testing.T) {
	secret := "hidden"
	tail := newTail(DiagnosticLimit + len(secret) - 1)
	_, _ = tail.Write([]byte(secret[:3]))
	_, _ = tail.Write([]byte(secret[3:] + strings.Repeat("x", DiagnosticLimit+1)))
	out := redactTail(tail, []string{secret}, len(secret)-1)
	if len(out) != DiagnosticLimit || strings.Contains(out, secret[3:]) {
		t.Fatalf("secret fragment leaked across discarded prefix: %q", out[:20])
	}
}

func TestRunRejectsRelativePathAndAmbientEnvironment(t *testing.T) {
	_, err := Run(context.Background(), Request{Executable: "devcontainer", Directory: t.TempDir()})
	if err == nil {
		t.Fatal("accepted unresolved executable")
	}
	req := helperRequest(t, "echo")
	req.Environment = nil
	out, err := Run(context.Background(), req)
	if err != nil || out.Stdout != "PASS\n" || out.Stderr != "" {
		t.Fatalf("unexpected empty-environment execution: %+v, %v", out, err)
	}
}
