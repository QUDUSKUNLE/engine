package ai

import (
	"context"
	"github.com/diagnoxix/core/ports"
)

var _ ports.AnomalyDetector = (*AnomalyDetection)(nil)

type AnomalyDetection struct{}

func (anomaly *AnomalyDetection) Detect(ctx context.Context, data []float64) ([]string, error) {
	return []string{}, nil
}

func NewAnomalyDetection() *AnomalyDetection {
	return &AnomalyDetection{}
}
