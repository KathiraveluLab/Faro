package detect

import "math"

// Levenshtein implements the SimilarityMeasure interface using Levenshtein distance.
type Levenshtein struct{}

func (l Levenshtein) Name() string {
	return "Levenshtein"
}

func (l Levenshtein) Compare(a, b string) float64 {
	d := levenshteinDistance(a, b)
	maxLength := math.Max(float64(len(a)), float64(len(b)))
	if maxLength == 0 {
		return 1.0
	}
	// Return similarity score: 1 - (distance / max length)
	return 1.0 - (float64(d) / maxLength)
}

func levenshteinDistance(a, b string) int {
	f := make([]int, len(b)+1)
	for j := range f {
		f[j] = j
	}
	for _, ca := range a {
		j := 1
		nw := f[0]
		f[0]++
		for _, cb := range b {
			cur := f[j]
			if ca == cb {
				f[j] = nw
			} else {
				f[j] = min(nw, min(f[j], f[j-1])) + 1
			}
			nw = cur
			j++
		}
	}
	return f[len(b)]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
