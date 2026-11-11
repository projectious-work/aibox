//go:build !linux && !darwin

package operation

import "errors"

// OperationRecordStore is unavailable where private descriptor-relative persistence
// has not been reviewed. No process-local fallback is permitted.
type OperationRecordStore struct{}

// OpenOperationRecordStore fails closed on unsupported platforms.
func OpenOperationRecordStore(string) (*OperationRecordStore, error) {
	return nil, errors.New("durable operation records are unsupported on this platform")
}

// Close is a no-op for an unavailable store.
func (*OperationRecordStore) Close() error { return nil }

// Save fails closed on unsupported platforms.
func (*OperationRecordStore) Save(OperationRecord) error {
	return errors.New("durable operation records are unsupported on this platform")
}

// Read fails closed on unsupported platforms.
func (*OperationRecordStore) Read(string) (OperationRecord, error) {
	return OperationRecord{}, errors.New("durable operation records are unsupported on this platform")
}
