# 🏥 Diagnoxix API Documentati
## 📋 Table of Contents

- [Overview](#overview)
- [Authentication](#authentication)
- [Base URL & Headers](#base-url--headers)
- [Response Format](#response-format)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)
- [API Endpoints](#api-endpoints)
  - [Authentication](#authentication-endpoints)
  - [User Management](#user-management)
  - [Diagnostic Centres](#diagnostic-centres)
  - [Appointments](#appointments)
  - [Medical Records](#medical-records)
  - [Payments](#payments)
  - [Schedules](#schedules)
  - [Availability](#availability)
  - [AI Services](#ai-services)
  - [WebSocket](#websocket)
  - [Notifications](#notifications)

---

## 📖 Overview

The Diagnoxix API is a comprehensive healthcare platform that provides endpoints for managing diagnostic centres, appointments, medical records, payments, and AI-powered medical analysis. The API follows RESTful principles and uses JSON for data exchange.

### Key Features
- 🔐 JWT-based authentication
- 🏥 Multi-role user system (Patients, Managers, Owners)
- 🤖 AI-powered medical analysis
- ⚡ Real-time WebSocket notifications
- 💳 Integrated payment processing
- 📱 Mobile-friendly endpoints

---

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header for protected endpoints.

### Authentication Flow
1. Register or login to get a JWT token
2. Include the token in subsequent requests
3. Token expires after configured time (default: 1080 hours)

### Headers
```http
Authorization: Bearer <your-jwt-token>
Content-Type: application/json
```

---

## 🌐 Base URL & Headers

### Base URL
```
Production: https://diagnoxix.onrender.com
Development: http://localhost:7556
```

### Required Headers
```http
Content-Type: application/json
Authorization: Bearer <jwt-token>  # For protected endpoints
```

---

## 📊 Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "status": 200,
  "success": true,
  "data": {
    // Response data here
  }
}
```

### Error Response
```json
{
  "status": 400,
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Invalid request parameters"
  }
}
```

---

## ⚠️ Error Handling

### HTTP Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict (Duplicate)
- `422` - Unprocessable Entity
- `500` - Internal Server Error

### Error Codes
- `BAD_REQUEST` - Invalid request parameters
- `UNAUTHORIZED_ERROR` - Authentication required
- `NOT_FOUND_ERROR` - Resource not found
- `DUPLICATE_ERROR` - Resource already exists
- `UNPROCESSED_ERROR` - Validation failed
- `INTERNAL_SERVER_ERROR` - Server error

---

## 🚦 Rate Limiting

- **Default**: 100 requests per minute per IP
- **WebSocket**: 60 messages per minute per connection
- **AI Endpoints**: 20 requests per minute per user

---

# 🔗 API Endpoints

## 🔐 Authentication Endpoints

### Register User
Create a new user account.

```http
POST /v1/register
```

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "password": "SecurePassword123",
  "confirm_password": "SecurePassword123",
  "user_type": "PATIENT"
}
```

**Response:**
```json
{
  "status": 201,
  "success": true,
  "data": {
    "id": "9c6bd6fe-493c-4b0a-a186-085039d24bf0",
    "email": "john.doe@example.com",
    "user_type": "PATIENT",
    "created_at": "2025-01-15T16:59:44.649447Z"
  }
}
```

### Login
Authenticate user and get JWT token.

```http
POST /v1/login
```

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "SecurePassword123"
}
```

**Response:**
```json
{
  "status": 200,
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "9c6bd6fe-493c-4b0a-a186-085039d24bf0",
      "email": "john.doe@example.com",
      "user_type": "PATIENT"
    }
  }
}
```

### Request Password Reset
Request a password reset token.

```http
POST /v1/request_password_reset
```

**Request Body:**
```json
{
  "email": "john.doe@example.com"
}
```

### Reset Password
Reset password using token.

```http
POST /v1/reset_password
```

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "token": "reset-token-here",
  "new_password": "NewSecurePassword123",
  "confirm_password": "NewSecurePassword123"
}
```

### Update Password
Update password for authenticated user.

```http
PUT /v1/update_password
```

**Request Body:**
```json
{
  "new_password": "NewSecurePassword123",
  "confirm_password": "NewSecurePassword123"
}
```

### Verify Email
Verify user email address.

```http
GET /v1/verify_email?email=user@example.com&token=verification-token
```

### Resend Verification
Resend email verification.

```http
POST /v1/resend_verification
```

**Request Body:**
```json
{
  "email": "john.doe@example.com"
}
```

---

## 👤 User Management

### Get User Profile
Get current user's profile.

```http
GET /v1/account/profile
```

**Response:**
```json
{
  "status": 200,
  "success": true,
  "data": {
    "id": "9c6bd6fe-493c-4b0a-a186-085039d24bf0",
    "email": "john.doe@example.com",
    "full_name": "John Doe",
    "phone_number": "+2348000000000",
    "user_type": "PATIENT",
    "created_at": "2025-01-15T16:59:44.649447Z"
  }
}
```

### Update User Profile
Update current user's profile.

```http
PUT /v1/account/profile
```

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "+2348000000000",
  "nin": "12345678901"
}
```

### List Managers
List diagnostic centre managers (Admin only).

```http
GET /v1/managers?assigned=true&page=1&per_page=20
```

**Query Parameters:**
- `assigned` (boolean, optional) - Filter by assignment status
- `page` (integer, optional) - Page number (default: 1)
- `per_page` (integer, optional) - Items per page (default: 20, max: 50)

---

## 🏥 Diagnostic Centres

### Create Diagnostic Centre
Create a new diagnostic centre.

```http
POST /v1/diagnostic_centres
```

**Request Body:**
```json
{
  "diagnostic_centre_name": "Diagnoxix Medical Centre",
  "latitude": 6.5244,
  "longitude": 3.3792,
  "address": {
    "street": "123 Main Street",
    "city": "Lagos",
    "state": "Lagos",
    "country": "Nigeria"
  },
  "contact": {
    "phone": ["+2348000000000", "+2348111111111"],
    "email": "info@diagnoxix.com"
  },
  "doctors": ["Male", "Female"],
  "available_tests": [
    {
      "test_type": "BLOOD_TEST",
      "price": 5000.00
    },
    {
      "test_type": "X_RAY",
      "price": 15000.00
    }
  ]
}
```

### Get Diagnostic Centre
Get diagnostic centre details.

```http
GET /v1/diagnostic_centres/{diagnostic_centre_id}
```

**Response:**
```json
{
  "status": 200,
  "success": true,
  "data": {
    "id": "dc-001",
    "name": "Diagnoxix Medical Centre",
    "latitude": 6.5244,
    "longitude": 3.3792,
    "address": "123 Main Street, Lagos, Lagos, Nigeria",
    "phone": ["+2348000000000"],
    "email": "info@diagnoxix.com",
    "doctors": ["Male", "Female"],
    "available_tests": ["BLOOD_TEST", "X_RAY"],
    "created_at": "2025-01-15T20:00:00Z"
  }
}
```

### Search Diagnostic Centres
Search for diagnostic centres with filters.

```http
GET /v1/diagnostic_centres?latitude=6.5244&longitude=3.3792&test=BLOOD_TEST&doctor=Male&page=1&per_page=20
```

**Query Parameters:**
- `latitude` (float, optional) - Latitude for location-based search
- `longitude` (float, optional) - Longitude for location-based search
- `test` (string, optional) - Filter by available test type
- `doctor` (string, optional) - Filter by doctor gender
- `day_of_week` (string, optional) - Filter by availability day
- `page` (integer, optional) - Page number
- `per_page` (integer, optional) - Items per page

### Update Diagnostic Centre
Update diagnostic centre details.

```http
PUT /v1/diagnostic_centres/{diagnostic_centre_id}
```

### Delete Diagnostic Centre
Delete a diagnostic centre.

```http
DELETE /v1/diagnostic_centres/{diagnostic_centre_id}
```

### Get Owner's Diagnostic Centres
Get diagnostic centres owned by current user.

```http
GET /v1/diagnostic_centres/owner?page=1&per_page=20
```

### Create Diagnostic Centre Manager
Create a manager for diagnostic centres.

```http
POST /v1/diagnostic_centres/manager
```

**Request Body:**
```json
{
  "email": "manager@example.com",
  "first_name": "Jane",
  "last_name": "Smith",
  "user_type": "DIAGNOSTIC_CENTRE_MANAGER"
}
```

### Assign Manager
Assign a manager to a diagnostic centre.

```http
POST /v1/diagnostic_centres/assign
```

**Request Body:**
```json
{
  "diagnostic_centre_id": "dc-001",
  "manager_id": "manager-001"
}
```

### Unassign Manager
Remove manager from diagnostic centre.

```http
POST /v1/diagnostic_centres/unassign
```

**Request Body:**
```json
{
  "diagnostic_centre_id": "dc-001"
}
```

### Submit KYC
Submit KYC documents for diagnostic centre owner.

```http
POST /v1/diagnostic_centres_owner/kyc
```

**Request Body:**
```json
{
  "nin": "12345678901",
  "passport": "A12345678",
  "driver_licence": "DL123456789"
}
```

---

## 📅 Appointments

### Create Appointment
Create a new appointment.

```http
POST /v1/appointments
```

**Request Body:**
```json
{
  "diagnostic_centre_id": "dc-001",
  "test_type": "BLOOD_TEST",
  "appointment_date": "2025-01-20T09:00:00Z",
  "amount": 5000.00,
  "preferred_doctor": "Male",
  "payment_provider": "PAYSTACK",
  "notes": "Fasting blood test required"
}
```

**Response:**
```json
{
  "status": 201,
  "success": true,
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "diagnostic_centre_id": "dc-001",
    "patient_id": "user-001",
    "test_type": "BLOOD_TEST",
    "appointment_date": "2025-01-20T09:00:00Z",
    "status": "pending",
    "amount": 5000.00,
    "created_at": "2025-01-15T20:00:00Z"
  }
}
```

### Get Appointment
Get appointment details.

```http
GET /v1/appointments/{appointment_id}
```

### List Appointments
List appointments with filters.

```http
GET /v1/appointments?diagnostic_centre_id=dc-001&status=pending&from_date=2025-01-01&to_date=2025-01-31&page=1&per_page=20
```

**Query Parameters:**
- `diagnostic_centre_id` (string, optional) - Filter by diagnostic centre
- `patient_id` (string, optional) - Filter by patient
- `status` (string, optional) - Filter by status (pending, confirmed, in_progress, completed, cancelled, rescheduled)
- `from_date` (string, optional) - Start date filter (ISO 8601)
- `to_date` (string, optional) - End date filter (ISO 8601)
- `page` (integer, optional) - Page number
- `per_page` (integer, optional) - Items per page

### Cancel Appointment
Cancel an appointment.

```http
POST /v1/appointments/{appointment_id}/cancel
```

**Request Body:**
```json
{
  "reason": "Patient unavailable"
}
```

### Reschedule Appointment
Reschedule an appointment.

```http
POST /v1/appointments/{appointment_id}/reschedule
```

**Request Body:**
```json
{
  "new_schedule_id": "schedule-002",
  "new_date": "2025-01-22T10:00:00Z",
  "new_time_slot": "10:00-11:00",
  "reschedule_reason": "Patient requested different time"
}
```

### Confirm Appointment
Confirm an appointment with payment.

```http
POST /v1/appointments/confirm_appointment
```

**Request Body:**
```json
{
  "appointment_id": "123e4567-e89b-12d3-a456-426614174000",
  "amount": 5000.00,
  "currency": "NGN",
  "payment_method": "card",
  "payment_provider": "PAYSTACK",
  "provider_reference": "paystack-ref-123"
}
```

---

## 📋 Medical Records

### Create Medical Record
Upload a new medical record.

```http
POST /v1/medical_records
```

**Request Body:**
```json
{
  "patient_id": "user-001",
  "diagnostic_centre_id": "dc-001",
  "document_type": "LAB_REPORT",
  "file_url": "https://storage.example.com/records/report.pdf",
  "file_type": "PDF",
  "description": "Blood test results",
  "test_date": "2025-01-15T09:00:00Z"
}
```

### Get Medical Record
Get medical record details.

```http
GET /v1/medical_records/{record_id}
```

### List Medical Records
Get all medical records for current user.

```http
GET /v1/medical_records?page=1&per_page=20
```

### Update Medical Record
Update medical record details.

```http
PUT /v1/medical_records
```

**Request Body:**
```json
{
  "record_id": "rec-001",
  "description": "Updated description",
  "notes": "Additional notes"
}
```

### Get Uploader Medical Record
Get specific medical record uploaded by diagnostic centre.

```http
GET /v1/medical_records/{record_id}/diagnostic_centre/{diagnostic_centre_id}
```

### Get Uploader Medical Records
Get medical records uploaded by specific diagnostic centre.

```http
GET /v1/medical_records/diagnostic_centre/{diagnostic_centre_id}?page=1&per_page=20
```

---

## 💳 Payments

### Create Payment
Process a new payment.

```http
POST /v1/payments
```

**Request Body:**
```json
{
  "appointment_id": "123e4567-e89b-12d3-a456-426614174000",
  "amount": 5000.00,
  "currency": "NGN",
  "payment_method": "card",
  "payment_provider": "PAYSTACK",
  "provider_reference": "paystack-ref-123"
}
```

### Get Payment
Get payment details.

```http
GET /v1/payments/{payment_id}
```

### List Payments
List payments with filters.

```http
GET /v1/payments?diagnostic_centre_id=dc-001&status=success&from_date=2025-01-01&to_date=2025-01-31&page=1&per_page=20
```

**Query Parameters:**
- `diagnostic_centre_id` (string, optional) - Filter by diagnostic centre
- `patient_id` (string, optional) - Filter by patient
- `status` (string, optional) - Filter by status (pending, success, failed, refunded, cancelled)
- `from_date` (string, optional) - Start date filter
- `to_date` (string, optional) - End date filter
- `page` (integer, optional) - Page number
- `per_page` (integer, optional) - Items per page

### Refund Payment
Process payment refund.

```http
POST /v1/payments/{payment_id}/refund
```

**Request Body:**
```json
{
  "refund_amount": 5000.00,
  "refund_reason": "Appointment cancelled by patient",
  "refunded_by": "admin-001"
}
```

### Verify Payment
Verify payment status with provider.

```http
GET /v1/payments/verify/{reference}
```

### Payment Webhook
Handle payment provider webhooks.

```http
POST /v1/payments/webhook
```

**Request Body:**
```json
{
  "payment_id": "payment-001",
  "status": "success",
  "transaction_id": "txn-123456",
  "metadata": {
    "provider": "paystack",
    "reference": "ref-123"
  }
}
```

---

## 📊 Schedules

### Get Schedule
Get schedule details.

```http
GET /v1/diagnostic_schedules/{schedule_id}
```

### Update Schedule
Update schedule details.

```http
PUT /v1/diagnostic_schedules/{schedule_id}
```

**Request Body:**
```json
{
  "schedule_time": "2025-01-20T10:00:00Z",
  "status": "available",
  "notes": "Updated schedule"
}
```

### List Schedules
Get all schedules.

```http
GET /v1/diagnostic_schedules?page=1&per_page=20
```

---

## ⏰ Availability

### Create Availability
Create availability slots for diagnostic centre.

```http
POST /v1/availability
```

**Request Body:**
```json
{
  "diagnostic_centre_id": "dc-001",
  "day_of_week": "monday",
  "start_time": "09:00",
  "end_time": "17:00",
  "slot_duration": 30,
  "break_times": [
    {
      "start": "12:00",
      "end": "13:00"
    }
  ]
}
```

### Get Availability
Get availability slots for diagnostic centre.

```http
GET /v1/availability/{diagnostic_centre_id}?day_of_week=monday
```

### Update Availability
Update availability slot.

```http
PUT /v1/availability/{diagnostic_centre_id}/{day_of_week}
```

**Request Body:**
```json
{
  "start_time": "08:00",
  "end_time": "18:00",
  "slot_duration": 30
}
```

### Update Multiple Availability
Update multiple availability slots.

```http
PUT /v1/availability/{diagnostic_centre_id}
```

**Request Body:**
```json
{
  "availability": [
    {
      "day_of_week": "monday",
      "start_time": "09:00",
      "end_time": "17:00"
    },
    {
      "day_of_week": "tuesday",
      "start_time": "09:00",
      "end_time": "17:00"
    }
  ]
}
```

### Delete Availability
Delete availability slot.

```http
DELETE /v1/availability/{diagnostic_centre_id}/{day_of_week}
```

---

## 🤖 AI Services

### Interpret Lab Results
Get AI analysis of laboratory test results.

```http
POST /v1/ai/interpret_lab
```

**Request Body:**
```json
{
  "patient_id": "patient-123",
  "diagnostic_centre_id": "center-456",
  "test_name": "Complete Blood Count (CBC)",
  "results": {
    "Hemoglobin": "8.5 g/dL",
    "White Blood Cells": "12,000 /μL",
    "Platelets": "150,000 /μL"
  },
  "reference_ranges": {
    "Hemoglobin": "12.0-15.5 g/dL",
    "White Blood Cells": "4,000-11,000 /μL",
    "Platelets": "150,000-450,000 /μL"
  }
}
```

**Response:**
```json
{
  "status": 200,
  "success": true,
  "data": {
    "summary": "Blood test shows low hemoglobin and elevated white blood cells",
    "abnormal_results": [
      {
        "test_name": "Hemoglobin",
        "value": "8.5 g/dL",
        "reference_range": "12.0-15.5 g/dL",
        "severity": "moderate",
        "explanation": "Low hemoglobin may indicate anemia"
      }
    ],
    "recommendations": [
      "Consult with healthcare provider",
      "Consider iron supplementation",
      "Follow-up blood work in 4-6 weeks"
    ],
    "urgency_level": "medium",
    "follow_up_required": true
  },
  "disclaimer": "This analysis is for informational purposes only and should not replace professional medical consultation."
}
```

### Analyze Symptoms
Get AI analysis of patient symptoms.

```http
POST /v1/ai/analyze_symptoms
```

**Request Body:**
```json
{
  "symptoms": ["persistent cough", "shortness of breath", "chest pain"],
  "age": 45,
  "gender": "male"
}
```

**Response:**
```json
{
  "status": 200,
  "success": true,
  "data": {
    "possible_conditions": [
      {
        "name": "Respiratory Infection",
        "probability": "high",
        "description": "Upper or lower respiratory tract infection",
        "symptoms": ["persistent cough", "shortness of breath"]
      }
    ],
    "urgency_level": "medium",
    "recommendations": [
      "Seek medical evaluation",
      "Monitor symptoms closely",
      "Stay hydrated and rest"
    ],
    "red_flags": [
      "Severe difficulty breathing",
      "Chest pain worsening",
      "High fever"
    ],
    "next_steps": [
      "Schedule appointment with healthcare provider",
      "Consider chest X-ray if symptoms persist"
    ]
  },
  "disclaimer": "This is preliminary analysis only. Please consult with a healthcare professional for proper diagnosis and treatment."
}
```

### Generate Report Summary
Create patient-friendly or professional summaries of medical reports.

```http
POST /v1/ai/summarize_report
```

**Request Body:**
```json
{
  "medical_report": "CHEST X-RAY REPORT\n\nPatient: John Doe, Age: 45, Male\nDate: 2024-01-15\n\nFINDINGS:\n- Bilateral lower lobe infiltrates consistent with pneumonia\n- No pleural effusion\n- Heart size within normal limits",
  "patient_friendly": true
}
```

**Response:**
```json
{
  "status": 200,
  "success": true,
  "summary": "Your chest X-ray shows signs of pneumonia in both lower parts of your lungs. The good news is that your heart appears normal and there's no fluid around your lungs. This type of infection is treatable with antibiotics. Please follow up with your doctor for treatment.",
  "type": "patient-friendly"
}
```

### Get AI Capabilities
Get list of available AI features.

```http
GET /v1/ai/capabilities
```

**Response:**
```json
{
  "status": 200,
  "success": true,
  "data": {
    "features": [
      {
        "name": "Lab Result Interpretation",
        "endpoint": "/v1/ai/interpret-lab",
        "description": "AI-powered analysis of laboratory test results",
        "method": "POST"
      },
      {
        "name": "Symptom Analysis",
        "endpoint": "/v1/ai/analyze-symptoms",
        "description": "Preliminary analysis of patient symptoms",
        "method": "POST"
      },
      {
        "name": "Report Summarization",
        "endpoint": "/v1/ai/summarize-report",
        "description": "Generate patient-friendly or professional summaries of medical reports",
        "method": "POST"
      }
    ],
    "disclaimer": "All AI features are for informational purposes only and should not replace professional medical consultation.",
    "version": "1.0"
  }
}
```

---

## ⚡ WebSocket

### Establish WebSocket Connection
Connect to real-time notification system.

```javascript
// JavaScript WebSocket Connection
const ws = new WebSocket('ws://localhost:7556/v1/ws/notifications?user_id=user-123');

ws.onopen = function(event) {
    console.log('Connected to WebSocket');
};

ws.onmessage = function(event) {
    const notification = JSON.parse(event.data);
    console.log('Received notification:', notification);
};
```

**Connection URL:**
```
ws://localhost:7556/v1/ws/notifications?user_id={user_id}
```

**Message Format:**
```json
{
  "type": "appointment_confirmed",
  "user_id": "user-123",
  "title": "Appointment Confirmed",
  "message": "Your appointment has been confirmed for tomorrow at 2:00 PM",
  "data": {
    "appointment_id": "appointment-456",
    "type": "appointment"
  },
  "timestamp": "2025-01-15T10:30:00Z",
  "id": "notification-789"
}
```

### Get WebSocket Statistics
Get current WebSocket connection statistics.

```http
GET /v1/ws/stats
```

**Response:**
```json
{
  "status": 200,
  "success": true,
  "data": {
    "connected_clients": 25,
    "status": "active"
  }
}
```

### Send Test Notification
Send test notification via WebSocket (for testing).

```http
POST /v1/ws/test-notification
```

**Request Body:**
```json
{
  "user_id": "user-123",
  "title": "Test Notification",
  "message": "This is a test notification",
  "data": {
    "test": true,
    "timestamp": "2025-01-15T10:30:00Z"
  }
}
```

---

## 🔔 Notifications

### Get Notifications
Get user notifications with pagination.

```http
GET /v1/notifications?page=1&per_page=20
```

**Response:**
```json
{
  "status": 200,
  "success": true,
  "data": {
    "notifications": [
      {
        "id": "notif-001",
        "type": "APPOINTMENT_CONFIRMED",
        "title": "Appointment Confirmed",
        "message": "Your appointment has been confirmed",
        "read": false,
        "created_at": "2025-01-15T10:30:00Z",
        "metadata": {
          "appointment_id": "appointment-456"
        }
      }
    ],
    "pagination": {
      "page": 1,
      "per_page": 20,
      "total": 45,
      "total_pages": 3
    }
  }
}
```

### Mark Notification as Read
Mark a specific notification as read.

```http
PUT /v1/notifications/{notification_id}/read
```

**Response:**
```json
{
  "status": 200,
  "success": true,
  "data": {
    "message": "Notification marked as read"
  }
}
```

---

## 📊 Health & Monitoring

### Health Check
Check API health status.

```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-01-15T10:30:00Z",
  "version": "1.0.0",
  "services": {
    "database": "healthy",
    "websocket": "healthy",
    "ai": "healthy"
  }
}
```

### Metrics
Get Prometheus metrics (for monitoring).

```http
GET /metrics
```

### API Documentation
Access Swagger documentation.

```http
GET /swagger/
```

---

## 📝 Data Models

### User Types
- `PATIENT` - Regular patients
- `DIAGNOSTIC_CENTRE_OWNER` - Diagnostic centre owners
- `DIAGNOSTIC_CENTRE_MANAGER` - Diagnostic centre managers

### Test Types
Available medical test types:
- `BLOOD_TEST`, `URINE_TEST`, `X_RAY`, `MRI`, `CT_SCAN`
- `ULTRASOUND`, `ECG`, `EEG`, `BIOPSY`, `SKIN_TEST`
- `IMMUNOLOGY_TEST`, `HORMONE_TEST`, `VIRAL_TEST`
- `BACTERIAL_TEST`, `PARASITIC_TEST`, `FUNGAL_TEST`
- `MOLECULAR_TEST`, `TOXICOLOGY_TEST`, `ECHO`
- `COVID_19_TEST`, `BLOOD_SUGAR_TEST`, `LIPID_PROFILE`
- `HEMOGLOBIN_TEST`, `THYROID_TEST`, `LIVER_FUNCTION_TEST`
- `KIDNEY_FUNCTION_TEST`, `URIC_ACID_TEST`, `VITAMIN_D_TEST`
- `VITAMIN_B12_TEST`, `HEMOGRAM`, `COMPLETE_BLOOD_COUNT`
- `BLOOD_GROUPING`, `HEPATITIS_B_TEST`, `HEPATITIS_C_TEST`
- `HIV_TEST`, `MALARIA_TEST`, `DENGUE_TEST`, `TYPHOID_TEST`
- `COVID_19_ANTIBODY_TEST`, `COVID_19_RAPID_ANTIGEN_TEST`
- `COVID_19_RT_PCR_TEST`, `PREGNANCY_TEST`, `ALLERGY_TEST`
- `GENETIC_TEST`, `OTHER`

### Appointment Status
- `pending` - Appointment created, awaiting confirmation
- `confirmed` - Appointment confirmed and paid
- `in_progress` - Appointment currently happening
- `completed` - Appointment finished
- `cancelled` - Appointment cancelled
- `rescheduled` - Appointment rescheduled

### Payment Status
- `pending` - Payment initiated
- `success` - Payment successful
- `failed` - Payment failed
- `refunded` - Payment refunded
- `cancelled` - Payment cancelled

### Payment Methods
- `card` - Credit/debit card
- `transfer` - Bank transfer
- `cash` - Cash payment
- `wallet` - Digital wallet

### Payment Providers
- `PAYSTACK` - Paystack payment gateway
- `FLUTTERWAVE` - Flutterwave payment gateway
- `STRIPE` - Stripe payment gateway
- `MONNIFY` - Monnify payment gateway

---

## 🔧 SDK & Integration Examples

### JavaScript/Node.js Example
```javascript
const DiagnoxixAPI = {
  baseURL: 'https://diagnoxix.onrender.com',
  token: null,

  setToken(token) {
    this.token = token;
  },

  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const headers = {
      'Content-Type': 'application/json',
      ...(this.token && { Authorization: `Bearer ${this.token}` }),
      ...options.headers
    };

    const response = await fetch(url, {
      ...options,
      headers
    });

    return response.json();
  },

  // Authentication
  async login(email, password) {
    const response = await this.request('/v1/login', {
      method: 'POST',
      body: JSON.stringify({ email, password })
    });
    
    if (response.success) {
      this.setToken(response.data.token);
    }
    
    return response;
  },

  // Appointments
  async createAppointment(appointmentData) {
    return this.request('/v1/appointments', {
      method: 'POST',
      body: JSON.stringify(appointmentData)
    });
  },

  // AI Services
  async interpretLabResults(labData) {
    return this.request('/v1/ai/interpret-lab', {
      method: 'POST',
      body: JSON.stringify(labData)
    });
  }
};

// Usage
const api = DiagnoxixAPI;
await api.login('user@example.com', 'password');
const appointment = await api.createAppointment({
  diagnostic_centre_id: 'dc-001',
  test_type: 'BLOOD_TEST',
  appointment_date: '2025-01-20T09:00:00Z',
  amount: 5000.00
});
```

### Python Example
```python
import requests
import json

class DiagnoxixAPI:
    def __init__(self, base_url='https://diagnoxix.onrender.com'):
        self.base_url = base_url
        self.token = None
        self.session = requests.Session()
    
    def set_token(self, token):
        self.token = token
        self.session.headers.update({'Authorization': f'Bearer {token}'})
    
    def login(self, email, password):
        response = self.session.post(f'{self.base_url}/v1/login', 
                                   json={'email': email, 'password': password})
        data = response.json()
        if data.get('success'):
            self.set_token(data['data']['token'])
        return data
    
    def create_appointment(self, appointment_data):
        response = self.session.post(f'{self.base_url}/v1/appointments', 
                                   json=appointment_data)
        return response.json()
    
    def interpret_lab_results(self, lab_data):
        response = self.session.post(f'{self.base_url}/v1/ai/interpret-lab', 
                                   json=lab_data)
        return response.json()

# Usage
api = DiagnoxixAPI()
api.login('user@example.com', 'password')
appointment = api.create_appointment({
    'diagnostic_centre_id': 'dc-001',
    'test_type': 'BLOOD_TEST',
    'appointment_date': '2025-01-20T09:00:00Z',
    'amount': 5000.00
})
```

---

## 🚀 Getting Started

1. **Register an account** using `/v1/register`
2. **Verify your email** using the verification link
3. **Login** using `/v1/login` to get your JWT token
4. **Include the token** in the Authorization header for protected endpoints
5. **Start making API calls** to manage appointments, records, and more

## 📞 Support

- **Documentation**: [API Docs](https://diagnoxix.onrender.com/swagger/)
- **Email**: support@diagnoxix.com
- **Status Page**: [status.diagnoxix.com](https://status.diagnoxix.com)

---

## 📄 License

This API documentation is proprietary to Diagnoxix. All rights reserved.

---

*Last updated: January 15, 2025*
