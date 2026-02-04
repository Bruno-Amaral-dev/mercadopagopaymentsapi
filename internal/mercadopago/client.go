package mercadopago

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents the MercadoPago API client
type Client struct {
	AccessToken string
	BaseURL     string
	HTTPClient  *http.Client
}

// NewClient creates a new MercadoPago client
func NewClient(accessToken, baseURL string) *Client {
	return &Client{
		AccessToken: accessToken,
		BaseURL:     baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateCardToken generates a card token for payment processing
func (c *Client) GenerateCardToken(cardNumber, holderName, docType, docNumber, securityCode, expirationMonth, expirationYear string) (*MercadoPagoCardTokenResponse, error) {
	tokenRequest := MercadoPagoCardTokenRequest{
		CardNumber:      cardNumber,
		SecurityCode:    securityCode,
		ExpirationMonth: expirationMonth,
		ExpirationYear:  expirationYear,
		Cardholder: struct {
			Name           string `json:"name"`
			Identification struct {
				Type   string `json:"type"`
				Number string `json:"number"`
			} `json:"identification"`
		}{
			Name: holderName,
			Identification: struct {
				Type   string `json:"type"`
				Number string `json:"number"`
			}{
				Type:   docType,
				Number: docNumber,
			},
		},
	}

	jsonData, err := json.Marshal(tokenRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal token request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/v1/card_tokens", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MercadoPago API error: status %d, response: %s", resp.StatusCode, string(body))
	}

	var tokenResponse MercadoPagoCardTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &tokenResponse, nil
}

// generateIdempotencyKey generates a unique idempotency key for the request
func generateIdempotencyKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreatePayment creates a payment with the given token and payment details
func (c *Client) CreatePayment(token, paymentMethodID string, amount float64, installments int, docType, docNumber, email string) (*MercadoPagoPaymentResponse, error) {
	paymentRequest := MercadoPagoPaymentRequest{
		TransactionAmount: amount,
		Installments:      installments,
		PaymentMethodID:   paymentMethodID,
		Token:             token,
		Payer: struct {
			Email          string `json:"email"`
			Type           string `json:"type"`
			Identification struct {
				Type   string `json:"type"`
				Number string `json:"number"`
			} `json:"identification"`
		}{
			Email: email,
			Type:  "customer",
			Identification: struct {
				Type   string `json:"type"`
				Number string `json:"number"`
			}{
				Type:   docType,
				Number: docNumber,
			},
		},
	}

	jsonData, err := json.Marshal(paymentRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/v1/payments", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	// Generate and set idempotency key
	idempotencyKey, err := generateIdempotencyKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate idempotency key: %w", err)
	}
	req.Header.Set("X-Idempotency-Key", idempotencyKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MercadoPago API error: status %d, response: %s", resp.StatusCode, string(body))
	}

	var paymentResponse MercadoPagoPaymentResponse
	if err := json.Unmarshal(body, &paymentResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &paymentResponse, nil
}
