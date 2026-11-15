// Package redact removes known secrets, secret-shaped strings and the home
// directory prefix from text before it leaves the process
// (settings-and-logging@v1, redaction). Redaction cannot be disabled.
package redact

import (
	"bytes"
	"io"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// Replacement is redaction.replacement.
const Replacement = "[redacted]"

// minSecretBytes is the shortest environment value treated as a secret.
const minSecretBytes = 8

// maxBufferedBytes bounds a Writer's line buffer; a longer line is redacted
// and written in pieces.
const maxBufferedBytes = 64 * 1024

// SecretNamePattern is redaction.secretNamePattern.
const SecretNamePattern = `(?i)(token|secret|password|passwd|credential|api_?key|private_?key)`

// SecretShapes is redaction.secretShapes.
var SecretShapes = []string{
	`gh[pousr]_[A-Za-z0-9]{36,}`,
	`github_pat_[A-Za-z0-9_]{22,}`,
	`sk-[A-Za-z0-9_-]{20,}`,
	`AKIA[0-9A-Z]{16}`,
	`-----BEGIN [A-Z ]*PRIVATE KEY-----`,
	`(?i)bearer [A-Za-z0-9._~+/-]{16,}=*`,
}

var secretName = regexp.MustCompile(SecretNamePattern)
var shapes = compileShapes()

func compileShapes() []*regexp.Regexp {
	compiled := make([]*regexp.Regexp, len(SecretShapes))
	for index, shape := range SecretShapes {
		compiled[index] = regexp.MustCompile(shape)
	}
	return compiled
}

// Redactor holds the known secrets and home prefixes of one process.
type Redactor struct {
	secrets []string
	homes   []string
}

// New collects the known secrets from the process environment (RED-001,
// D10) and the home directory, in its given and symlink-resolved forms.
func New(environ []string, home string) *Redactor {
	redactor := &Redactor{}
	for _, entry := range environ {
		name, value, found := strings.Cut(entry, "=")
		if found && len(value) >= minSecretBytes && secretName.MatchString(name) {
			redactor.secrets = append(redactor.secrets, value)
		}
	}
	// Longer secrets first, so that a secret containing another is replaced whole.
	sort.SliceStable(redactor.secrets, func(i, j int) bool { return len(redactor.secrets[i]) > len(redactor.secrets[j]) })
	if home != "" && home != "/" && filepath.IsAbs(home) {
		home = filepath.Clean(home)
		redactor.homes = append(redactor.homes, home)
		if resolved, err := filepath.EvalSymlinks(home); err == nil && resolved != home && resolved != "/" {
			redactor.homes = append(redactor.homes, resolved)
		}
	}
	return redactor
}

// Secrets replaces known secrets and secret shapes (RED-001, RED-002).
func (r *Redactor) Secrets(text string) string {
	if r != nil {
		for _, secret := range r.secrets {
			text = strings.ReplaceAll(text, secret, Replacement)
		}
	}
	for _, shape := range shapes {
		text = shape.ReplaceAllString(text, Replacement)
	}
	return text
}

// Home writes the home directory prefix of every path in text as ~ (RED-003).
func (r *Redactor) Home(text string) string {
	if r == nil {
		return text
	}
	for _, home := range r.homes {
		text = replaceHome(text, home)
	}
	return text
}

// All applies Home and then Secrets, so that a secret inside a path under the
// home directory is still found.
func (r *Redactor) All(text string) string {
	return r.Secrets(r.Home(text))
}

const pathBoundaries = " \t\r\n\"'=:,;()[]{}<>"

// replaceHome replaces home where it starts a path: at the start of text or
// after a boundary character, and followed by a slash, a boundary or the end.
func replaceHome(text, home string) string {
	var out strings.Builder
	for {
		index := strings.Index(text, home)
		if index < 0 {
			out.WriteString(text)
			return out.String()
		}
		end := index + len(home)
		before := index == 0 || strings.IndexByte(pathBoundaries, text[index-1]) >= 0
		after := end == len(text) || text[end] == '/' || strings.IndexByte(pathBoundaries, text[end]) >= 0
		out.WriteString(text[:index])
		if before && after {
			out.WriteString("~")
		} else {
			out.WriteString(home)
		}
		text = text[end:]
	}
}

// Writer redacts a stream of child output whose secrets may be split over
// several writes (RED-004). It buffers up to a line ending and redacts whole
// lines; Close writes the rest.
type Writer struct {
	mu       sync.Mutex
	redactor *Redactor
	out      io.Writer
	buffer   []byte
}

// NewWriter returns a Writer that writes redacted lines to out.
func (r *Redactor) NewWriter(out io.Writer) *Writer {
	return &Writer{redactor: r, out: out}
}

// Write accepts any fragment and writes every completed line, redacted.
func (w *Writer) Write(fragment []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buffer = append(w.buffer, fragment...)
	for {
		index := bytes.IndexByte(w.buffer, '\n')
		if index < 0 {
			break
		}
		if err := w.emit(w.buffer[:index+1]); err != nil {
			return len(fragment), err
		}
		w.buffer = w.buffer[index+1:]
	}
	if len(w.buffer) > maxBufferedBytes {
		if err := w.emit(w.buffer); err != nil {
			return len(fragment), err
		}
		w.buffer = nil
	}
	return len(fragment), nil
}

// Close writes the buffered rest of the stream, redacted.
func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.buffer) == 0 {
		return nil
	}
	err := w.emit(w.buffer)
	w.buffer = nil
	return err
}

func (w *Writer) emit(chunk []byte) error {
	_, err := io.WriteString(w.out, w.redactor.Secrets(string(chunk)))
	return err
}
