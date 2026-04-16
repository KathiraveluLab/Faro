package detect

import (
	"regexp"
	"strings"
)

// ClinicalNormalizer handles medical-specific text normalization.
type ClinicalNormalizer struct {
	acronyms map[string]string
}

func NewClinicalNormalizer() *ClinicalNormalizer {
	return &ClinicalNormalizer{
		acronyms: map[string]string{
			"ASA":  "Aspirin",
			"APAP": "Acetaminophen",
			"HCTZ": "Hydrochlorothiazide",
		},
	}
}

func (n *ClinicalNormalizer) Normalize(input string) string {
	// 1. Convert to uppercase for acronym matching, then process
	upper := strings.ToUpper(strings.TrimSpace(input))

	// 2. Expand acronyms
	res := upper
	for k, v := range n.acronyms {
		res = strings.ReplaceAll(res, k, strings.ToUpper(v))
	}

	// 3. Lowercase for general similarity
	res = strings.ToLower(res)

	// 4. Unit normalization (e.g., "10 mg" -> "10mg")
	unitRegex := regexp.MustCompile(`(\d+)\s*(mg|g|mcg|ml|l|milligram|gram)`)
	res = unitRegex.ReplaceAllStringFunc(res, func(m string) string {
		m = strings.ReplaceAll(m, " ", "")
		m = strings.ReplaceAll(m, "milligram", "mg")
		m = strings.ReplaceAll(m, "gram", "g")
		return m
	})

	return strings.TrimSpace(res)
}
