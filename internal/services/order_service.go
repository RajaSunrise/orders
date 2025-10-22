package services

import (
	"errors"
	"github.com/rajasunsire/orders/internal/models"
)

// validateOrder validates the order
func validateOrder(order models.Order) bool {
	validUsers := map[string]bool{
		"user123": true,
		"user456": true,
	}
	validProducts := map[string]bool{
		"laptop":     true,
		"smartphone": true,
	}

	if _, userOk := validUsers[order.UserID]; !userOk {
		return false
	}
	if _, productOk := validProducts[order.Product]; !productOk {
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
	result := DB.Where("id = ?", id).First(&order)
	if result.Error != nil {
		return models.Order{}, errors.New("order not found")
	}
	return order, nil
}

// CreateOrder creates a new order
func CreateOrder(order models.Order) error {
	if order.ID == "" {
		return errors.New("order ID is required")
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
	var order models.Order
	result := DB.Where("id = ?", id).First(&order)
	if result.Error != nil {
		return errors.New("order not found")
	}
	updatedOrder.ID = id
	DB.Save(&updatedOrder)
	return nil
}

// DeleteOrder deletes an order by ID
func DeleteOrder(id string) error {
	result := DB.Where("id = ?", id).Delete(&models.Order{})
	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}
