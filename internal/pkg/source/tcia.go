package source

import (
	"encoding/json"
	"faro/internal/pkg/types"
	"fmt"
	"net/http"
)

// TciaSource implements the Source interface for The Cancer Imaging Archive REST API.
type TciaSource struct {
	BaseURL string
	APIKey  string
}

func NewTciaSource(baseURL, apiKey string) *TciaSource {
	return &TciaSource{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

func (s *TciaSource) FetchRecords() ([]types.Record, error) {
	// TCIA is primarily for image metadata, but we can fetch patient collections
	// as textual records for discovery purposes.
	return nil, fmt.Errorf("textual record fetching not implemented for TCIA source")
}

func (s *TciaSource) FetchMetadata(patientID string) (types.MedicalImageMetadata, error) {
	// 1. Fetch Studies for the Patient
	studies, err := s.getStudies(patientID)
	if err != nil {
		return types.MedicalImageMetadata{}, err
	}

	// 2. Wrap in the hierarchical model
	result := types.MedicalImageMetadata{
		PatientID: patientID,
		Studies:   studies,
	}

	return result, nil
}

func (s *TciaSource) getStudies(patientID string) ([]types.StudyLevel, error) {
	url := fmt.Sprintf("%s/getPatientStudy?PatientID=%s&format=json", s.BaseURL, patientID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if s.APIKey != "" {
		req.Header.Add("api-key", s.APIKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TCIA API returned status %d", resp.StatusCode)
	}

	var tciaStudies []struct {
		StudyInstanceUID string `json:"StudyInstanceUID"`
		StudyDate        string `json:"StudyDate"`
		StudyDescription string `json:"StudyDescription"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tciaStudies); err != nil {
		return nil, err
	}

	var studies []types.StudyLevel
	for _, ts := range tciaStudies {
		studies = append(studies, types.StudyLevel{
			StudyInstanceUID: ts.StudyInstanceUID,
			StudyDate:        ts.StudyDate,
			StudyDescription: ts.StudyDescription,
		})
	}

	return studies, nil
}
