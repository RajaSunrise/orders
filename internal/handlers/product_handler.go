package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rajasunsire/orders/internal/models"
	"github.com/rajasunsire/orders/internal/services"
)

// GetProducts handles GET /products
func GetProducts(c *fiber.Ctx) error {
	products, err := services.GetAllProducts()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve products"})
	}
	return c.JSON(products)
}

// GetProduct handles GET /products/:id
func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	product, err := services.GetProductByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(product)
}

// CreateProduct handles POST /products
func CreateProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	if err := services.CreateProduct(product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(product)
}

// UpdateProduct handles PUT /products/:id
func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var updatedProduct models.Product
	if err := c.BodyParser(&updatedProduct); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	if err := services.UpdateProduct(id, updatedProduct); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updatedProduct)
}

// DeleteProduct handles DELETE /products/:id
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := services.DeleteProduct(id); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Product deleted"})
}
