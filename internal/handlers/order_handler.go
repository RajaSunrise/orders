package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rajasunsire/orders/internal/models"
	"github.com/rajasunsire/orders/internal/services"
)

// GetOrders handles GET /orders
func GetOrders(c *fiber.Ctx) error {
	orders, err := services.GetAllOrders()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve orders"})
	}
	return c.JSON(orders)
}

// GetOrder handles GET /orders/:id
func GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	order, err := services.GetOrderByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(order)
}

// CreateOrder handles POST /orders
func CreateOrder(c *fiber.Ctx) error {
	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	if err := services.CreateOrder(order); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Send to Kafka for processing
	if err := services.SendOrderToKafka(order); err != nil {
		log.Printf("Error sending order to Kafka: %v", err)
		// Don't fail the request, just log
	}
	return c.Status(201).JSON(order)
}

// UpdateOrder handles PUT /orders/:id
func UpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var updatedOrder models.Order
	if err := c.BodyParser(&updatedOrder); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	if err := services.UpdateOrder(id, updatedOrder); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updatedOrder)
}

// DeleteOrder handles DELETE /orders/:id
func DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := services.DeleteOrder(id); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Order deleted"})
}
