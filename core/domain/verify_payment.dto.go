package domain

import "time"

// PaystackVerificationResponse represents a payment verification response from Paystack
type PaystackVerificationResponse struct {
	Status  bool                     `json:"status"`
	Message string                   `json:"message"`
	Data    PaystackVerificationData `json:"data"`
}

// PaystackVerificationData contains the verification response data
type PaystackVerificationData struct {
	Status    string                 `json:"status"`
	Reference string                 `json:"reference"`
	Amount    float64                `json:"amount"`
	Currency  string                 `json:"currency"`
	Channel   string                 `json:"channel"`
	PaidAt    time.Time              `json:"paid_at"`
	CreatedAt time.Time              `json:"created_at"`
	Customer  PaystackCustomer       `json:"customer"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// PaystackCustomer contains customer information from Paystack
type PaystackCustomer struct {
	ID        int                    `json:"id"`
	FirstName string                 `json:"first_name"`
	LastName  string                 `json:"last_name"`
	Email     string                 `json:"email"`
	Phone     string                 `json:"phone"`
	Metadata  map[string]interface{} `json:"metadata"`
}
