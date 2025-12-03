package handlers

// UserSwagger is used for Swagger documentation only
// @Description User response for Swagger
// @name UserSwagger
type UserSwagger struct {
	Status  int64 `json:"status" example:"201"`
	Success bool  `json:"success" example:"true"`
	Data    struct {
		ID        string `json:"id" example:"9c6bd6fe-493c-4b0a-a186-085039d24bf0"`
		Email     string `json:"email" example:"emi.user@diagnoxix.com"`
		Nin       string `json:"nin" example:"null"`
		CreatedAt string `json:"created_at" example:"2025-08-13T16:59:44.649447+01:00"`
		UpdatedAt string `json:"updated_at" example:"2025-08-13T16:59:44.649447+01:00"`
	} `json:"data"`
}

// DiagnosticCentreSwagger is used for Swagger documentation only
// @Description Diagnostic centre response for Swagger
// @name DiagnosticCentreSwagger
type DiagnosticCentreSwagger struct {
	ID        string  `json:"id" example:"dc-001"`
	Name      string  `json:"name" example:"DiagnoxixAI Diagnostics"`
	Latitude  float64 `json:"latitude" example:"6.5244"`
	Longitude float64 `json:"longitude" example:"3.3792"`
	Address   string  `json:"address" example:"123 Main St, Lagos"`
	Phone     string  `json:"phone" example:"+2348000000000"`
	Email     string  `json:"email" example:"info@diagnoxix.com"`
	CreatedAt string  `json:"created_at" example:"2025-06-26T20:00:00Z"`
	UpdatedAt string  `json:"updated_at" example:"2025-06-26T20:30:00Z"`
}

type ManagerSwagger struct {
	ID        string `json:"id" example:"user-001"`
	Email     string `json:"email" example:"user@example.com"`
	FullName  string `json:"full_name" example:"John Doe"`
	Role      string `json:"role" example:"DIAGNOSTIC_CENTRE_MANAGER"`
	CreatedAt string `json:"created_at" example:"2025-06-27T12:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2025-06-27T12:00:00Z"`
}

// ScheduleSwagger is used for Swagger documentation only
// @Description Diagnostic schedule response for Swagger
// @name ScheduleSwagger
type ScheduleSwagger struct {
	ID                 string `json:"id" example:"sched-001"`
	UserID             string `json:"user_id" example:"user-001"`
	DiagnosticCentreID string `json:"diagnostic_centre_id" example:"dc-001"`
	TestType           string `json:"test_type" example:"Blood Test"`
	ScheduleTime       string `json:"schedule_time" example:"2025-06-26T09:00:00Z"`
	Status             string `json:"status" example:"pending"`
	CreatedAt          string `json:"created_at" example:"2025-06-25T20:00:00Z"`
	UpdatedAt          string `json:"updated_at" example:"2025-06-25T21:00:00Z"`
}

// MedicalRecordSwagger is used for Swagger documentation only
// @Description Medical record response for Swagger
// @name MedicalRecordSwagger
type MedicalRecordSwagger struct {
	ID                 string `json:"id" example:"rec-001"`
	PatientID          string `json:"patient_id" example:"user-001"`
	DiagnosticCentreID string `json:"diagnostic_centre_id" example:"dc-001"`
	FileType           string `json:"file_type" example:"PDF"`
	FileURL            string `json:"file_url" example:"https://medivue.com/records/rec-001.pdf"`
	CreatedAt          string `json:"created_at" example:"2025-06-26T20:00:00Z"`
	UpdatedAt          string `json:"updated_at" example:"2025-06-26T20:30:00Z"`
}

// AppointmentSwagger is used for Swagger documentation only
// @Description Appointment response for Swagger
// @name AppointmentSwagger
// @property appointment_date string "Appointment date in RFC3339 format" example:"2025-06-26T21:00:00Z"
type AppointmentSwagger struct {
	ID                 string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	AppointmentDate    string `json:"appointment_date" example:"2025-06-26T21:00:00Z"` // format: date-time
	DiagnosticCentreID string `json:"diagnostic_centre_id" example:"dc-001"`
	PatientID          string `json:"patient_id" example:"user-001"`
	Status             string `json:"status" example:"pending"`
	CreatedAt          string `json:"created_at" example:"2025-06-26T20:00:00Z"`
	UpdatedAt          string `json:"updated_at" example:"2025-06-26T20:30:00Z"`
}

type SUCCESS_RESPONSE struct {
	Status  int64    `json:"status" example:"200"`
	Success bool     `json:"success" example:"true"`
	Data    struct{} `json:"data"`
}

