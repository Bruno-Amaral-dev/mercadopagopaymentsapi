package internal

import (
	"log"
	"os"

	"mercadopago-demo/internal/mercadopago"

	"github.com/gofiber/fiber/v2"
)

// PaymentHandler handles payment requests
type PaymentHandler struct {
	mpClient *mercadopago.Client
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(accessToken, baseURL string) *PaymentHandler {
	return &PaymentHandler{
		mpClient: mercadopago.NewClient(accessToken, baseURL),
	}
}

// ProcessPayment handles the payment processing endpoint
func (h *PaymentHandler) ProcessPayment(c *fiber.Ctx) error {
	var request PaymentRequest

	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(PaymentResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
	}

	// Validate required fields
	if request.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(PaymentResponse{
			Success: false,
			Error:   "Amount must be greater than 0",
		})
	}

	if request.DocType == "" || request.DocNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(PaymentResponse{
			Success: false,
			Error:   "Document type and number are required",
		})
	}

	if request.Card.Number == "" || request.Card.HolderName == "" ||
		request.Card.ExpirationMonth == "" || request.Card.ExpirationYear == "" ||
		request.Card.SecurityCode == "" || request.Card.PaymentMethods == "" {
		return c.Status(fiber.StatusBadRequest).JSON(PaymentResponse{
			Success: false,
			Error:   "All card fields are required",
		})
	}

	if request.Installments <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(PaymentResponse{
			Success: false,
			Error:   "Installments must be greater than 0",
		})
	}

	log.Printf("Processing payment for amount: %d, doc: %s %s", request.Amount, request.DocType, request.DocNumber)

	// Generate card token
	tokenResponse, err := h.mpClient.GenerateCardToken(
		request.Card.Number,
		request.Card.HolderName,
		request.DocType,
		request.DocNumber,
		request.Card.SecurityCode,
		request.Card.ExpirationMonth,
		request.Card.ExpirationYear,
	)

	if err != nil {
		log.Printf("Error generating card token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(PaymentResponse{
			Success: false,
			Error:   "Failed to generate card token: " + err.Error(),
		})
	}

	log.Printf("Card token generated successfully: %s", tokenResponse.ID)

	// Create payment
	paymentResponse, err := h.mpClient.CreatePayment(
		tokenResponse.ID,
		request.Card.PaymentMethods,
		float64(request.Amount),
		request.Installments,
		request.DocType,
		request.DocNumber,
		"user@example.com",
	)

	if err != nil {
		log.Printf("Error creating payment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(PaymentResponse{
			Success: false,
			Error:   "Failed to create payment: " + err.Error(),
		})
	}

	log.Printf("Payment created successfully with ID: %d, Status: %s", paymentResponse.ID, paymentResponse.Status)

	// Format response data
	responseData := map[string]interface{}{
		"payment_id":     paymentResponse.ID,
		"status":         paymentResponse.Status,
		"status_detail":  paymentResponse.StatusDetail,
		"amount":         paymentResponse.TransactionAmount,
		"installments":   paymentResponse.Installments,
		"payment_method": paymentResponse.PaymentMethodID,
		"card_last_four": paymentResponse.Card.LastFourDigits,
		"created_at":     paymentResponse.DateCreated,
	}

	return c.JSON(PaymentResponse{
		Success: true,
		Message: "Payment processed successfully",
		Data:    responseData,
	})
}

// HealthCheck handles health check requests
func (h *PaymentHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "mercadopago-payment-api",
	})
}

// GetEnv gets environment variable with fallback
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Add to handlers.go
func (h *PaymentHandler) HandleWebhook(c *fiber.Ctx) error {
	var notification map[string]interface{}
	if err := c.BodyParser(&notification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(PaymentResponse{
			Success: false,
			Error:   "Invalid webhook payload",
		})
	}

	// Verify webhook signature (optional but recommended)
	signature := c.Get("x-signature")
	if signature == "" {
		log.Printf("Warning: Webhook received without signature")
	}

	// Process the notification
	log.Printf("Received webhook notification: %+v", notification)

	// Handle payment status updates here
	// Update your database, send notifications, etc.

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "received",
	})
}
