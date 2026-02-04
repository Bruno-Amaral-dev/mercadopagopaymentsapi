package main

import (
	"log"
	"mercadopago-demo/internal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(internal.PaymentResponse{
				Success: false,
				Error:   err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load environment variables: %v", err)
	}

	// Get MercadoPago credentials from environment variables
	accessToken := internal.GetEnv("MERCADOPAGO_ACCESS_TOKEN", "")
	baseURL := internal.GetEnv("MERCADOPAGO_BASE_URL", "")

	if accessToken == "" {
		log.Fatal("MERCADOPAGO_ACCESS_TOKEN environment variable is required")
	}

	// Initialize payment handler
	paymentHandler := internal.NewPaymentHandler(accessToken, baseURL)

	// Routes
	app.Get("/health", paymentHandler.HealthCheck)
	app.Post("/api/payments", paymentHandler.ProcessPayment)
	app.Post("/api/webhooks", paymentHandler.HandleWebhook)

	// Start server
	port := internal.GetEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)
	log.Printf("MercadoPago API URL: %s", baseURL)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
