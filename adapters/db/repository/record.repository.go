package repository

import (
	"context"

	"github.com/medivue/adapters/db"
	"github.com/medivue/core/ports"
)

// Ensure Repository implements RecordRepository
var _ ports.RecordRepository = (*Repository)(nil)

// CreateMedicalRecord implements the ports.RecordRepository interface.
func (r *Repository) CreateMedicalRecord(ctx context.Context, record db.CreateMedicalRecordParams) (*db.MedicalRecord, error) {
	// TODO: implement the actual logic to create a medical record in the database
	return r.database.CreateMedicalRecord(ctx, record)
}

// GetMedicalRecord implements the ports.RecordRepository interface.
func (r *Repository) GetMedicalRecord(ctx context.Context, params db.GetMedicalRecordParams) (*db.GetMedicalRecordRow, error) {
	// TODO: implement the actual logic to get a medical record from the database
	return r.database.GetMedicalRecord(ctx, params)
}

// GetMedicalRecords implements the ports.RecordRepository interface.
func (r *Repository) GetMedicalRecords(ctx context.Context, params db.GetMedicalRecordsParams) ([]*db.GetMedicalRecordsRow, error) {
	// TODO: implement the actual logic to get multiple medical records from the database
	return r.database.GetMedicalRecords(ctx, params)
}

// GetUploaderMedicalRecord implements the ports.RecordRepository interface.
func (r *Repository) GetUploaderMedicalRecord(ctx context.Context, params db.GetUploaderMedicalRecordParams) (*db.GetUploaderMedicalRecordRow, error) {
	// TODO: implement the actual logic to get a medical record uploaded by a specific uploader from the database
	return r.database.GetUploaderMedicalRecord(ctx, params)
}

// GetUploaderMedicalRecords implements the ports.RecordRepository interface.
func (r *Repository) GetUploaderMedicalRecords(ctx context.Context, params db.GetUploaderMedicalRecordsParams) ([]*db.GetUploaderMedicalRecordsRow, error) {
	// TODO: implement the actual logic to get multiple medical records uploaded by a specific uploader from the database
	return r.database.GetUploaderMedicalRecords(ctx, params)
}

// UpdateMedicalRecord implements the ports.RecordRepository interface.
func (r *Repository) UpdateMedicalRecord(ctx context.Context, params db.UpdateMedicalRecordByUploaderParams) (*db.MedicalRecord, error) {
	// TODO: implement the actual logic to update a medical record in the database
	return r.database.UpdateMedicalRecordByUploader(ctx, params)
}
