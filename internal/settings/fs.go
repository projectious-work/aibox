package settings

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// ErrNotRegular reports a selected file that is not a regular file after
// symlink resolution.
var ErrNotRegular = errors.New("not a regular file")

// ErrTooLarge reports a settings or environment file above its byte limit.
var ErrTooLarge = errors.New("file exceeds its size limit")

// FileSystem is the read-only port through which settings files and
// environment files are read. ReadFile returns fs.ErrNotExist for a missing
// file, ErrNotRegular for a file that is not regular after symlink
// resolution, and any other error for a file that cannot be read.
type FileSystem interface {
	// ReadFile returns the symlink-resolved path and the bytes of the file.
	ReadFile(path string, limit int64) (resolved string, contents []byte, err error)
}

// OSFileSystem reads the host file system. It never opens a file for writing.
type OSFileSystem struct{}

// ReadFile resolves symlinks, checks for a regular file and reads at most
// limit bytes.
func (OSFileSystem) ReadFile(path string, limit int64) (string, []byte, error) {
	if _, err := os.Lstat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return path, nil, fs.ErrNotExist
		}
		return path, nil, err
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return path, nil, fs.ErrNotExist
		}
		return path, nil, err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return resolved, nil, err
	}
	if !info.Mode().IsRegular() {
		return resolved, nil, ErrNotRegular
	}
	file, err := os.Open(resolved)
	if err != nil {
		return resolved, nil, err
	}
	defer file.Close()
	contents, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return resolved, nil, err
	}
	if int64(len(contents)) > limit {
		return resolved, nil, ErrTooLarge
	}
	return resolved, contents, nil
}
