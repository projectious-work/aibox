package process

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
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
		time.Sleep(10 * time.Second)
	}
	os.Exit(0)
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
