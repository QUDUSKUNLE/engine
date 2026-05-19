package ai

import (
	"context"
	"github.com/diagnoxix/core/ports"
)

var _ ports.DecisionSupporter = (*DecisionSupport)(nil)

type DecisionSupport struct{}

// Recommend implements ports.DecisionSupporter.
func (d *DecisionSupport) Recommend(ctx context.Context, patientID string) ([]string, error) {
	panic("unimplemented")
}

func NewDecisionSupport() *DecisionSupport {
	return &DecisionSupport{}
}
