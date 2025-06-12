package ex

import (
	"context"
	"os"

	"github.com/medicue/core/ports"
)

var _ ports.FileService = (*LocalFileService)(nil)

type LocalFileService struct{}

// UploadFile implements the ports.FileService interface.
func (l *LocalFileService) UploadFile(ctx context.Context, filePath []byte) (string, error) {
	return cloudinaryUploader(ctx, filePath)
}

// SaveFile implements the ports.FileService interface.
func (l *LocalFileService) SaveFile(filePath string, data []byte) error {
	// TODO: implement file saving logic
	return os.WriteFile(filePath, data, 0644)
}

// ExtractTextWithMonai implements the ports.FileService interface.
func (l *LocalFileService) ExtractTextWithMonai(imagePath string, model string) (string, error) {
	// TODO: implement MONAI-based text extraction logic
	return "", nil
}

// ExtractTextFromImage implements the ports.FileService interface.
func (l *LocalFileService) ExtractTextFromImage(imagePath string) (map[string]interface{}, error) {
	// TODO: implement text extraction logic
	return map[string]interface{}{}, nil
}

// DeleteFile implements the ports.FileService interface.
func (l *LocalFileService) DeleteFile(filePath string) error {
	// TODO: implement file deletion logic
	return os.Remove(filePath)
}
