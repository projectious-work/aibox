//go:build linux || darwin

package operation

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/projectious-work/aibox/internal/contract"
)

func recordFixture(t *testing.T) OperationRecord {
	t.Helper()
	now := time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)
	return OperationRecord{
		SchemaVersion: OperationRecordSchemaVersion,
		OperationID:   "operation-1", RequestFingerprint: "sha256:" + strings.Repeat("a", 64),
		Operation: contract.StartEnvironment, ProjectRoot: t.TempDir(),
		InputDigest: "sha256:" + strings.Repeat("b", 64),
		CreatedAt:   now, UpdatedAt: now, Executor: "operator", State: "validated",
		LastConfirmedStep: "validated", Resources: []Resource{},
		CompletedEffects: []contract.Effect{}, UnknownEffects: []contract.Effect{},
	}
}

func TestOperationRecordStorePersistsAndReplacesAtomically(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "operations")
	if err := os.Mkdir(directory, 0700); err != nil {
		t.Fatal(err)
	}
	store, err := OpenOperationRecordStore(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	record := recordFixture(t)
	if err := store.Save(record); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(directory, record.OperationID+".json")
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("record permissions: %v %v", info, err)
	}
	read, err := store.Read(record.OperationID)
	if err != nil || read.State != "validated" || read.RequestFingerprint != record.RequestFingerprint {
		t.Fatalf("read record: %+v %v", read, err)
	}
	record.State = "authorized"
	record.LastConfirmedStep = "authorized"
	if err := store.Save(record); err != nil {
		t.Fatal(err)
	}
	record.State = "executing"
	if err := store.Save(record); err != nil {
		t.Fatal(err)
	}
	record.State = "inspected"
	record.LastConfirmedStep = "postcondition checked"
	record.UpdatedAt = time.Date(2026, 9, 28, 12, 1, 0, 0, time.UTC).Format(time.RFC3339)
	if err := store.Save(record); err != nil {
		t.Fatal(err)
	}
	read, err = store.Read(record.OperationID)
	if err != nil || read.State != "inspected" {
		t.Fatalf("updated record: %+v %v", read, err)
	}
	if matches, err := filepath.Glob(filepath.Join(directory, "*.tmp")); err != nil || len(matches) != 0 {
		t.Fatalf("temporary records remain: %v %v", matches, err)
	}
}

func TestOperationRecordStoreRejectsIdentityChangeAndBackwardState(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "operations")
	if err := os.Mkdir(directory, 0700); err != nil {
		t.Fatal(err)
	}
	store, err := OpenOperationRecordStore(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	record := recordFixture(t)
	if err := store.Save(record); err != nil {
		t.Fatal(err)
	}
	changed := record
	changed.RequestFingerprint = "sha256:" + strings.Repeat("c", 64)
	if err := store.Save(changed); err == nil {
		t.Fatal("accepted changed request fingerprint")
	}
	changed = record
	changed.State = "executing"
	if err := store.Save(changed); err == nil {
		t.Fatal("accepted skipped authorization")
	}
	record.State = "authorized"
	if err := store.Save(record); err != nil {
		t.Fatal(err)
	}
	record.State = "validated"
	if err := store.Save(record); err == nil {
		t.Fatal("accepted backward state")
	}
}

func TestOperationRecordStoreRejectsUnsafeStorageAndRecords(t *testing.T) {
	parent := t.TempDir()
	public := filepath.Join(parent, "public")
	if err := os.Mkdir(public, 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenOperationRecordStore(public); err == nil {
		t.Fatal("accepted public record store")
	}
	link := filepath.Join(parent, "link")
	if err := os.Symlink(public, link); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenOperationRecordStore(link); err == nil {
		t.Fatal("accepted symlinked record store")
	}
	private := filepath.Join(parent, "private")
	if err := os.Mkdir(private, 0700); err != nil {
		t.Fatal(err)
	}
	store, err := OpenOperationRecordStore(private)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	record := recordFixture(t)
	record.OperationID = "../escape"
	if err := store.Save(record); err == nil {
		t.Fatal("accepted traversal ID")
	}
	if _, err := store.Read("../escape"); err == nil {
		t.Fatal("read traversal ID")
	}
	record.OperationID = "operation-1"
	record.UnknownEffects = nil
	if err := store.Save(record); err == nil {
		t.Fatal("accepted omitted effect list")
	}
	if err := os.WriteFile(filepath.Join(private, "operation-1.json"), []byte(`{"unexpected":true}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Read("operation-1"); err == nil {
		t.Fatal("accepted malformed record")
	}
	if err := os.Remove(filepath.Join(private, "operation-1.json")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(parent, "secret"), filepath.Join(private, "operation-1.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Read("operation-1"); err == nil || errors.Is(err, os.ErrNotExist) {
		t.Fatalf("symlinked record should fail explicitly: %v", err)
	}
}
