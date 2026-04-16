package detect

import (
	"encoding/json"
	"faro/internal/pkg/types"
)

// MetadataDetector implements duplicate detection based on image metadata hierarchy.
type MetadataDetector struct{}

// CompareStudies checks for identical StudyInstanceUIDs between two patient records.
func (d MetadataDetector) CompareStudies(a, b types.MedicalImageMetadata) []types.SimilarityResult {
	var results []types.SimilarityResult

	for _, studyA := range a.Studies {
		for _, studyB := range b.Studies {
			if studyA.StudyInstanceUID == studyB.StudyInstanceUID {
				// We found a duplicate study! 
				// Now let's calculate the overlap at the series level.
				overlap := calculateSeriesOverlap(studyA.Series, studyB.Series)
				
				metaJSON, _ := json.Marshal(studyA)
				results = append(results, types.SimilarityResult{
					RecordA:     a.PatientID + ":" + studyA.StudyInstanceUID,
					RecordB:     b.PatientID + ":" + studyB.StudyInstanceUID,
					Score:       overlap,
					Algorithm:   "HierarchicalMetadata",
					IsDuplicate: overlap > 0.5, // 50% series overlap threshold
					Metadata:    string(metaJSON),
				})
			}
		}
	}

	return results
}

func calculateSeriesOverlap(seriesA, seriesB []types.SeriesLevel) float64 {
	if len(seriesA) == 0 || len(seriesB) == 0 {
		return 0.0
	}

	matchCount := 0
	for _, sA := range seriesA {
		for _, sB := range seriesB {
			if sA.SeriesInstanceUID == sB.SeriesInstanceUID {
				matchCount++
				break
			}
		}
	}

	// Simple Jaccard-like overlap for series
	totalUnique := len(seriesA) + len(seriesB) - matchCount
	if totalUnique == 0 {
		return 0.0
	}
	return float64(matchCount) / float64(totalUnique)
}
