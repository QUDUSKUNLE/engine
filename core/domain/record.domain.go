package domain

import (
	"github.com/google/uuid"
	"github.com/medicue/adapters/db"
)

// File represents an uploaded file's metadata and content.
// Update this struct as needed for your application's requirements.
type File struct {
	FileName string
	FileSize int64
	Content  []byte
}

type (
	CreateMedicalRecordDTO struct {
		UserID     uuid.UUID `json:"user_id" validate:"required,uuid"`
		UploaderID uuid.UUID `json:"uploader_id" validate:"required,uuid"`
		// Diagnostic Centre ID or Doctor ID
		UploaderAdminID uuid.UUID       `json:"uploader_admin_id"`
		UploaderType    db.UserEnum     `json:"uploader_type"`
		ScheduleID      uuid.UUID       `json:"schedule_id" validate:"required,uuid"`
		Title           string          `json:"title" validate:"required,min=12"`
		DocumentType    db.DocumentType `json:"document_type" validate:"oneof=LAB_REPORT PRESCRIPTION DISCHARGE_SUMMARY IMAGING VACCINATION ALLERGY SURGERY CHRONIC_CONDITION FAMILY_HISTORY"`
		FileUpload      File            // Define File type below or import from the correct package
		FilePath        string          `json:"file_path"`
		FileType        string          `json:"file_type"`
		UploadedAt      string          `json:"uploaded_at"`
		ProviderName    string          `json:"provider_name"`
		Specialty       string          `json:"specialty"`
		IsShared        bool            `json:"is_shared"`
		SharedUntil     string          `json:"shared_until" validate:"omitempty,datetime=2006-01-02"`
	}
)
