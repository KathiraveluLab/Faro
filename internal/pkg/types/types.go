package types

import "time"

// Record represents a normalized medical record for duplicate detection.
type Record struct {
	ID         string            `json:"id"`
	Source     string            `json:"source"`
	Timestamp  time.Time         `json:"timestamp"`
	Attributes map[string]string `json:"attributes"` // Normalized textual attributes
	BinaryRef  string            `json:"binary_ref"` // Reference to images/binary data
}

// SimilarityResult captures the outcome of a comparison between two records.
type SimilarityResult struct {
	RecordA    string  `json:"record_a"`
	RecordB    string  `json:"record_b"`
	Score      float64 `json:"score"`      // 0.0 to 1.0
	Algorithm  string  `json:"algorithm"`
	IsDuplicate bool    `json:"is_duplicate"`
	Metadata    string  `json:"metadata,omitempty"` // Raw JSON metadata for images
}
