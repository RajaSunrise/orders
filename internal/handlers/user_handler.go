package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rajasunsire/orders/internal/models"
	"github.com/rajasunsire/orders/internal/services"
)

// GetUsers handles GET /users
func GetUsers(c *fiber.Ctx) error {
	users, err := services.GetAllUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve users"})
	}
	return c.JSON(users)
}

// GetUser handles GET /users/:id
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := services.GetUserByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// CreateUser handles POST /users
func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	if err := services.CreateUser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(user)
}

// UpdateUser handles PUT /users/:id
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var updatedUser models.User
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	if err := services.UpdateUser(id, updatedUser); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updatedUser)
}

// DeleteUser handles DELETE /users/:id
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := services.DeleteUser(id); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "User deleted"})
}
