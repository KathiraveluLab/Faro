package types

// ImageLevel represents the individual image metadata.
type ImageLevel struct {
	SOPInstanceUID string `json:"sop_instance_uid"`
	ImageNumber    int    `json:"image_number"`
}

// SeriesLevel represents a collection of images.
type SeriesLevel struct {
	SeriesInstanceUID string       `json:"series_instance_uid"`
	SeriesDescription string       `json:"series_description"`
	Modality          string       `json:"modality"`
	Images            []ImageLevel `json:"images"`
}

// StudyLevel represents a clinical imaging study (top level for most duplicate detection).
type StudyLevel struct {
	StudyInstanceUID string        `json:"study_instance_uid"`
	StudyDate        string        `json:"study_date"`
	StudyDescription string        `json:"study_description"`
	Series           []SeriesLevel `json:"series"`
}

// MedicalImageMetadata is the consolidated metadata for a medical imaging record.
type MedicalImageMetadata struct {
	PatientID string     `json:"patient_id"`
	Studies   []StudyLevel `json:"studies"`
}
