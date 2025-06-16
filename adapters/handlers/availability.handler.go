package handlers

import (
	"github.com/labstack/echo/v4"
)

// @Summary Create availability for a diagnostic centre
// @Description Create a new availability slot for the diagnostic centre
// @Tags Availability
// @Accept json
// @Produce json
// @Param availability body domain.CreateAvailabilityDTO true "Availability information"
// @Success 201 {object} domain.AvailabilitySlot
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/v1/availability [post]
func (handler *HTTPHandler) CreateAvailability(context echo.Context) error {
	return nil
}

// @Summary Get availability for a diagnostic centre
// @Description Get availability slots for the diagnostic centre, optionally filtered by day of week
// @Tags Availability
// @Accept json
// @Produce json
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID"
// @Param day_of_week query string false "Day of week (monday, tuesday, etc.)"
// @Success 200 {array} domain.AvailabilitySlot
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/v1/availability/{diagnostic_centre_id} [get]
func (handler *HTTPHandler) GetAvailability(context echo.Context) error {
	return nil
}

// @Summary Update availability for a diagnostic centre
// @Description Update an existing availability slot for the diagnostic centre
// @Tags Availability
// @Accept json
// @Produce json
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID"
// @Param day_of_week path string true "Day of week (monday, tuesday, etc.)"
// @Param availability body domain.UpdateAvailabilityDTO true "Updated availability information"
// @Success 200 {object} domain.AvailabilitySlot
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/v1/availability/{diagnostic_centre_id}/{day_of_week} [put]
func (h *HTTPHandler) UpdateAvailability(c echo.Context) error {
	return nil
}

// @Summary Bulk update availability for a diagnostic centre
// @Description Update multiple availability slots for the diagnostic centre
// @Tags Availability
// @Accept json
// @Produce json
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID"
// @Param availability body domain.UpdateManyAvailabilityDTO true "Updated availability information"
// @Success 200 {array} domain.AvailabilitySlot
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/v1/availability/{diagnostic_centre_id} [put]
func (handler *HTTPHandler) UpdateManyAvailability(context echo.Context) error {
	return handler.service.UpdateManyAvailability(context)
}

// @Summary Delete availability for a diagnostic centre
// @Description Delete an availability slot for the diagnostic centre
// @Tags Availability
// @Accept json
// @Produce json
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID"
// @Param day_of_week path string true "Day of week (monday, tuesday, etc.)"
// @Success 204 "No Content"
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/v1/availability/{diagnostic_centre_id}/{day_of_week} [delete]
func (h *HTTPHandler) DeleteAvailability(c echo.Context) error {

	return nil
}
