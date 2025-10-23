package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Health check endpoint
func HealthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Order Processing Service is running",
		"time":    time.Now().Format(time.RFC3339),
	})
}
