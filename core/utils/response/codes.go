package response

import "net/http"

// Standard error codes
const (
	// Generic error codes
	CodeBadRequest      = "BAD_REQUEST"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden       = "FORBIDDEN"
	CodeNotFound        = "NOT_FOUND"
	CodeInternalError   = "INTERNAL_ERROR"
	CodeValidationError = "VALIDATION_ERROR"
	CodeDuplicateError  = "DUPLICATE_ERROR"

	// Authentication error codes
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeTokenExpired       = "TOKEN_EXPIRED"
	CodeInvalidToken       = "INVALID_TOKEN"

	// Resource error codes
	CodeResourceNotFound = "RESOURCE_NOT_FOUND"
	CodeResourceConflict = "RESOURCE_CONFLICT"

	// Database error codes
	CodeDatabaseError   = "DATABASE_ERROR"
	CodeUniqueViolation = "UNIQUE_VIOLATION"
)

// Error messages
const (
	MsgBadRequest         = "Invalid request"
	MsgUnauthorized       = "Authentication required"
	MsgForbidden          = "Permission denied"
	MsgNotFound           = "Resource not found"
	MsgInternalError      = "Internal server error"
	MsgValidationError    = "Validation failed"
	MsgInvalidCredentials = "Invalid email or password"
	MsgTokenExpired       = "Token has expired"
	MsgInvalidToken       = "Invalid or malformed token"
	MsgDuplicateError     = "Resource already exists"
	MsgDatabaseError      = "Database operation failed"
)

// HTTP status code to error code mapping
var StatusToCode = map[int]string{
	http.StatusBadRequest:          CodeBadRequest,
	http.StatusUnauthorized:        CodeUnauthorized,
	http.StatusForbidden:           CodeForbidden,
	http.StatusNotFound:            CodeNotFound,
	http.StatusInternalServerError: CodeInternalError,
	http.StatusConflict:            CodeDuplicateError,
}

// Error code to default message mapping
var CodeToMessage = map[string]string{
	CodeBadRequest:         MsgBadRequest,
	CodeUnauthorized:       MsgUnauthorized,
	CodeForbidden:          MsgForbidden,
	CodeNotFound:           MsgNotFound,
	CodeInternalError:      MsgInternalError,
	CodeValidationError:    MsgValidationError,
	CodeInvalidCredentials: MsgInvalidCredentials,
	CodeTokenExpired:       MsgTokenExpired,
	CodeInvalidToken:       MsgInvalidToken,
	CodeDuplicateError:     MsgDuplicateError,
	CodeDatabaseError:      MsgDatabaseError,
}
