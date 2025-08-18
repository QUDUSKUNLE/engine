package ai

import (
	"context"
	"github.com/diagnoxix/core/ports"
)

var _ ports.ReportGenerator = (*AutomatedReport)(nil)

type AutomatedReport struct{}

// Generate implements ports.ReportGenerator.
func (a *AutomatedReport) Generate(ctx context.Context, input map[string]interface{}) (string, error) {
	panic("unimplemented")
}

func NewAutomatedReport() *AutomatedReport {
	return &AutomatedReport{}
}
