//go:build linux || darwin

package operation

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

// OperationRecordStore holds a private directory descriptor so individual records
// cannot redirect writes through a symlink. Callers serialize updates under
// the exact-environment Lock and revalidate the request before replay.
type OperationRecordStore struct{ root *os.Root }

// OpenOperationRecordStore opens an existing, operator-owned 0700 directory. It does
// not create an authority-bearing store inside a project or follow symlinks.
func OpenOperationRecordStore(path string) (*OperationRecordStore, error) {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return nil, errors.New("record store path must be absolute and clean")
	}
	real, err := filepath.EvalSymlinks(path)
	if err != nil || real != path {
		return nil, errors.New("record store must not use symlinks")
	}
	root, err := os.OpenRoot(path)
	if err != nil {
		return nil, fmt.Errorf("open record store: %w", err)
	}
	directory, err := root.Open(".")
	if err != nil {
		_ = root.Close()
		return nil, err
	}
	if err := verifyPrivateOperationRecordFile(directory, true); err != nil {
		_ = directory.Close()
		_ = root.Close()
		return nil, err
	}
	_ = directory.Close()
	return &OperationRecordStore{root: root}, nil
}

// Close releases the store descriptor; it does not delete durable records.
func (s *OperationRecordStore) Close() error { return s.root.Close() }

// Save atomically replaces a validated record and fsyncs the file and its
// parent directory. A failure returns an error rather than claiming evidence
// was durably recorded. The caller must hold the environment lock.
func (s *OperationRecordStore) Save(record OperationRecord) error {
	data, err := marshalOperationRecord(record)
	if err != nil {
		return err
	}
	previous, err := s.Read(record.OperationID)
	if err == nil {
		if err := validateOperationRecordTransition(previous, record); err != nil {
			return err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read previous record: %w", err)
	} else if record.State != "validated" {
		return errors.New("first record state must be validated")
	}
	random := make([]byte, 12)
	if _, err := rand.Read(random); err != nil {
		return err
	}
	name := record.OperationID + ".json"
	temporary := name + "." + hex.EncodeToString(random) + ".tmp"
	file, err := s.root.OpenFile(temporary, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("create temporary record: %w", err)
	}
	committed := false
	defer func() {
		_ = file.Close()
		if !committed {
			_ = s.root.Remove(temporary)
		}
	}()
	if err := verifyPrivateOperationRecordFile(file, false); err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write record: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("sync record: %w", err)
	}
	if err := s.root.Rename(temporary, name); err != nil {
		return fmt.Errorf("replace record: %w", err)
	}
	committed = true
	directory, err := s.root.Open(".")
	if err != nil {
		return fmt.Errorf("open record directory for sync: %w", err)
	}
	defer directory.Close()
	if err := directory.Sync(); err != nil {
		return fmt.Errorf("sync record directory: %w", err)
	}
	return nil
}

// Read returns one validated record by exact operation ID. It never follows
// a symlink or accepts an unbounded, unknown-field or malformed record.
func (s *OperationRecordStore) Read(id string) (OperationRecord, error) {
	if !recordIDPattern.MatchString(id) {
		return OperationRecord{}, errors.New("invalid operation ID")
	}
	file, err := s.root.Open(id + ".json")
	if err != nil {
		return OperationRecord{}, fmt.Errorf("open record: %w", err)
	}
	defer file.Close()
	if err := verifyPrivateOperationRecordFile(file, false); err != nil {
		return OperationRecord{}, err
	}
	data, err := io.ReadAll(io.LimitReader(file, maxOperationRecordBytes+1))
	if err != nil || len(data) > maxOperationRecordBytes {
		return OperationRecord{}, errors.New("record exceeds read limit")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var record OperationRecord
	if err := decoder.Decode(&record); err != nil {
		return OperationRecord{}, fmt.Errorf("decode record: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return OperationRecord{}, errors.New("record contains trailing JSON")
	}
	if err := record.Validate(); err != nil {
		return OperationRecord{}, err
	}
	if record.OperationID != id {
		return OperationRecord{}, errors.New("record ID does not match its filename")
	}
	return record, nil
}

func verifyPrivateOperationRecordFile(file *os.File, directory bool) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok || stat.Uid != uint32(os.Geteuid()) || info.Mode().Perm()&0o077 != 0 || directory != info.IsDir() {
		return errors.New("record store or file must be private and operator-owned")
	}
	if !directory && !info.Mode().IsRegular() {
		return errors.New("record is not a regular file")
	}
	return nil
}
