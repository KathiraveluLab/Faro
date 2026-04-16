package storage

import (
	"encoding/json"
	"faro/internal/pkg/types"
	"fmt"
	"github.com/dgraph-io/badger/v4"
)

// BadgerStore implements the Store interface using BadgerDB for persistence.
type BadgerStore struct {
	db *badger.DB
}

// NewBadgerStore initializes a new BadgerDB at the specified path.
func NewBadgerStore(path string) (*BadgerStore, error) {
	opts := badger.DefaultOptions(path).WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger db: %w", err)
	}
	return &BadgerStore{db: db}, nil
}

func (s *BadgerStore) PutRecord(record types.Record) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("rec:"+record.ID), data)
	})
}

func (s *BadgerStore) GetRecord(id string) (types.Record, error) {
	var record types.Record
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("rec:" + id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &record)
		})
	})
	if err != nil {
		return types.Record{}, err
	}
	return record, nil
}

func (s *BadgerStore) PutDuplicate(result types.SimilarityResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("dup:%s:%s", result.RecordA, result.RecordB)
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

func (s *BadgerStore) ListRecords() ([]types.Record, error) {
	var records []types.Record
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte("rec:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var r types.Record
				if err := json.Unmarshal(val, &r); err != nil {
					return err
				}
				records = append(records, r)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return records, err
}

func (s *BadgerStore) Close() error {
	return s.db.Close()
}
