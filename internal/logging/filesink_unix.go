//go:build linux || darwin

package logging

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// FileSink appends events to the log file by logging.fileSink.append: per
// event it takes an exclusive flock on <file>.lock, rotates if needed, opens
// the file by path with O_APPEND, writes the line with one write call and
// closes it. Only the lock file's descriptor is kept between events.
type FileSink struct {
	config     FileSinkConfig
	lock       *os.File
	retainDone bool
}

// OpenFileSink prepares the sink: it creates missing directories with mode
// 0700 (LOG-015), opens the lock file and checks that the log file can be
// opened and is private (LOG-020). On failure it returns open_failed or
// insecure_mode; an existing insecure file is left unchanged.
func OpenFileSink(config FileSinkConfig) (*FileSink, string) {
	if config.Now == nil {
		config.Now = time.Now
	}
	if err := os.MkdirAll(filepath.Dir(config.Path), 0o700); err != nil {
		return nil, OpenFailed
	}
	if info, err := os.Lstat(config.Path); err == nil {
		if !info.Mode().IsRegular() {
			return nil, OpenFailed
		}
		if info.Mode().Perm()&0o077 != 0 {
			return nil, InsecureMode
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, OpenFailed
	}
	fd, err := syscall.Open(config.Path+".lock", syscall.O_CREAT|syscall.O_RDWR|syscall.O_CLOEXEC|syscall.O_NOFOLLOW, 0o600)
	if err != nil {
		return nil, OpenFailed
	}
	sink := &FileSink{config: config, lock: os.NewFile(uintptr(fd), config.Path+".lock")}
	if err := sink.locked(sink.probe); err != nil {
		_ = sink.lock.Close()
		if errors.Is(err, ErrInsecureMode) {
			return nil, InsecureMode
		}
		return nil, OpenFailed
	}
	return sink, ""
}

// probe opens (and if needed creates) the log file without writing.
func (s *FileSink) probe() error {
	file, err := s.open()
	if err != nil {
		return err
	}
	return file.Close()
}

// Append writes one line under the lock (LOG-016 to LOG-018).
func (s *FileSink) Append(line []byte) error {
	return s.locked(func() error {
		if !s.retainDone {
			s.retainDone = true
			s.retain()
		}
		if info, err := os.Lstat(s.config.Path); err == nil && info.Size() >= int64(s.config.RotationMiB)*1024*1024 {
			if err := s.rotate(); err != nil {
				return err
			}
		}
		file, err := s.open()
		if err != nil {
			return err
		}
		written, err := syscall.Write(int(file.Fd()), line)
		closeErr := file.Close()
		if err != nil || written != len(line) {
			return fmt.Errorf("write the log file: %v", err)
		}
		return closeErr
	})
}

// Close releases the lock file descriptor.
func (s *FileSink) Close() error {
	return s.lock.Close()
}

func (s *FileSink) open() (*os.File, error) {
	fd, err := syscall.Open(s.config.Path, syscall.O_WRONLY|syscall.O_APPEND|syscall.O_CREAT|syscall.O_CLOEXEC|syscall.O_NOFOLLOW, 0o600)
	if err != nil {
		return nil, err
	}
	file := os.NewFile(uintptr(fd), s.config.Path)
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() {
		_ = file.Close()
		return nil, errors.New("the log file is not a regular file")
	}
	if info.Mode().Perm()&0o077 != 0 {
		_ = file.Close()
		return nil, ErrInsecureMode
	}
	return file, nil
}

// locked runs action while holding the exclusive lock, waiting at most
// LockWait for it.
func (s *FileSink) locked(action func() error) error {
	deadline := time.Now().Add(s.config.LockWait)
	pause := time.Millisecond
	for {
		err := syscall.Flock(int(s.lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			break
		}
		if !errors.Is(err, syscall.EWOULDBLOCK) && !errors.Is(err, syscall.EAGAIN) && !errors.Is(err, syscall.EINTR) {
			return err
		}
		if time.Now().After(deadline) {
			return errLockTimeout
		}
		time.Sleep(pause)
		pause = min(2*pause, 10*time.Millisecond)
	}
	defer syscall.Flock(int(s.lock.Fd()), syscall.LOCK_UN)
	return action()
}

func (s *FileSink) rotated(n int) string {
	return fmt.Sprintf("%s.%d", s.config.Path, n)
}

// rotate follows logging.fileSink.rotation inside the lock.
func (s *FileSink) rotate() error {
	if s.config.RetentionFiles <= 1 {
		return os.Truncate(s.config.Path, 0)
	}
	if err := removeIfPresent(s.rotated(s.config.RetentionFiles - 1)); err != nil {
		return err
	}
	for n := maxRotated - 1; n >= 1; n-- {
		if _, err := os.Lstat(s.rotated(n)); err == nil {
			if err := os.Rename(s.rotated(n), s.rotated(n+1)); err != nil {
				return err
			}
		}
	}
	if err := os.Rename(s.config.Path, s.rotated(1)); err != nil {
		return err
	}
	s.retain()
	return nil
}

// retain deletes rotated files beyond RetentionFiles and rotated files
// older than RetentionDays (LOG-017). It touches only <file>.<n>.
func (s *FileSink) retain() {
	cutoff := s.config.Now().Add(-time.Duration(s.config.RetentionDays) * 24 * time.Hour)
	for n := 1; n <= maxRotated; n++ {
		info, err := os.Lstat(s.rotated(n))
		if err != nil {
			continue
		}
		if n >= s.config.RetentionFiles || info.ModTime().Before(cutoff) {
			_ = os.Remove(s.rotated(n))
		}
	}
}

func removeIfPresent(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}
