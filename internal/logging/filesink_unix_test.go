//go:build linux || darwin

package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

func sinkConfig(path string) FileSinkConfig {
	return FileSinkConfig{Path: path, RotationMiB: 10, RetentionFiles: 7, RetentionDays: 7, LockWait: DefaultLimits.Flush}
}

func openSink(t *testing.T, config FileSinkConfig) *FileSink {
	t.Helper()
	sink, failure := OpenFileSink(config)
	if failure != "" {
		t.Fatalf("open: %s", failure)
	}
	t.Cleanup(func() { _ = sink.Close() })
	return sink
}

func mode(t *testing.T, path string) os.FileMode {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	return info.Mode().Perm()
}

func TestFileSinkCreatesPrivateFiles(t *testing.T) {
	// LOG-012, LOG-015.
	state := filepath.Join(t.TempDir(), "state")
	path := filepath.Join(state, "aibox", "logs", "run.log")
	sink := openSink(t, sinkConfig(path))
	if err := sink.Append([]byte("one\n")); err != nil {
		t.Fatal(err)
	}
	if err := sink.Append([]byte("two\n")); err != nil {
		t.Fatal(err)
	}
	contents, _ := os.ReadFile(path)
	if string(contents) != "one\ntwo\n" || mode(t, path) != 0o600 || mode(t, path+".lock") != 0o600 ||
		mode(t, filepath.Dir(path)) != 0o700 || mode(t, filepath.Join(state, "aibox")) != 0o700 {
		t.Fatalf("contents %q modes %v %v %v", contents, mode(t, path), mode(t, filepath.Dir(path)), mode(t, filepath.Join(state, "aibox")))
	}
}

