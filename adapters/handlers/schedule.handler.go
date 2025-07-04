package handlers

import (
	"github.com/labstack/echo/v4"
)

// ScheduleSwagger is used for Swagger documentation only
// @Description Diagnostic schedule response for Swagger
// @name ScheduleSwagger
type ScheduleSwagger struct {
	ID                 string `json:"id" example:"sched-001"`
	UserID             string `json:"user_id" example:"user-001"`
	DiagnosticCentreID string `json:"diagnostic_centre_id" example:"dc-001"`
	TestType           string `json:"test_type" example:"Blood Test"`
	ScheduleTime       string `json:"schedule_time" example:"2025-06-26T09:00:00Z"` // format: date-time
	Status             string `json:"status" example:"pending"`
	CreatedAt          string `json:"created_at" example:"2025-06-25T20:00:00Z"` // format: date-time
	UpdatedAt          string `json:"updated_at" example:"2025-06-25T21:00:00Z"` // format: date-time
	// ...add other fields as needed for docs
}

// CreateSchedule godoc
// @Summary Create a new diagnostic schedule
// @Description Schedule a diagnostic test at a diagnostic centre. Requires user authentication.
// @Tags Schedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param schedule body domain.CreateScheduleDTO true "Schedule details"
// @Success 201 {object} handlers.ScheduleSwagger "Schedule created successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid schedule data/Invalid datetime format"
// @Failure 401 {object} handlers.ErrorResponse "Authentication required"
// @Failure 403 {object} handlers.ErrorResponse "Insufficient permissions"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
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
// @Param schedule_id path string true "Schedule ID (UUID format)" format(uuid)
// @Success 200 {object} handlers.ScheduleSwagger "Schedule details retrieved successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid schedule ID"
// @Failure 401 {object} handlers.ErrorResponse "Authentication required"
// @Failure 403 {object} handlers.ErrorResponse "Not the schedule owner"
// @Failure 404 {object} handlers.ErrorResponse "Schedule not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
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
// @Failure 401 {object} handlers.ErrorResponse "Authentication required"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
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
// @Param schedule_id path string true "Schedule ID (UUID format)" format(uuid)
// @Param schedule body domain.UpdateScheduleDTO true "Updated schedule details"
// @Success 200 {object} handlers.ScheduleSwagger "Schedule updated successfully"
// @Success 204 {object} nil "Schedule updated successfully with no content to return"
// @Failure 400 {object} handlers.ErrorResponse "Invalid schedule data/Invalid datetime format"
// @Failure 401 {object} handlers.ErrorResponse "Authentication required"
// @Failure 403 {object} handlers.ErrorResponse "Not the schedule owner"
// @Failure 404 {object} handlers.ErrorResponse "Schedule not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
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
// @Param schedule_id path string true "Schedule ID (UUID format)" format(uuid)
// @Success 204 {object} nil "Schedule deleted successfully"
// @Failure 401 {object} handlers.ErrorResponse "Authentication required"
// @Failure 403 {object} handlers.ErrorResponse "Not the schedule owner"
// @Failure 404 {object} handlers.ErrorResponse "Schedule not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
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
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid)
// @Param schedule_id path string true "Schedule ID" format(uuid)
// @Success 200 {object} handlers.ScheduleSwagger "Schedule details"
// @Failure 400 {object} handlers.ErrorResponse "Invalid ID format"
// @Failure 401 {object} handlers.ErrorResponse "Authentication required"
// @Failure 403 {object} handlers.ErrorResponse "Not authorized"
// @Failure 404 {object} handlers.ErrorResponse "Schedule not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
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
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid)
// @Param limit query integer false "Number of records to return" minimum(1) maximum(100) default(10)
// @Param offset query integer false "Number of records to skip" minimum(0) default(0)
// @Success 200 {array} handlers.ScheduleSwagger "List of schedules"
// @Failure 400 {object} handlers.ErrorResponse "Invalid diagnostic centre ID"
// @Failure 401 {object} handlers.ErrorResponse "Authentication required"
// @Failure 403 {object} handlers.ErrorResponse "Not authorized"
// @Failure 404 {object} handlers.ErrorResponse "Diagnostic centre not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
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
// @Param diagnostic_centre_id path string true "Diagnostic Centre ID" format(uuid)
// @Param schedule_id path string true "Schedule ID" format(uuid)
// @Param schedule body domain.UpdateDiagnosticScheduleByCentreDTO true "Updated schedule status"
// @Success 200 {object} handlers.ScheduleSwagger "Schedule updated successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid input data"
// @Failure 401 {object} handlers.ErrorResponse "Authentication required"
// @Failure 403 {object} handlers.ErrorResponse "Not authorized"
// @Failure 404 {object} handlers.ErrorResponse "Schedule not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /v1/diagnostic_centres/{diagnostic_centre_id}/diagnostic_schedules/{schedule_id} [put]
func (handler *HTTPHandler) UpdateDiagnosticScheduleByCentre(context echo.Context) error {
	return handler.service.UpdateDiagnosticScheduleByCentre(context)
}
