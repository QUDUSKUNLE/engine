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
)

func cloudinaryUploader(context context.Context, file []byte) (string, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		fmt.Printf("Error handshaking cloudinary %s", err.Error())
		return "", err
	}
	uploadResult, err := cld.Upload.Upload(context, strings.NewReader(string(file)), uploader.UploadParams{Folder: "halalmeat"})
	if err != nil {
		fmt.Printf("Error uploading %s to cloudinary", err.Error())
		return "", err
	}
	fmt.Println("File uploaded successfully.")
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
