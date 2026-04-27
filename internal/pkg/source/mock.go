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
		{ID: "REC005", Source: "Postgres", Attributes: map[string]string{"name": "Metformin 500mg", "patient": "Bob Wilson"}},
		{ID: "REC006", Source: "Dicom", Attributes: map[string]string{"name": "Glucophage 500 mg", "patient": "Bob Wilson"}},
		{ID: "REC007", Source: "Excel", Attributes: map[string]string{"name": "Lisinopril 10mg", "patient": "Alice Brown"}},
		{ID: "REC008", Source: "Redis", Attributes: map[string]string{"name": "Prinivil 10 mg", "patient": "Alice Brown"}},
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
	if patientID == "PAT002" {
		return types.MedicalImageMetadata{
			PatientID: "PAT002",
			Studies: []types.StudyLevel{
				{
					StudyInstanceUID: "9.8.7.6.5",
					Series: []types.SeriesLevel{
						{SeriesInstanceUID: "9.8.7.6.5.1", Modality: "MR"},
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
