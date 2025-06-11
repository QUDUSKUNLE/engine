package utils

const (
	Limit                  = 100
	Offset                 = 0
	ValidatedBodyDTO       = "validatedBodyDTO"
	ValidatedQueryParamDTO = "validatedQueryParamDTO"
	ScheduleID             = "schedule_id"
	DiagnosticCentreID     = "diagnostic_centre_id"
)

const (
	ErrScheduleTimeFormatRFC3339 = "invalid schedule time format, expected RFC3339 or RFC3339Nano"
	ErrScheduleTimeFormatParsing = "invalid schedule time format after parsing"
)
