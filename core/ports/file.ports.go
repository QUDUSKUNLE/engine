package ports

import (
	"context"
)

type FileService interface {
	SaveFile(filePath string, content []byte) error
	UploadFile(context context.Context, filePath []byte) (string, error)
	DeleteFile(filePath string) error
	// ExtractTextFromScan(filePath string) thirdparty.Report
	ExtractTextFromImage(filePath string) (map[string]interface{}, error)
	ExtractTextWithMonai(filePath string, monaiURL string) (string, error)
	// MedicalAIAssistant(payload thirdparty.FileUpload) string
}
