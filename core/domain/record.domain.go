package domain

import (
	"github.com/google/uuid"
	"github.com/medivue/adapters/db"
)

type (
	File struct {
		FileName string
		FileSize int64
		Content  []byte
	}
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
	GetMedicalRecordParamsDTO struct {
		RecordID uuid.UUID `json:"record_id" validate:"required,uuid"`
		UserID   uuid.UUID `json:"user_id"`
	}
	GetMedicalRecordsParamQueryDTO struct {
		Limit  int32 `query:"limit" validate:"omitempty,gte=0"`
		Offset int32 `query:"offset" validate:"omitempty,gte=0"`
	}
	GetUploaderMedicalRecordParamsDTO struct {
		RecordID   uuid.UUID `json:"record_id" validate:"required,uuid"`
		UploaderID uuid.UUID `json:"diagnostic_centre_id" validate:"required,uuid"`
	}
	GetUploaderMedicalRecordsParamQueryDTO struct {
		UploaderID uuid.UUID `json:"diagnostic_centre_id" validate:"required,uuid"`
		Limit      int32     `query:"limit" validate:"omitempty,gte=0"`
		Offset     int32     `query:"offset" validate:"omitempty,gte=0"`
	}
	UpdateMedicalRecordDTO struct {
		RecordID        uuid.UUID       `json:"record_id" validate:"required,uuid"`
		UploaderID      uuid.UUID       `json:"diagnostic_centre_id" validate:"required,uuid"`
		UploaderAdminID uuid.UUID       `json:"uploader_admin_id"`
		Title           string          `json:"title" validate:"required,min=12"`
		DocumentType    db.DocumentType `json:"document_type" validate:"oneof=LAB_REPORT PRESCRIPTION DISCHARGE_SUMMARY IMAGING VACCINATION ALLERGY SURGERY CHRONIC_CONDITION FAMILY_HISTORY"`
		FileUpload      File            // Define File type below or import from the correct package
		FilePath        string          `json:"file_path"`
		FileType        string          `json:"file_type"`
		DocumentDate    string          `json:"document_date"`
		UploadedAt      string          `json:"uploaded_at"`
		ProviderName    string          `json:"provider_name"`
		Specialty       string          `json:"specialty"`
		IsShared        bool            `json:"is_shared"`
		SharedUntil     string          `json:"shared_until" validate:"omitempty,datetime=2006-01-02"`
	}
)

// GetLimit returns the limit value
func (q *GetMedicalRecordsParamQueryDTO) GetLimit() int32 {
	return q.Limit
}

// GetOffset returns the offset value
func (q *GetMedicalRecordsParamQueryDTO) GetOffset() int32 {
	return q.Offset
}

// SetLimit sets the limit value
func (q *GetMedicalRecordsParamQueryDTO) SetLimit(limit int32) {
	q.Limit = limit
}

// SetOffset sets the offset value
func (q *GetMedicalRecordsParamQueryDTO) SetOffset(offset int32) {
	q.Offset = offset
}

// GetLimit returns the limit value for uploader records
func (q *GetUploaderMedicalRecordsParamQueryDTO) GetLimit() int32 {
	return q.Limit
}

// GetOffset returns the offset value for uploader records
func (q *GetUploaderMedicalRecordsParamQueryDTO) GetOffset() int32 {
	return q.Offset
}

// SetLimit sets the limit value for uploader records
func (q *GetUploaderMedicalRecordsParamQueryDTO) SetLimit(limit int32) {
	q.Limit = limit
}

// SetOffset sets the offset value for uploader records
func (q *GetUploaderMedicalRecordsParamQueryDTO) SetOffset(offset int32) {
	q.Offset = offset
}
