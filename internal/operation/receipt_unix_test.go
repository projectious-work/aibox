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

func receiptFixture(t *testing.T) Receipt {
	t.Helper()
	now := time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)
	return Receipt{
		SchemaVersion: ReceiptSchemaVersion,
		OperationID:   "operation-1", RequestFingerprint: "sha256:" + strings.Repeat("a", 64),
		Operation: contract.StartEnvironment, ProjectRoot: t.TempDir(),
		InputDigest: "sha256:" + strings.Repeat("b", 64),
		CreatedAt:   now, UpdatedAt: now, Executor: "operator", State: "validated",
		LastConfirmedStep: "validated", Resources: []Resource{},
		CompletedEffects: []contract.Effect{}, UnknownEffects: []contract.Effect{},
	}
}

func TestReceiptStorePersistsAndReplacesAtomically(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "operations")
	if err := os.Mkdir(directory, 0700); err != nil {
		t.Fatal(err)
	}
	store, err := OpenReceiptStore(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	receipt := receiptFixture(t)
	if err := store.Save(receipt); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(directory, receipt.OperationID+".json")
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("receipt permissions: %v %v", info, err)
	}
	read, err := store.Read(receipt.OperationID)
	if err != nil || read.State != "validated" || read.RequestFingerprint != receipt.RequestFingerprint {
		t.Fatalf("read receipt: %+v %v", read, err)
	}
	receipt.State = "authorized"
	receipt.LastConfirmedStep = "authorized"
	if err := store.Save(receipt); err != nil {
		t.Fatal(err)
	}
	receipt.State = "executing"
	if err := store.Save(receipt); err != nil {
		t.Fatal(err)
	}
	receipt.State = "inspected"
	receipt.LastConfirmedStep = "postcondition checked"
	receipt.UpdatedAt = time.Date(2026, 9, 28, 12, 1, 0, 0, time.UTC).Format(time.RFC3339)
	if err := store.Save(receipt); err != nil {
		t.Fatal(err)
	}
	read, err = store.Read(receipt.OperationID)
	if err != nil || read.State != "inspected" {
		t.Fatalf("updated receipt: %+v %v", read, err)
	}
	if matches, err := filepath.Glob(filepath.Join(directory, "*.tmp")); err != nil || len(matches) != 0 {
		t.Fatalf("temporary receipts remain: %v %v", matches, err)
	}
}

func TestReceiptStoreRejectsIdentityChangeAndBackwardState(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "operations")
	if err := os.Mkdir(directory, 0700); err != nil {
		t.Fatal(err)
	}
	store, err := OpenReceiptStore(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	receipt := receiptFixture(t)
	if err := store.Save(receipt); err != nil {
		t.Fatal(err)
	}
	changed := receipt
	changed.RequestFingerprint = "sha256:" + strings.Repeat("c", 64)
	if err := store.Save(changed); err == nil {
		t.Fatal("accepted changed request fingerprint")
	}
	changed = receipt
	changed.State = "executing"
	if err := store.Save(changed); err == nil {
		t.Fatal("accepted skipped authorization")
	}
	receipt.State = "authorized"
	if err := store.Save(receipt); err != nil {
		t.Fatal(err)
	}
	receipt.State = "validated"
	if err := store.Save(receipt); err == nil {
		t.Fatal("accepted backward state")
	}
}

func TestReceiptStoreRejectsUnsafeStorageAndRecords(t *testing.T) {
	parent := t.TempDir()
	public := filepath.Join(parent, "public")
	if err := os.Mkdir(public, 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenReceiptStore(public); err == nil {
		t.Fatal("accepted public receipt store")
	}
	link := filepath.Join(parent, "link")
	if err := os.Symlink(public, link); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenReceiptStore(link); err == nil {
		t.Fatal("accepted symlinked receipt store")
	}
	private := filepath.Join(parent, "private")
	if err := os.Mkdir(private, 0700); err != nil {
		t.Fatal(err)
	}
	store, err := OpenReceiptStore(private)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	receipt := receiptFixture(t)
	receipt.OperationID = "../escape"
	if err := store.Save(receipt); err == nil {
		t.Fatal("accepted traversal ID")
	}
	if _, err := store.Read("../escape"); err == nil {
		t.Fatal("read traversal ID")
	}
	receipt.OperationID = "operation-1"
	receipt.UnknownEffects = nil
	if err := store.Save(receipt); err == nil {
		t.Fatal("accepted omitted effect list")
	}
	if err := os.WriteFile(filepath.Join(private, "operation-1.json"), []byte(`{"unexpected":true}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Read("operation-1"); err == nil {
		t.Fatal("accepted malformed receipt")
	}
	if err := os.Remove(filepath.Join(private, "operation-1.json")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(parent, "secret"), filepath.Join(private, "operation-1.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Read("operation-1"); err == nil || errors.Is(err, os.ErrNotExist) {
		t.Fatalf("symlinked receipt should fail explicitly: %v", err)
	}
}
