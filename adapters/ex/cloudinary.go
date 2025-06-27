package ex

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/medivue/core/utils"
)

func cloudinaryUploader(ctx context.Context, file []byte) (string, error) {
	utils.Info("Initializing Cloudinary connection",
		utils.LogField{Key: "cloud_name", Value: os.Getenv("CLOUDINARY_CLOUD_NAME")})

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		utils.Error("Failed to initialize Cloudinary",
			utils.LogField{Key: "error", Value: err.Error()})
		return "", err
	}

	utils.Info("Starting file upload to Cloudinary",
		utils.LogField{Key: "file_size", Value: len(file)})

	uploadResult, err := cld.Upload.Upload(ctx, strings.NewReader(string(file)), uploader.UploadParams{
		Folder:       "medivue",
		ResourceType: "auto",
	})
	if err != nil {
		utils.Error("Failed to upload file to Cloudinary",
			utils.LogField{Key: "error", Value: err.Error()},
			utils.LogField{Key: "file_size", Value: len(file)})
		return "", err
	}

	utils.Info("File uploaded successfully to Cloudinary",
		utils.LogField{Key: "public_id", Value: uploadResult.PublicID},
		utils.LogField{Key: "resource_type", Value: uploadResult.ResourceType},
		utils.LogField{Key: "size", Value: uploadResult.Bytes})

	return uploadResult.SecureURL, nil
}

func ExtractValue(text, pattern string) string {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func CalculateEDD(scanDate time.Time, ga string) time.Time {
	re := regexp.MustCompile(`([0-9]+)[wW]\s*([0-9]+)[dD]?`)
	match := re.FindStringSubmatch(ga)
	if len(match) < 3 {
		return scanDate
	}
	weeks := toInt(match[1])
	days := toInt(match[2])
	gaDays := weeks*7 + days
	return scanDate.AddDate(0, 0, 280-gaDays)
}

func toInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
