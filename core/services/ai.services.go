package services

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/medivue/core/domain"
	"github.com/medivue/core/utils"
)

func (service *ServicesHandler) InterpretLabHandler(ctx echo.Context) error {
	dto, _ := ctx.Get(utils.ValidatedBodyDTO).(*domain.LabTest)

	// Format a prompt for AI
	prompt := "Interpret the following lab test for a patient:\n"
	prompt += fmt.Sprintf("Test: %s\n", dto.TestName)
	for k, v := range dto.Results {
		ref := dto.ReferenceRanges[k]
		prompt += fmt.Sprintf("%s: %s (Normal: %s)\n", k, v, ref)
	}
	prompt += "Provide a summary for both patient and doctor."

	interpretation, err := service.aiService.InterpretLabResult(prompt)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, ctx)
	}

	return utils.ResponseMessage(http.StatusOK, interpretation, ctx)
}
