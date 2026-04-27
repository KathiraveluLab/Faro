package storage

import "faro/internal/pkg/types"

// Store defines the interface for hybrid storage (In-memory + KV).
type Store interface {
	// PutRecord saves a record to the persistent store.
	PutRecord(record types.Record) error
	// GetRecord retrieves a record from the store.
	GetRecord(id string) (types.Record, error)
	// PutDuplicate marks two records as duplicates.
	PutDuplicate(result types.SimilarityResult) error
	// ListRecords returns a slice of records for processing.
	ListRecords() ([]types.Record, error)
	// GetDuplicates retrieves all found duplicates.
	GetDuplicates() ([]types.SimilarityResult, error)
	// Close simplifies resource management.
	Close() error
}
