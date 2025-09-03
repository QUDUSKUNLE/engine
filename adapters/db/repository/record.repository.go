package repository

import (
	"context"

	"github.com/diagnoxix/adapters/db"
	"github.com/diagnoxix/core/ports"
)

// Ensure Repository implements RecordRepository
var _ ports.RecordRepository = (*Repository)(nil)

// CreateMedicalRecord implements the ports.RecordRepository interface.
func (r *Repository) CreateMedicalRecord(
	ctx context.Context,
	record db.CreateMedicalRecordParams,
) (*db.MedicalRecord, error) {
	return r.database.CreateMedicalRecord(ctx, record)
}

func (r *Repository) UpdateFilePath(
	ctx context.Context,
	req db.UpdateFilePathParams,
) (*db.MedicalRecord, error) {
	return r.database.UpdateFilePath(ctx, req)
}

// GetMedicalRecord implements the ports.RecordRepository interface.
func (r *Repository) GetMedicalRecord(
	ctx context.Context,
	params db.GetMedicalRecordParams,
) (*db.GetMedicalRecordRow, error) {
	return r.database.GetMedicalRecord(ctx, params)
}

// GetMedicalRecords implements the ports.RecordRepository interface.
func (r *Repository) GetMedicalRecords(
	ctx context.Context,
	params db.GetMedicalRecordsParams,
) ([]*db.GetMedicalRecordsRow, error) {
	return r.database.GetMedicalRecords(ctx, params)
}

// GetUploaderMedicalRecord implements the ports.RecordRepository interface.
func (r *Repository) GetUploaderMedicalRecord(
	ctx context.Context,
	params db.GetUploaderMedicalRecordParams,
) (*db.GetUploaderMedicalRecordRow, error) {
	return r.database.GetUploaderMedicalRecord(ctx, params)
}

// GetUploaderMedicalRecords implements the ports.RecordRepository interface.
func (r *Repository) GetUploaderMedicalRecords(
	ctx context.Context,
	params db.GetUploaderMedicalRecordsParams,
) ([]*db.GetUploaderMedicalRecordsRow, error) {
	return r.database.GetUploaderMedicalRecords(ctx, params)
}

// UpdateMedicalRecord implements the ports.RecordRepository interface.
func (r *Repository) UpdateMedicalRecord(
	ctx context.Context,
	params db.UpdateMedicalRecordByUploaderParams,
) (*db.MedicalRecord, error) {
	return r.database.UpdateMedicalRecordByUploader(ctx, params)
}
