package ex

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/diagnoxix/core/ports"
	"github.com/diagnoxix/core/utils"
)

// Supported file types and their extensions
var (
	supportedImageTypes = map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true,
		".gif": true, ".bmp": true, ".webp": true,
	}
	supportedDocTypes = map[string]bool{
		".pdf": true, ".doc": true, ".docx": true,
	}
	maxFileSize = int64(32 << 20) // 32MB
)

var _ ports.FileService = (*LocalFileService)(nil)

type LocalFileService struct {
	uploadDir string // Base directory for file uploads
}

// NewLocalFileService creates a new LocalFileService instance
func NewLocalFileService() *LocalFileService {
	return &LocalFileService{}
}

// UploadFile implements the ports.FileService interface.
func (l *LocalFileService) UploadFile(ctx context.Context, fileData []byte) (string, error) {
	if len(fileData) == 0 {
		utils.Error("Empty file data provided")
		return "", errors.New("empty file data")
	}

	if int64(len(fileData)) > maxFileSize {
		utils.Error("File size exceeds limit",
			utils.LogField{Key: "size", Value: len(fileData)},
			utils.LogField{Key: "max_size", Value: maxFileSize})
		return "", errors.New("file size exceeds maximum limit")
	}

	utils.Info("Uploading file to Cloudinary",
		utils.LogField{Key: "size", Value: len(fileData)})

	return cloudinaryUploader(ctx, fileData)
}

// SaveFile implements the ports.FileService interface.
func (l *LocalFileService) SaveFile(filePath string, data []byte) error {
	if !l.isValidFile(filePath, data) {
		return errors.New("invalid file type or size")
	}

	fullPath := filepath.Join(l.uploadDir, filePath)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		utils.Error("Failed to create directory",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "path", Value: dir})
		return fmt.Errorf("failed to create directory: %w", err)
	}

	utils.Info("Saving file",
		utils.LogField{Key: "path", Value: fullPath},
		utils.LogField{Key: "size", Value: len(data)})

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		utils.Error("Failed to save file",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "path", Value: fullPath})
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// ExtractTextWithMonai implements the ports.FileService interface.
func (l *LocalFileService) ExtractTextWithMonai(imagePath string, model string) (string, error) {
	if !l.fileExists(imagePath) {
		return "", errors.New("image file not found")
	}

	if !l.isImageFile(imagePath) {
		return "", errors.New("invalid image file type")
	}

	utils.Info("Extracting text with MONAI",
		utils.LogField{Key: "path", Value: imagePath},
		utils.LogField{Key: "model", Value: model})

	// TODO: Integration with MONAI service
	// This would typically involve:
	// 1. Loading the image
	// 2. Preprocessing the image
	// 3. Sending to MONAI service
	// 4. Processing the response

	return "", errors.New("MONAI integration not implemented")
}

// ExtractTextFromImage implements the ports.FileService interface.
func (l *LocalFileService) ExtractTextFromImage(imagePath string) (map[string]interface{}, error) {
	if !l.fileExists(imagePath) {
		return nil, errors.New("image file not found")
	}

	if !l.isImageFile(imagePath) {
		return nil, errors.New("invalid image file type")
	}

	utils.Info("Extracting text from image",
		utils.LogField{Key: "path", Value: imagePath})

	// TODO: Implement OCR integration
	// This would typically involve:
	// 1. Loading the image
	// 2. Using OCR library (e.g., Tesseract)
	// 3. Processing and structuring the results

	return map[string]interface{}{
		"status": "not_implemented",
		"path":   imagePath,
		"time":   time.Now(),
	}, errors.New("OCR integration not implemented")
}

// DeleteFile implements the ports.FileService interface.
func (l *LocalFileService) DeleteFile(filePath string) error {
	fullPath := filepath.Join(l.uploadDir, filePath)

	if !l.fileExists(fullPath) {
		utils.Error("File not found",
			utils.LogField{Key: "path", Value: fullPath})
		return errors.New("file not found")
	}

	utils.Info("Deleting file",
		utils.LogField{Key: "path", Value: fullPath})

	if err := os.Remove(fullPath); err != nil {
		utils.Error("Failed to delete file",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "path", Value: fullPath})
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// Helper functions

func (l *LocalFileService) isValidFile(filePath string, data []byte) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	// Check file type
	if !supportedImageTypes[ext] && !supportedDocTypes[ext] {
		utils.Error("Unsupported file type",
			utils.LogField{Key: "extension", Value: ext},
			utils.LogField{Key: "path", Value: filePath})
		return false
	}

	// Check file size
	if int64(len(data)) > maxFileSize {
		utils.Error("File size exceeds limit",
			utils.LogField{Key: "size", Value: len(data)},
			utils.LogField{Key: "max_size", Value: maxFileSize},
			utils.LogField{Key: "path", Value: filePath})
		return false
	}

	return true
}

func (l *LocalFileService) isImageFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return supportedImageTypes[ext]
}

func (l *LocalFileService) fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
