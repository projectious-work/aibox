package logging

import (
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/projectious-work/aibox/internal/redact"
)

// Limits is logging.queue.
type Limits struct {
	// Capacity is the number of queued events.
	Capacity int
	// Flush bounds the shutdown flush and every wait on a sink.
	Flush time.Duration
	// DirectWrite bounds a warn or error event written past a full queue.
	DirectWrite time.Duration
}

// DefaultLimits are the values of logging.queue.
var DefaultLimits = Limits{Capacity: 1024, Flush: 2 * time.Second, DirectWrite: 100 * time.Millisecond}

// FileAppender is the port of the file sink: it appends one complete JSON
// line or fails.
type FileAppender interface {
	Append(line []byte) error
}

// Config configures one Logger.
type Config struct {
	// Level is logging.level.
	Level string
	// Format is logging.format of the stderr sink.
	Format Format
	// Stderr receives every emitted event.
	Stderr io.Writer
	// File is the file sink, or nil when logging.file is unset or the sink
	// failed to open; FileFailure then holds the reason.
	File FileAppender
	// FileFailure is open_failed or insecure_mode when the file sink could
	// not be opened, otherwise empty.
	FileFailure string
	// Redactor is applied to every member and field before an event is
	// queued; fields also get the home prefix written as ~.
	Redactor *redact.Redactor
	// Now is the clock.
	Now func() time.Time
	// Limits are the queue limits.
	Limits Limits
}

// queued is one encoded event on its way to the sinks.
type queued struct {
	stderr     []byte
	file       []byte
	scope      *Scope
	stderrOnly bool
}

// Logger fans events out to the sinks through a bounded queue.
type Logger struct {
	config  Config
	minimum Severity
	stderr  *stderrSink
	queue   chan queued
	done    chan struct{}

	// send guards the queue against Close: emitters hold it for reading.
	send   sync.RWMutex
	closed bool

	state       sync.Mutex
	dropped     int
	fileFailed  bool
	failure     string
	failureSent bool
}

// Start creates the logger and its writer.
func Start(config Config) *Logger {
	if config.Now == nil {
		config.Now = time.Now
	}
	if config.Limits.Capacity <= 0 {
		config.Limits = DefaultLimits
	}
	logger := &Logger{
		config: config, minimum: threshold(config.Level), stderr: &stderrSink{out: config.Stderr},
		queue: make(chan queued, config.Limits.Capacity), done: make(chan struct{}),
		fileFailed: config.FileFailure != "", failure: config.FileFailure,
	}
	go logger.write()
	return logger
}

// Scope numbers the events of one request or of the process.
type Scope struct {
	logger    *Logger
	requestID string
	mu        sync.Mutex
	next      int
}

// Scope returns the event scope of a request ID; its sequence starts at 0.
func (l *Logger) Scope(requestID string) *Scope {
	return &Scope{logger: l, requestID: requestID}
}

func (s *Scope) sequence() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	value := s.next
	s.next++
	return value
}

// Emit sends an event of this scope, if its severity passes logging.level.
// It never blocks longer than the direct-write limit.
func (s *Scope) Emit(event Event) {
	s.logger.emit(s, event, false)
}

// FileSinkFailed reports whether the file sink has failed (LOG-030).
func (l *Logger) FileSinkFailed() bool {
	l.state.Lock()
	defer l.state.Unlock()
	return l.fileFailed
}

// ReportOpenFailure emits logging.sink_degraded once when the file sink
// could not be opened at start.
func (l *Logger) ReportOpenFailure(scope *Scope) {
	l.state.Lock()
	reason, send := l.failure, l.fileFailed && !l.failureSent
	l.failureSent = l.failureSent || send
	l.state.Unlock()
	if send {
		l.emit(scope, degradedEvent(reason), true)
	}
}

func degradedEvent(reason string) Event {
	return Event{Severity: Warn, Name: "logging.sink_degraded", Component: "logging",
		Fields: map[string]string{"sink": "file", "reason": reason}}
}

func (l *Logger) prepare(scope *Scope, event Event) (Event, bool) {
	if event.Severity < l.minimum {
		return event, false
	}
	event.Timestamp = l.config.Now()
	if event.Executor == "" {
		event.Executor = Executor
	}
	event.RequestID = scope.requestID
	event.Sequence = scope.sequence()
	return l.redact(event).fit(), true
}

// redact applies RED-001 and RED-002 to every member and field and RED-003
// to fields, before the event reaches a queue or a sink.
func (l *Logger) redact(event Event) Event {
	r := l.config.Redactor
	event.Name, event.Component, event.Operation = r.Secrets(event.Name), r.Secrets(event.Component), r.Secrets(event.Operation)
	event.RequestID, event.Initiator, event.Executor = r.Secrets(event.RequestID), r.Secrets(event.Initiator), r.Secrets(event.Executor)
	event.Outcome, event.ErrorCategory = r.Secrets(event.Outcome), r.Secrets(event.ErrorCategory)
	if event.Fields != nil {
		fields := make(map[string]string, len(event.Fields))
		for name, value := range event.Fields {
			fields[r.Secrets(name)] = r.All(value)
		}
		event.Fields = fields
	}
	return event
}

func (l *Logger) encode(scope *Scope, event Event, stderrOnly bool) queued {
	item := queued{file: event.encodeJSON(), scope: scope, stderrOnly: stderrOnly}
	item.stderr = item.file
	if l.config.Format == Plain {
		item.stderr = event.encodePlain()
	}
	return item
}

