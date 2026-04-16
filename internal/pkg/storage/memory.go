package storage

import (
	"faro/internal/pkg/types"
	"fmt"
	"sync"
)

// MemoryStore is an in-memory implementation of the Store interface.
type MemoryStore struct {
	mu         sync.RWMutex
	records    map[string]types.Record
	duplicates []types.SimilarityResult
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		records: make(map[string]types.Record),
	}
}

func (m *MemoryStore) PutRecord(record types.Record) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records[record.ID] = record
	return nil
}

func (m *MemoryStore) GetRecord(id string) (types.Record, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	record, ok := m.records[id]
	if !ok {
		return types.Record{}, fmt.Errorf("record %s not found", id)
	}
	return record, nil
}

func (m *MemoryStore) PutDuplicate(result types.SimilarityResult) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.duplicates = append(m.duplicates, result)
	return nil
}

func (m *MemoryStore) ListRecords() ([]types.Record, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	records := make([]types.Record, 0, len(m.records))
	for _, r := range m.records {
		records = append(records, r)
	}
	return records, nil
}

func (m *MemoryStore) Close() error {
	return nil
}
