package detect

import "faro/internal/pkg/types"

// SimilarityMeasure defines the interface for different comparison algorithms.
type SimilarityMeasure interface {
	Name() string
	Compare(a, b string) float64
}

// Detector orchestrates the duplicate detection process.
type Detector interface {
	// Detect finds potential duplicates for a given record against a target set.
	Detect(record types.Record, targets []types.Record) []types.SimilarityResult
}

// Normalizer prepares raw strings for comparison (e.g., unit conversion, acronym expansion).
type Normalizer interface {
    Normalize(input string) string
}
