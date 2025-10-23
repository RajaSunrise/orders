package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rajasunsire/orders/internal/models"
)

// GetAllProducts returns all products
func GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	result := DB.Find(&products)
	return products, result.Error
}

// GetProductByID returns product by ID
func GetProductByID(id string) (models.Product, error) {
	var product models.Product
	result := DB.Where("id = ?", id).First(&product)
	if result.Error != nil {
		return models.Product{}, errors.New("product not found")
	}
	return product, nil
}

// CreateProduct creates a new product
func CreateProduct(product models.Product) error {
	if product.ID.String() == "00000000-0000-0000-0000-000000000000" {
		product.ID = uuid.New()
	}
	result := DB.Create(&product)
	return result.Error
}

// UpdateProduct updates an existing product
func UpdateProduct(id string, updatedProduct models.Product) error {
	var product models.Product
	result := DB.Where("id = ?", id).First(&product)
	if result.Error != nil {
		return errors.New("product not found")
	}
	updatedProduct.ID = product.ID
	DB.Save(&updatedProduct)
	return nil
}

// DeleteProduct deletes a product by ID
func DeleteProduct(id string) error {
	result := DB.Where("id = ?", id).Delete(&models.Product{})
	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}
