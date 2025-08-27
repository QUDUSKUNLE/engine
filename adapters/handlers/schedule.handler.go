package handlers

import (
	"github.com/labstack/echo/v4"
)

// CreateSchedule godoc
// @Summary Create a new diagnostic schedule
// @Description Schedule a diagnostic test at a diagnostic centre. Requires user authentication.
// @Tags Schedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param schedule body domain.CreateScheduleDTO true "Schedule details"
// @Success 201 {object} handlers.ScheduleSwagger "SUCCESS_RESPONSE"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_schedules [post]
func (handler *HTTPHandler) CreateSchedule(context echo.Context) error {
	return handler.service.CreateSchedule(context)
}

// GetSchedule godoc
// @Summary Get a specific diagnostic schedule
// @Description Retrieve details of a specific diagnostic schedule by ID. Only accessible by the schedule owner.
// @Tags Schedule
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param schedule_id path string true "Schedule ID (UUID format)" format(uuid) default(123e4567-e89b-12d3-a456-426614174000)
// @Success 200 {object} handlers.ScheduleSwagger "SUCCESS_RESPONSE"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_schedules/{schedule_id} [get]
func (handler *HTTPHandler) GetSchedule(context echo.Context) error {
	return handler.service.GetDiagnosticSchedule(context)
}

// GetSchedules godoc
// @Summary List user's diagnostic schedules
// @Description Get all diagnostic schedules for the authenticated user with pagination
// @Tags Schedule
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param limit query int false "Number of records to return" minimum(1) maximum(100) default(10)
// @Param offset query int false "Number of records to skip" minimum(0) default(0)
// @Success 200 {array} handlers.ScheduleSwagger "List of schedules"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_schedules [get]
func (handler *HTTPHandler) GetSchedules(context echo.Context) error {
	return handler.service.GetDiagnosticSchedules(context)
}

// UpdateSchedule godoc
// @Summary Update a diagnostic schedule
// @Description Update an existing diagnostic schedule. Only accessible by the schedule owner.
// @Tags Schedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param schedule_id path string true "Schedule ID (UUID format)" format(uuid) default(123e4567-e89b-12d3-a456-426614174000)
// @Param schedule body domain.UpdateScheduleDTO true "Updated schedule details"
// @Success 200 {object} handlers.ScheduleSwagger "Schedule updated successfully"
// @Success 204 {object} nil "Schedule updated successfully with no content to return"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_schedules/{schedule_id} [put]
func (handler *HTTPHandler) UpdateSchedule(context echo.Context) error {
	return handler.service.UpdateDiagnosticSchedule(context)
}

// DeleteDiagnosticSchedule godoc
// @Summary Delete a diagnostic schedule
// @Description Delete an existing diagnostic schedule. Only accessible by the schedule owner.
// @Tags Schedule
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param schedule_id path string true "Schedule ID (UUID format)" format(uuid) default(123e4567-e89b-12d3-a456-426614174000)
// @Success 204 {object} nil "Schedule deleted successfully"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_schedules/{schedule_id} [delete]
func (handler *HTTPHandler) DeleteDiagnosticSchedule(context echo.Context) error {
	return nil
}

// GetDiagnosticScheduleByCentre godoc
// @Summary Get a schedule by diagnostic centre
// @Description Retrieve a specific schedule for a diagnostic centre. Accessible by diagnostic centre staff.
// @Tags Schedule
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid) default(123e4567-e89b-12d3-a456-426614174000)
// @Param schedule_id path string true "Schedule ID" format(uuid) default(098e4567-e89b-12d3-a456-426614174000)
// @Success 200 {object} handlers.ScheduleSwagger "Schedule details"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/{diagnostic_centre_id}/diagnostic_schedules/{schedule_id} [get]
func (handler *HTTPHandler) GetDiagnosticScheduleByCentre(context echo.Context) error {
	return handler.service.GetDiagnosticScheduleByCentre(context)
}

// GetDiagnosticSchedulesByCentre godoc
// @Summary List schedules for a diagnostic centre
// @Description Get all schedules for a specific diagnostic centre with pagination
// @Tags Schedule
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid) default(123e4567-e89b-12d3-a456-426614174000)
// @Param limit query integer false "Number of records to return" minimum(1) maximum(100) default(10)
// @Param offset query integer false "Number of records to skip" minimum(0) default(0)
// @Success 200 {array} handlers.ScheduleSwagger "SUCCESS_RESPONSE"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/{diagnostic_centre_id}/diagnostic_schedules [get]
func (handler *HTTPHandler) GetDiagnosticSchedulesByCentre(context echo.Context) error {
	return handler.service.GetDiagnosticSchedulesByCentre(context)
}

// UpdateDiagnosticScheduleByCentre godoc
// @Summary Update schedule status by centre
// @Description Update a schedule's acceptance status. Only accessible by diagnostic centre staff.
// @Tags Schedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid) default(123e4567-e89b-12d3-a456-426614174000)
// @Param schedule_id path string true "Schedule ID" format(uuid) default(234e4567-e89b-12d3-a456-426614174000)
// @Param schedule body domain.UpdateDiagnosticScheduleByCentreDTO true "Updated schedule status"
// @Success 200 {object} handlers.ScheduleSwagger "SUCCESS_RESPONSE"
// @Failure 400 {object} handlers.BAD_REQUEST "BAD_REQUEST"
// @Failure 401 {object} handlers.UNAUTHORIZED_ERROR "UNAUTHORIZED_ERROR"
// @Failure 404 {object} handlers.NOT_FOUND_ERROR "NOT_FOUND_ERROR"
// @Failure 500 {object} handlers.INTERNAL_SERVER_ERROR "INTERNAL_SERVER_ERROR"
// @Router /v1/diagnostic_centres/{diagnostic_centre_id}/diagnostic_schedules/{schedule_id} [put]
func (handler *HTTPHandler) UpdateDiagnosticScheduleByCentre(context echo.Context) error {
	return handler.service.UpdateDiagnosticScheduleByCentre(context)
}
