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
)

const DiagnosticLimit = 256 * 1024
const maxRedactionPatterns = 64
const maxRedactionPatternBytes = 4096

type Request struct {
	Executable  string
	Arguments   []string
	Directory   string
	Environment []string
	Secrets     []string
}

type Output struct {
	Stdout          string
	Stderr          string
	StdoutTruncated bool
	StderrTruncated bool
	ExitCode        int
}

// Run captures a bounded, redacted tail of each stream. A nonzero child exit
// returns its output and an *exec.ExitError; cancellation returns ctx.Err().
// Process-group termination and observed-effect classification are separate
// V1-03 work and must be added before using this for mutating lifecycle calls.
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
	cmd := exec.CommandContext(ctx, req.Executable, req.Arguments...)
	cmd.Dir = req.Directory
	cmd.Env = make([]string, len(req.Environment))
	copy(cmd.Env, req.Environment)
	stdout := newTail(DiagnosticLimit + guard)
	stderr := newTail(DiagnosticLimit + guard)
	cmd.Stdout, cmd.Stderr = stdout, stderr
	err := cmd.Run()
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

func newTail(limit int) *tail { return &tail{limit: limit} }

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
