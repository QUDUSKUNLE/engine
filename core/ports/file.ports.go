package ports

import (
	"context"
)

type FileService interface {
	SaveFile(filePath string, content []byte) error
	UploadFile(context context.Context, filePath []byte) (string, error)
	DeleteFile(filePath string) error
	ExtractTextFromImage(filePath string) (map[string]interface{}, error)
	ExtractTextWithMonai(filePath, monaiURL string) (string, error)
}
