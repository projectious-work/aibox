// Package process runs a previously authorized executable without a shell.
// Policy resolution belongs to the caller; this package never reads project
// configuration to choose a program or inherits the ambient environment.
package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

// DiagnosticLimit is the maximum returned diagnostic tail per child stream.
// The runner temporarily retains a few more bytes to avoid leaking a secret
// fragment where an older prefix is discarded.
const DiagnosticLimit = 256 * 1024

const terminationGrace = 5 * time.Second

const maxRedactionPatterns = 64
const maxRedactionPatternBytes = 4096

// Request describes one already-authorized child invocation. Executable and
// Directory must be absolute. Environment is the complete child environment;
// a nil slice means no inherited variables. Secrets lists known values that
// must be masked in returned diagnostics; it is never serialized or logged.
type Request struct {
	// Executable is a trusted, policy-resolved absolute program path.
	Executable string
	// Arguments are literal argv elements and are never shell-interpolated.
	Arguments []string
	// Directory is the absolute child working directory.
	Directory string
	// Environment is an explicit list of KEY=VALUE entries.
	Environment []string
	// Secrets contains known byte strings to mask in both diagnostic streams.
	Secrets []string
}

// Output contains bounded diagnostic tails and the child's observed exit
// status. Truncation flags refer to the original unredacted stream lengths.
// ExitCode is only meaningful when the child reached a process exit.
type Output struct {
	// Stdout is the redacted stdout tail, not authoritative machine output.
	Stdout string
	// Stderr is the redacted stderr tail.
	Stderr string
	// StdoutTruncated reports that stdout exceeded DiagnosticLimit.
	StdoutTruncated bool
	// StderrTruncated reports that stderr exceeded DiagnosticLimit.
	StderrTruncated bool
	// ExitCode is the child's exit status when ProcessState is available.
	ExitCode int
}

// Run captures a bounded, redacted tail of each stream. A nonzero child exit
// returns its output and an *exec.ExitError; cancellation returns ctx.Err().
// On supported Unix targets cancellation signals the entire child process
// group, then escalates after terminationGrace. Effect classification remains
// the caller's responsibility before using this for mutating lifecycle calls.
func Run(ctx context.Context, req Request) (Output, error) {
	if !filepath.IsAbs(req.Executable) || !filepath.IsAbs(req.Directory) {
		return Output{}, errors.New("executable and working directory must be absolute")
	}
	for _, entry := range req.Environment {
		if len(entry) == 0 || entry[0] == '=' || !bytes.ContainsRune([]byte(entry), '=') {
			return Output{}, fmt.Errorf("invalid environment entry")
		}
	}
	if len(req.Secrets) > maxRedactionPatterns {
		return Output{}, errors.New("too many redaction patterns")
	}
	// Retain the longest secret length minus one before the visible tail.
	// A secret crossing that discard boundary can then be masked in full.
	guard := 0
	for _, secret := range req.Secrets {
		if len(secret) > maxRedactionPatternBytes {
			return Output{}, errors.New("redaction pattern exceeds limit")
		}
		if len(secret) > guard {
			guard = len(secret)
		}
	}
	if guard > 0 {
		guard--
	}
	cmd := exec.Command(req.Executable, req.Arguments...)
	cmd.Dir = req.Directory
	// An allocated empty slice is intentional: nil would inherit os.Environ.
	cmd.Env = make([]string, len(req.Environment))
	copy(cmd.Env, req.Environment)
	stdout := newTail(DiagnosticLimit + guard)
	stderr := newTail(DiagnosticLimit + guard)
	cmd.Stdout, cmd.Stderr = stdout, stderr
	err := runCommand(ctx, cmd, terminationGrace)
	out := Output{
		Stdout:          redactTail(stdout, req.Secrets, guard),
		Stderr:          redactTail(stderr, req.Secrets, guard),
		StdoutTruncated: stdout.total > DiagnosticLimit,
		StderrTruncated: stderr.total > DiagnosticLimit,
	}
	if cmd.ProcessState != nil {
		out.ExitCode = cmd.ProcessState.ExitCode()
	}
	if ctx.Err() != nil {
		return out, ctx.Err()
	}
	return out, err
}

type tail struct {
	data  []byte
	limit int
	total uint64
}

// newTail bounds retained bytes regardless of how much the child writes.
func newTail(limit int) *tail { return &tail{limit: limit} }

// Write implements io.Writer while retaining only the most recent bytes.
func (t *tail) Write(p []byte) (int, error) {
	n := len(p)
	t.total += uint64(n)
	if n >= t.limit {
		t.data = append(t.data[:0], p[n-t.limit:]...)
		return n, nil
	}
	excess := len(t.data) + n - t.limit
	if excess > 0 {
		copy(t.data, t.data[excess:])
		t.data = t.data[:len(t.data)-excess]
	}
	t.data = append(t.data, p...)
	return n, nil
}

// redactTail masks every occurrence before removing the guard prefix.
// Fixed-width masking preserves byte offsets, so a secret overlapping the
// discarded prefix cannot reappear as an unmatched suffix in the output.
func redactTail(t *tail, secrets []string, guard int) string {
	data := bytes.Clone(t.data)
	masked := make([]bool, len(data))
	for _, secret := range secrets {
		if secret == "" {
			continue
		}
		needle := []byte(secret)
		for offset := 0; offset < len(data); {
			index := bytes.Index(data[offset:], needle)
			if index < 0 {
				break
			}
			start := offset + index
			for i := start; i < start+len(needle); i++ {
				masked[i] = true
			}
			offset = start + 1
		}
	}
	for i, hide := range masked {
		if hide {
			data[i] = '*'
		}
	}
	if t.total > uint64(len(t.data)) && guard > 0 {
		data = data[min(guard, len(data)):]
	}
	if len(data) > DiagnosticLimit {
		data = data[len(data)-DiagnosticLimit:]
	}
	return string(data)
}