func TestFileSinkRefusesInsecureAndUnopenableFiles(t *testing.T) {
	// LOG-020; LOG-030 open failure.
	dir := t.TempDir()
	wide := filepath.Join(dir, "wide.log")
	if err := os.WriteFile(wide, []byte("keep\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(wide, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, failure := OpenFileSink(sinkConfig(wide)); failure != InsecureMode {
		t.Fatalf("failure %q", failure)
	}
	if contents, _ := os.ReadFile(wide); string(contents) != "keep\n" || mode(t, wide) != 0o644 {
		t.Fatal("insecure file changed")
	}
	if _, err := os.Stat(wide + ".lock"); err == nil {
		t.Fatal("lock file created beside an insecure file")
	}
	if os.Geteuid() == 0 {
		return
	}
	locked := filepath.Join(dir, "logs")
	if err := os.Mkdir(locked, 0o500); err != nil {
		t.Fatal(err)
	}
	if _, failure := OpenFileSink(sinkConfig(filepath.Join(locked, "new.log"))); failure != OpenFailed {
		t.Fatalf("failure %q", failure)
	}
}

// fill writes size bytes of lines whose first line is first.
func fill(t *testing.T, path, first string, size int) {
	t.Helper()
	var b bytes.Buffer
	b.WriteString(first + "\n")
	for b.Len() < size {
		b.WriteString(strings.Repeat("f", 99) + "\n")
	}
	contents := b.Bytes()[:size]
	contents[size-1] = '\n'
	if err := os.WriteFile(path, contents, 0o600); err != nil {
		t.Fatal(err)
	}
}

func firstLine(t *testing.T, path string) string {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	line, _, _ := strings.Cut(string(contents), "\n")
	return line
}

func names(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var list []string
	for _, entry := range entries {
		list = append(list, entry.Name())
	}
	sort.Strings(list)
	return list
}

// run is one process lifetime: open, append a run's events, close.
func run(t *testing.T, config FileSinkConfig, events ...string) {
	t.Helper()
	sink := openSink(t, config)
	for _, event := range events {
		if err := sink.Append([]byte(event + "\n")); err != nil {
			t.Fatal(err)
		}
	}
	_ = sink.Close()
}

func TestRotationAndRetention(t *testing.T) {
	// LOG-016, LOG-017, LOG-019, D9 (J-LOG-07 and J-LOG-09 at component level).
	dir := t.TempDir()
	path := filepath.Join(dir, "r.log")
	foreign := []string{filepath.Join(dir, "other.txt"), filepath.Join(dir, "r.log.x"), filepath.Join(dir, "r.logbook")}
	for _, name := range foreign {
		if err := os.WriteFile(name, []byte("foreign"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	config := sinkConfig(path)
	config.RotationMiB, config.RetentionFiles = 1, 3
	for _, first := range []string{"A", "B", "C", "D"} {
		fill(t, path, first, 1024*1024)
		run(t, config, "event-"+first+"-1", "event-"+first+"-2")
	}
	want := []string{"other.txt", "r.log", "r.log.1", "r.log.2", "r.log.lock", "r.log.x", "r.logbook"}
	if got := names(t, dir); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("files %v", got)
	}
	if firstLine(t, path+".1") != "D" || firstLine(t, path+".2") != "C" {
		t.Fatal("rotation order")
	}
	if contents, _ := os.ReadFile(path); string(contents) != "event-D-1\nevent-D-2\n" {
		t.Fatalf("active file %q", contents)
	}
	old := time.Now().Add(-30 * 24 * time.Hour)
	if err := os.Chtimes(path+".2", old, old); err != nil {
		t.Fatal(err)
	}
	fill(t, path, "E", 1024*1024)
	run(t, config, "event-E")
	if firstLine(t, path+".2") != "D" || firstLine(t, path+".1") != "E" {
		t.Fatal("retention by age")
	}
	for _, name := range foreign {
		if contents, _ := os.ReadFile(name); string(contents) != "foreign" {
			t.Fatalf("%s changed", name)
		}
	}

	single := t.TempDir()
	config = sinkConfig(filepath.Join(single, "r.log"))
	config.RotationMiB, config.RetentionFiles = 1, 1
	for _, first := range []string{"A", "B"} {
		fill(t, config.Path, first, 1024*1024)
		run(t, config, "event-"+first)
		if got := names(t, single); strings.Join(got, ",") != "r.log,r.log.lock" {
			t.Fatalf("files %v", got)
		}
	}
}

func TestRetentionAtFirstAppendRemovesExcessFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "r.log")
	for n := 1; n <= 5; n++ {
		if err := os.WriteFile(fmt.Sprintf("%s.%d", path, n), []byte("old"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	config := sinkConfig(path)
	config.RetentionFiles = 3
	run(t, config, "x")
	if got := names(t, dir); strings.Join(got, ",") != "r.log,r.log.1,r.log.2,r.log.lock" {
		t.Fatalf("files %v", got)
	}
}

func TestSixteenWritersShareTheLock(t *testing.T) {
	// LOG-018: 16 writers contend for the lock; all finish within the lock
	// wait limit, none reaches the failure path, and every line is whole.
	// The second case rotates once in the middle and loses no line.
	for _, prefill := range []int{0, 1024*1024 - 300} {
		dir := t.TempDir()
		path := filepath.Join(dir, "c.log")
		if prefill > 0 {
			fill(t, path, "P", prefill)
		}
		config := sinkConfig(path)
		config.RotationMiB, config.RetentionFiles = 1, 100
		var wg sync.WaitGroup
		errs := make(chan error, 64)
		for writer := range 16 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sink, failure := OpenFileSink(config)
				if failure != "" {
					errs <- fmt.Errorf("open: %s", failure)
					return
				}
				defer sink.Close()
				for event := range 4 {
					line, _ := json.Marshal(map[string]any{"writer": writer, "event": event, "pad": strings.Repeat("p", 200)})
					started := time.Now()
					if err := sink.Append(append(line, '\n')); err != nil {
						errs <- err
					}
					if waited := time.Since(started); waited > config.LockWait {
						errs <- fmt.Errorf("append waited %v", waited)
					}
				}
			}()
		}
		wg.Wait()
		close(errs)
		for err := range errs {
			t.Fatal(err)
		}
		seen := map[string]int{}
		prefilled := 0
		for _, name := range names(t, dir) {
			if name == "c.log.lock" {
				continue
			}
			contents, _ := os.ReadFile(filepath.Join(dir, name))
			for _, line := range strings.Split(strings.TrimSuffix(string(contents), "\n"), "\n") {
				if strings.HasPrefix(line, "{") {
					var event map[string]any
					if err := json.Unmarshal([]byte(line), &event); err != nil {
						t.Fatalf("partial line in %s: %q", name, line)
					}
					seen[fmt.Sprint(event["writer"], "/", event["event"])]++
				} else if line != "" {
					prefilled++
				}
			}
		}
		if len(seen) != 64 {
			t.Fatalf("prefill %d: %d distinct events", prefill, len(seen))
		}
		for key, count := range seen {
			if count != 1 {
				t.Fatalf("event %s written %d times", key, count)
			}
		}
		if prefill > 0 && prefilled == 0 {
			t.Fatal("prefilled lines lost")
		}
	}
}

func TestLockWaitIsBounded(t *testing.T) {
	// A writer that cannot get the lock within LockWait fails the append
	// instead of blocking the operation (LOG-034, N13).
	path := filepath.Join(t.TempDir(), "w.log")
	config := sinkConfig(path)
	config.LockWait = 200 * time.Millisecond
	holder := openSink(t, config)
	waiter := openSink(t, config)
	release := make(chan struct{})
	held := make(chan struct{})
	go func() {
		_ = holder.locked(func() error {
			close(held)
			<-release
			return nil
		})
	}()
	<-held
	started := time.Now()
	err := waiter.Append([]byte("x\n"))
	close(release)
	if err == nil || time.Since(started) > time.Second {
		t.Fatalf("append %v after %v", err, time.Since(started))
	}
}
