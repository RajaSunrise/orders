package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rajasunsire/orders/internal/models"
)

// validateOrder validates the order
func validateOrder(order models.Order) bool {
	var user models.User
	result := DB.Where("id = ?", order.UserID).First(&user)
	if result.Error != nil {
		return false
	}
	var product models.Product
	result = DB.Where("id = ?", order.ProductID).First(&product)
	if result.Error != nil {
		return false
	}
	if order.Quantity <= 0 || order.TotalAmount <= 0 {
		return false
	}
	return true
}

// GetAllOrders returns all orders
func GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	result := DB.Find(&orders)
	return orders, result.Error
}

// GetOrderByID returns order by ID
func GetOrderByID(id string) (models.Order, error) {
	var order models.Order
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.Order{}, errors.New("invalid ID")
	}
	result := DB.Where("id = ?", parsedID).First(&order)
	if result.Error != nil {
		return models.Order{}, errors.New("order not found")
	}
	return order, nil
}

// CreateOrder creates a new order
func CreateOrder(order models.Order) error {
	if order.ID.String() == "00000000-0000-0000-0000-000000000000" {
		order.ID = uuid.New()
	}
	// Validate and set status
	if validateOrder(order) {
		order.Status = "approved"
	} else {
		order.Status = "rejected"
	}
	result := DB.Create(&order)
	return result.Error
}

// UpdateOrder updates an existing order
func UpdateOrder(id string, updatedOrder models.Order) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid ID")
	}
	var order models.Order
	result := DB.Where("id = ?", parsedID).First(&order)
	if result.Error != nil {
		return errors.New("order not found")
	}
	updatedOrder.ID = parsedID
	DB.Save(&updatedOrder)
	return nil
}

// DeleteOrder deletes an order by ID
func DeleteOrder(id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid ID")
	}
	result := DB.Where("id = ?", parsedID).Delete(&models.Order{})
	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}