func (l *Logger) emit(scope *Scope, event Event, stderrOnly bool) {
	l.send.RLock()
	defer l.send.RUnlock()
	if l.closed || event.Severity < l.minimum {
		return
	}
	l.reportDropped(scope)
	prepared, ok := l.prepare(scope, event)
	if !ok {
		return
	}
	item := l.encode(scope, prepared, stderrOnly)
	select {
	case l.queue <- item:
		return
	default:
	}
	// The queue is full (LOG-033): trace and debug are dropped; info waits
	// for room up to the direct-write limit; warn and error go to stderr
	// directly within that limit.
	switch {
	case prepared.Severity >= Warn:
		l.writeDirect(item.stderr)
	case prepared.Severity == Info:
		timer := time.NewTimer(l.config.Limits.DirectWrite)
		defer timer.Stop()
		select {
		case l.queue <- item:
		case <-timer.C:
			l.countDropped()
		}
	default:
		l.countDropped()
	}
}

func (l *Logger) countDropped() {
	l.state.Lock()
	l.dropped++
	l.state.Unlock()
}

// reportDropped queues one logging.dropped event with the count once the
// queue has room. The caller holds the send lock for reading.
func (l *Logger) reportDropped(scope *Scope) {
	if len(l.queue) >= cap(l.queue) || Warn < l.minimum {
		return
	}
	l.state.Lock()
	count := l.dropped
	l.dropped = 0
	l.state.Unlock()
	if count == 0 {
		return
	}
	prepared, _ := l.prepare(scope, Event{Severity: Warn, Name: "logging.dropped", Component: "logging",
		Fields: map[string]string{"count": strconv.Itoa(count)}})
	select {
	case l.queue <- l.encode(scope, prepared, false):
	default:
		l.writeDirect(l.encode(scope, prepared, false).stderr)
	}
}

// writeDirect writes a line to stderr past the queue, waiting at most the
// direct-write limit for a sink that blocks.
func (l *Logger) writeDirect(line []byte) {
	written := make(chan struct{})
	go func() {
		l.stderr.write(line)
		close(written)
	}()
	timer := time.NewTimer(l.config.Limits.DirectWrite)
	defer timer.Stop()
	select {
	case <-written:
	case <-timer.C:
	}
}

// write is the single consumer of the queue. A file sink failure is
// reported once, directly on stderr, with the next sequence of the request
// whose event failed.
func (l *Logger) write() {
	defer close(l.done)
	for item := range l.queue {
		l.stderr.write(item.stderr)
		if item.stderrOnly || l.config.File == nil || l.FileSinkFailed() {
			continue
		}
		if err := l.config.File.Append(item.file); err != nil {
			l.failWrite(item.scope, err)
		}
	}
}

func (l *Logger) failWrite(scope *Scope, err error) {
	reason := "write_failed"
	if err == ErrInsecureMode {
		reason = "insecure_mode"
	}
	l.state.Lock()
	send := !l.failureSent
	l.fileFailed, l.failure, l.failureSent = true, reason, true
	l.state.Unlock()
	if !send {
		return
	}
	if prepared, ok := l.prepare(scope, degradedEvent(reason)); ok {
		l.stderr.write(l.encode(scope, prepared, true).stderr)
	}
}

// Close flushes the sinks for at most the flush limit (LOG-034). If the
// flush does not complete, it writes logging.flush_incomplete directly to
// stderr. Events emitted after Close are discarded.
func (l *Logger) Close(scope *Scope) {
	l.send.RLock()
	if !l.closed {
		l.reportDropped(scope)
	}
	l.send.RUnlock()
	l.send.Lock()
	if l.closed {
		l.send.Unlock()
		return
	}
	l.closed = true
	close(l.queue)
	l.send.Unlock()
	timer := time.NewTimer(l.config.Limits.Flush)
	defer timer.Stop()
	select {
	case <-l.done:
		return
	case <-timer.C:
	}
	prepared, ok := l.prepare(scope, Event{Severity: Warn, Name: "logging.flush_incomplete", Component: "logging",
		Fields: map[string]string{"sink": "file"}})
	if ok {
		l.writeDirect(l.encode(scope, prepared, true).stderr)
	}
}

// stderrSink writes whole lines to stderr. A closed stderr or a reader that
// has gone ends stderr output silently: no event, no warning (LOG-032).
type stderrSink struct {
	mu     sync.Mutex
	out    io.Writer
	failed bool
}

func (s *stderrSink) write(line []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.failed || s.out == nil {
		return
	}
	if _, err := s.out.Write(line); err != nil {
		s.failed = true
	}
}

// WriteRejection writes one settings.rejected event directly to stderr in
// the format of logging.rejectionFormat, with sequence 0 (LOG-006).
func WriteRejection(out io.Writer, format Format, requestID string, fields map[string]string, redactor *redact.Redactor, now time.Time) {
	logger := &Logger{config: Config{Redactor: redactor, Format: format}}
	event := RejectedEvent(fields)
	event.Timestamp, event.RequestID = now, requestID
	item := logger.encode(nil, logger.redact(event).fit(), true)
	_, _ = out.Write(item.stderr)
}

// EmitStderr sends an event of this scope to stderr only, as LOG-006
// requires for settings.rejected.
func (s *Scope) EmitStderr(event Event) {
	s.logger.emit(s, event, true)
}

// RejectedEvent is the settings.rejected event with the given fields.
func RejectedEvent(fields map[string]string) Event {
	return Event{Severity: Error, Name: "settings.rejected", Component: "settings", Executor: Executor,
		Outcome: "failed", ErrorCategory: "invalid_input", Fields: fields}
}
