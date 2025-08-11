package ports

import (
	"context"

	"github.com/diagnoxix/adapters/db"
)

type RecordRepository interface {
	// CreateMedicalRecord creates a new medical record with the provided parameters.
	CreateMedicalRecord(ctx context.Context, req db.CreateMedicalRecordParams) (*db.MedicalRecord, error)

	// GetMedicalRecord retrieves a single medical record based on the given parameters.
	GetMedicalRecord(ctx context.Context, req db.GetMedicalRecordParams) (*db.GetMedicalRecordRow, error)

	// GetMedicalRecords retrieves multiple medical records matching the specified criteria.
	GetMedicalRecords(ctx context.Context, req db.GetMedicalRecordsParams) ([]*db.GetMedicalRecordsRow, error)

	// GetUploaderMedicalRecord retrieves a single medical record uploaded by a specific uploader.
	GetUploaderMedicalRecord(ctx context.Context, req db.GetUploaderMedicalRecordParams) (*db.GetUploaderMedicalRecordRow, error)

	// GetUploaderMedicalRecords retrieves multiple medical records uploaded by a specific uploader.
	GetUploaderMedicalRecords(ctx context.Context, req db.GetUploaderMedicalRecordsParams) ([]*db.GetUploaderMedicalRecordsRow, error)

	// UpdateMedicalRecord updates an existing medical record by the uploader with the provided parameters.
	UpdateMedicalRecord(ctx context.Context, req db.UpdateMedicalRecordByUploaderParams) (*db.MedicalRecord, error)
}
