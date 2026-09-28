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

// ReceiptStore holds a private directory descriptor so individual receipts
// cannot redirect writes through a symlink. Callers serialize updates under
// the exact-environment Lock and revalidate the request before replay.
type ReceiptStore struct{ root *os.Root }

// OpenReceiptStore opens an existing, operator-owned 0700 directory. It does
// not create an authority-bearing store inside a project or follow symlinks.
func OpenReceiptStore(path string) (*ReceiptStore, error) {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return nil, errors.New("receipt store path must be absolute and clean")
	}
	real, err := filepath.EvalSymlinks(path)
	if err != nil || real != path {
		return nil, errors.New("receipt store must not use symlinks")
	}
	root, err := os.OpenRoot(path)
	if err != nil {
		return nil, fmt.Errorf("open receipt store: %w", err)
	}
	directory, err := root.Open(".")
	if err != nil {
		_ = root.Close()
		return nil, err
	}
	if err := verifyPrivateReceiptFile(directory, true); err != nil {
		_ = directory.Close()
		_ = root.Close()
		return nil, err
	}
	_ = directory.Close()
	return &ReceiptStore{root: root}, nil
}

// Close releases the store descriptor; it does not delete durable records.
func (s *ReceiptStore) Close() error { return s.root.Close() }

// Save atomically replaces a validated receipt and fsyncs the file and its
// parent directory. A failure returns an error rather than claiming evidence
// was durably recorded. The caller must hold the environment lock.
func (s *ReceiptStore) Save(receipt Receipt) error {
	data, err := marshalReceipt(receipt)
	if err != nil {
		return err
	}
	previous, err := s.Read(receipt.OperationID)
	if err == nil {
		if err := validateReceiptTransition(previous, receipt); err != nil {
			return err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read previous receipt: %w", err)
	} else if receipt.State != "validated" {
		return errors.New("first receipt state must be validated")
	}
	random := make([]byte, 12)
	if _, err := rand.Read(random); err != nil {
		return err
	}
	name := receipt.OperationID + ".json"
	temporary := name + "." + hex.EncodeToString(random) + ".tmp"
	file, err := s.root.OpenFile(temporary, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("create temporary receipt: %w", err)
	}
	committed := false
	defer func() {
		_ = file.Close()
		if !committed {
			_ = s.root.Remove(temporary)
		}
	}()
	if err := verifyPrivateReceiptFile(file, false); err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write receipt: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("sync receipt: %w", err)
	}
	if err := s.root.Rename(temporary, name); err != nil {
		return fmt.Errorf("replace receipt: %w", err)
	}
	committed = true
	directory, err := s.root.Open(".")
	if err != nil {
		return fmt.Errorf("open receipt directory for sync: %w", err)
	}
	defer directory.Close()
	if err := directory.Sync(); err != nil {
		return fmt.Errorf("sync receipt directory: %w", err)
	}
	return nil
}

// Read returns one validated receipt by exact operation ID. It never follows
// a symlink or accepts an unbounded, unknown-field or malformed record.
func (s *ReceiptStore) Read(id string) (Receipt, error) {
	if !receiptIDPattern.MatchString(id) {
		return Receipt{}, errors.New("invalid operation ID")
	}
	file, err := s.root.Open(id + ".json")
	if err != nil {
		return Receipt{}, fmt.Errorf("open receipt: %w", err)
	}
	defer file.Close()
	if err := verifyPrivateReceiptFile(file, false); err != nil {
		return Receipt{}, err
	}
	data, err := io.ReadAll(io.LimitReader(file, maxReceiptBytes+1))
	if err != nil || len(data) > maxReceiptBytes {
		return Receipt{}, errors.New("receipt exceeds read limit")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var receipt Receipt
	if err := decoder.Decode(&receipt); err != nil {
		return Receipt{}, fmt.Errorf("decode receipt: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return Receipt{}, errors.New("receipt contains trailing JSON")
	}
	if err := receipt.Validate(); err != nil {
		return Receipt{}, err
	}
	if receipt.OperationID != id {
		return Receipt{}, errors.New("receipt ID does not match its filename")
	}
	return receipt, nil
}

func verifyPrivateReceiptFile(file *os.File, directory bool) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok || stat.Uid != uint32(os.Geteuid()) || info.Mode().Perm()&0o077 != 0 || directory != info.IsDir() {
		return errors.New("receipt store or file must be private and operator-owned")
	}
	if !directory && !info.Mode().IsRegular() {
		return errors.New("receipt is not a regular file")
	}
	return nil
}