// BAD_REQUEST is used for Swagger documentation only
// @Description BAD_REQUEST response for Swagger
// @name BAD_REQUEST
// @property code string "BAD_REQUEST" example:"BAD_REQUEST"
// @property message string "Error message" example:"Invalid request"
type BAD_REQUEST struct {
	Status  int64 `json:"status" example:"400"`
	Success bool  `json:"success" example:"false"`
	Error   struct {
		Code    string `json:"code" example:"BAD_REQUEST"`
		Message string `json:"message" example:"Invalid request"`
	} `json:"error"`
}

// UNAUTHORIZED_ERROR is used for Swagger documentation only
// @Description UNAUTHORIZED_ERROR response for Swagger
// @name UNAUTHORIZED_ERROR
// @property code string "UNAUTHORIZED_ERROR" example:"UNAUTHORIZED_ERROR"
// @property message string "Error message" example:"Unauthorized to perform this operation"
type UNAUTHORIZED_ERROR struct {
	Status  int64 `json:"status" example:"401"`
	Success bool  `json:"success" example:"false"`
	Error   struct {
		Code    string `json:"code" example:"UNAUTHORIZED_ERROR"`
		Message string `json:"message" example:"Unauthorized to perform this operation"`
	} `json:"error"`
}

// NOT_FOUND_ERROR is used for Swagger documentation only
// @Description NOT_FOUND_ERROR response for Swagger
// @name NOT_FOUND_ERROR
// @property code string "NOT_FOUND_ERROR" example:"NOT_FOUND_ERROR"
// @property message string "Error message" example:"Not found error"
type NOT_FOUND_ERROR struct {
	Status  int64 `json:"status" example:"404"`
	Success bool  `json:"success" example:"false"`
	Error   struct {
		Code    string `json:"code" example:"NOT_FOUND_ERROR"`
		Message string `json:"message" example:"Data not found"`
	} `json:"error"`
}

// DUPLICATE_ERROR is used for Swagger documentation only
// @Description DUPLICATE_ERROR response for Swagger
// @name DUPLICATE_ERROR
// @property code string "DUPLICATE_ERROR" example:"DUPLICATE_ERROR"
// @property message string "Error message" example:"Duplicate error"
type DUPLICATE_ERROR struct {
	Status  int64 `json:"status" example:"409"`
	Success bool  `json:"success" example:"false"`
	Error   struct {
		Code    string `json:"code" example:"DUPLICATE_ERROR"`
		Message string `json:"message" example:"Duplicate error"`
	} `json:"error"`
}

// UNPROCESSED_ERROR is used for Swagger documentation only
// @Description UNPROCESSED_ERROR response for Swagger
// @name DUPLICATE_ERROR
// @property code string "UNPROCESSED_ERROR" example:"UNPROCESSED_ERROR"
// @property message string "Error message" example:"Unprocessable entity"
type UNPROCESSED_ERROR struct {
	Status  int64 `json:"status" example:"422"`
	Success bool  `json:"success" example:"false"`
	Error   struct {
		Code    string `json:"code" example:"UNPROCESSED_ERROR"`
		Message string `json:"message" example:"Unprocessable entity"`
	} `json:"error"`
}

// INTERNAL_SERVER_ERROR is used for Swagger documentation only
// @Description INTERNAL_SERVER_ERROR response for Swagger
// @name INTERNAL_SERVER_ERROR
// @property code string "INTERNAL_SERVER_ERROR" example:"INTERNAL_SERVER_ERROR"
// @property message string "Error message" example:"Internal server error"
type INTERNAL_SERVER_ERROR struct {
	Status  int64 `json:"status" example:"500"`
	Success bool  `json:"success" example:"false"`
	Error   struct {
		Code    string `json:"code" example:"INTERNAL_SERVER_ERROR"`
		Message string `json:"message" example:"Internal Server Error"`
	} `json:"error"`
}

// UNAUTHORIZED_ERROR is used for Swagger documentation only
// @Description UNAUTHORIZED_ERROR response for Swagger
// @name UNAUTHORIZED_ERROR
// @property code string "UNAUTHORIZED_ERROR" example:"UNAUTHORIZED_ERROR"
// @property message string "Error message" example:"Unauthorized to perform this operation"
type ErrorResponse struct {
	Status  int64 `json:"status" example:"401"`
	Success bool  `json:"success" example:"false"`
	Error   struct {
		Code    string `json:"code" example:"BAD_REQUEST"`
		Message string `json:"message" example:"Invalid request"`
	} `json:"error"`
}
