package response

import "net/http"

// Standard error codes
const (
	// Generic error codes
	CodeNotFound        = "NOT_FOUND"
	CodeForbidden       = "FORBIDDEN"
	CodeBadRequest      = "BAD_REQUEST"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeInternalError   = "INTERNAL_ERROR"
	CodeDuplicateError  = "DUPLICATE_ERROR"
	CodeValidationError = "VALIDATION_ERROR"

	// Authentication error codes
	CodeInvalidToken       = "INVALID_TOKEN"
	CodeTokenExpired       = "TOKEN_EXPIRED"
	CodeInvalidCredentials = "INVALID_CREDENTIALS"

	// Resource error codes
	CodeResourceNotFound = "RESOURCE_NOT_FOUND"
	CodeResourceConflict = "RESOURCE_CONFLICT"

	// Database error codes
	CodeDatabaseError   = "DATABASE_ERROR"
	CodeUniqueViolation = "UNIQUE_VIOLATION"
)

// Error messages
const (
	MsgNotFound           = "Resource not found"
	MsgForbidden          = "Permission denied"
	MsgBadRequest         = "Invalid request"
	MsgTokenExpired       = "Token has expired"
	MsgInvalidToken       = "Invalid or malformed token"
	MsgUnauthorized       = "Authentication required"
	MsgInternalError      = "Internal server error"
	MsgDatabaseError      = "Database operation failed"
	MsgDuplicateError     = "Resource already exists"
	MsgValidationError    = "Validation failed"
	MsgInvalidCredentials = "Invalid email or password"
)

// HTTP status code to error code mapping
var StatusToCode = map[int]string{
	http.StatusNotFound:            CodeNotFound,
	http.StatusConflict:            CodeDuplicateError,
	http.StatusForbidden:           CodeForbidden,
	http.StatusBadRequest:          CodeBadRequest,
	http.StatusUnauthorized:        CodeUnauthorized,
	http.StatusInternalServerError: CodeInternalError,
}

// Error code to default message mapping
var CodeToMessage = map[string]string{
	CodeNotFound:           MsgNotFound,
	CodeForbidden:          MsgForbidden,
	CodeBadRequest:         MsgBadRequest,
	CodeTokenExpired:       MsgTokenExpired,
	CodeInvalidToken:       MsgInvalidToken,
	CodeUnauthorized:       MsgUnauthorized,
	CodeInternalError:      MsgInternalError,
	CodeDatabaseError:      MsgDatabaseError,
	CodeDuplicateError:     MsgDuplicateError,
	CodeValidationError:    MsgValidationError,
	CodeInvalidCredentials: MsgInvalidCredentials,
}
