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
