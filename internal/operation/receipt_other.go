//go:build !linux && !darwin

package operation

import "errors"

// ReceiptStore is unavailable where private descriptor-relative persistence
// has not been reviewed. No process-local fallback is permitted.
type ReceiptStore struct{}

// OpenReceiptStore fails closed on unsupported platforms.
func OpenReceiptStore(string) (*ReceiptStore, error) {
	return nil, errors.New("durable operation receipts are unsupported on this platform")
}

// Close is a no-op for an unavailable store.
func (*ReceiptStore) Close() error { return nil }

// Save fails closed on unsupported platforms.
func (*ReceiptStore) Save(Receipt) error {
	return errors.New("durable operation receipts are unsupported on this platform")
}

// Read fails closed on unsupported platforms.
func (*ReceiptStore) Read(string) (Receipt, error) {
	return Receipt{}, errors.New("durable operation receipts are unsupported on this platform")
}
