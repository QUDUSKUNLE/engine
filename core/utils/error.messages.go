package utils

import "errors"

const (
	ErrBadRequest          = "bad Request"
	InvalidRequest         = "invalid request"
	PermissionDenied       = "Permission denied"
	InvalidLoginRequest    = "invalid login credentials"
	InvalidRequestBody     = "Invalid request body"
	AuthenticationRequired = "Authentication required"
	MissingParameters      = "Missing required parameters"

	// Schedule
	ScheduleNotFound       = "Schedule not found"
	FailedToUpdateSchedule = "Failed to update schedule"
)

// Common error types
var (
	ErrNotFound         = errors.New("resource not found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidInput     = errors.New("invalid input")
	ErrDatabaseError    = errors.New("database error")
	ErrConflictError    = errors.New("ERROR: duplicate key value violates unique constraint \"diagnostic_centres_latitude_longitude_key\" (SQLSTATE 23505)")
)
