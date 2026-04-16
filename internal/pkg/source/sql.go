package source

import (
	"database/sql"
	"faro/internal/pkg/types"
	"fmt"
)

// SqlSource implements the Source interface for clinical SQL databases.
type SqlSource struct {
	db        *sql.DB
	tableName string
}

func NewSqlSource(db *sql.DB, tableName string) *SqlSource {
	return &SqlSource{
		db:        db,
		tableName: tableName,
	}
}

func (s *SqlSource) FetchRecords() ([]types.Record, error) {
	query := fmt.Sprintf("SELECT id, source, patient_name, record_text FROM %s", s.tableName)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []types.Record
	for rows.Next() {
		var id, src, patient, text string
		if err := rows.Scan(&id, &src, &patient, &text); err != nil {
			return nil, err
		}
		records = append(records, types.Record{
			ID:     id,
			Source: src,
			Attributes: map[string]string{
				"name":    text,
				"patient": patient,
			},
		})
	}
	return records, nil
}

func (s *SqlSource) FetchMetadata(patientID string) (types.MedicalImageMetadata, error) {
	// In a real application, imaging metadata might also be cached in a sidecar SQL table.
	// For this prototype, we'll return an empty set or a placeholder error.
	return types.MedicalImageMetadata{}, fmt.Errorf("metadata fetching not implemented for SQL source")
}
