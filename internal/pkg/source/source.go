package source

import "faro/internal/pkg/types"

// Source defines the interface for medical data providers.
type Source interface {
	// FetchRecords retrieves clinical textual records from the source.
	FetchRecords() ([]types.Record, error)
	// FetchMetadata retrieves hierarchical imaging metadata for a specific patient.
	FetchMetadata(patientID string) (types.MedicalImageMetadata, error)
}
