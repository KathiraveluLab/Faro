package source

import (
	"faro/internal/pkg/types"
)

// MockSource is a hardcoded data provider for testing purposes.
type MockSource struct{}

func (s *MockSource) FetchRecords() ([]types.Record, error) {
	return []types.Record{
		{ID: "REC001", Source: "MySQL", Attributes: map[string]string{"name": "ASA 81mg", "patient": "John Doe"}},
		{ID: "REC002", Source: "MongoDB", Attributes: map[string]string{"name": "Aspirin 81 mg", "patient": "John Doe"}},
		{ID: "REC003", Source: "CSV", Attributes: map[string]string{"name": "HCTZ 25 milligram", "patient": "Jane Smith"}},
		{ID: "REC004", Source: "SQL", Attributes: map[string]string{"name": "Hydrochlorothiazide 25mg", "patient": "Jane Smith"}},
	}, nil
}

func (s *MockSource) FetchMetadata(patientID string) (types.MedicalImageMetadata, error) {
	if patientID == "PAT001" {
		return types.MedicalImageMetadata{
			PatientID: "PAT001",
			Studies: []types.StudyLevel{
				{
					StudyInstanceUID: "1.2.3.4.5",
					Series: []types.SeriesLevel{
						{SeriesInstanceUID: "1.2.3.4.5.1", Modality: "CT"},
						{SeriesInstanceUID: "1.2.3.4.5.2", Modality: "CT"},
					},
				},
			},
		}, nil
	}
	// Default conflict for demo
	return types.MedicalImageMetadata{
		PatientID: "PAT001-CONFLICT",
		Studies: []types.StudyLevel{
			{
				StudyInstanceUID: "1.2.3.4.5",
				Series: []types.SeriesLevel{
					{SeriesInstanceUID: "1.2.3.4.5.1", Modality: "CT"},
					{SeriesInstanceUID: "1.2.3.4.5.3", Modality: "CT"},
				},
			},
		},
	}, nil
}
