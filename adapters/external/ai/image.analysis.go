package ai

import (
	"context"
	"github.com/diagnoxix/core/ports"
)

var _ ports.ImageAnalyzer = (*ImageAnalysis)(nil)

type ImageAnalysis struct{}

// Analyze implements ports.ImageAnalyzer.
func (i *ImageAnalysis) Analyze(ctx context.Context, imagePath string) (map[string]interface{}, error) {
	panic("unimplemented")
}

func NewImageAnalysis() *ImageAnalysis {
	return &ImageAnalysis{}
}
