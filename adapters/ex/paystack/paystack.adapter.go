package paystack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type PaystackConfig struct {
	SecretKey string
	BaseURL   string
}

type PaystackTransactionResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

type PaystackVerificationResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Status    string  `json:"status"`
		Reference string  `json:"reference"`
		Amount    float64 `json:"amount"`
		Channel   string  `json:"channel"`
		Currency  string  `json:"currency"`
		PaidAt    string  `json:"paid_at"`
		CreatedAt string  `json:"created_at"`
		Customer  struct {
			ID           int    `json:"id"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Email        string `json:"email"`
			CustomerCode string `json:"customer_code"`
			Phone        string `json:"phone"`
		} `json:"customer"`
		Authorization struct {
			AuthorizationCode string `json:"authorization_code"`
			CardType          string `json:"card_type"`
			Last4             string `json:"last4"`
			ExpMonth          string `json:"exp_month"`
			ExpYear           string `json:"exp_year"`
			Bank              string `json:"bank"`
			Brand             string `json:"brand"`
		} `json:"authorization"`
		Fees     float64                `json:"fees"`
		Metadata map[string]interface{} `json:"metadata"`
	} `json:"data"`
}

type Customer struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	CustomerCode string `json:"customer_code"`
	Phone        string `json:"phone"`
}

type PaystackAdapter struct {
	config *PaystackConfig
}

func NewPaystackAdapter(con *PaystackConfig) *PaystackAdapter {
	return &PaystackAdapter{
		config: con,
	}
}

func (p *PaystackAdapter) InitializeTransaction(email string, amount float64, reference string, metadata map[string]interface{}) (*PaystackTransactionResponse, error) {
	url := fmt.Sprintf("%s/transaction/initialize", p.config.BaseURL)

	// Paystack expects amount in kobo (multiply by 100)
	amountInKobo := int64(amount * 100)

	payload := map[string]interface{}{
		"email":     email,
		"amount":    amountInKobo,
		"reference": reference,
		"metadata":  metadata,
		"currency":  "NGN",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.SecretKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PaystackTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// VerifyTransaction verifies a transaction using its reference
func (p *PaystackAdapter) VerifyTransaction(reference string) (*PaystackVerificationResponse, error) {
	url := fmt.Sprintf("%s/transaction/verify/%s", p.config.BaseURL, reference)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.SecretKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result PaystackVerificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}
