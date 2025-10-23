package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rajasunsire/orders/internal/models"
)

// GetAllUsers returns all users
func GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := DB.Find(&users)
	return users, result.Error
}

// GetUserByID returns user by ID
func GetUserByID(id string) (models.User, error) {
	var user models.User
	result := DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

// CreateUser creates a new user
func CreateUser(user models.User) error {
	if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
		user.ID = uuid.New()
	}
	result := DB.Create(&user)
	return result.Error
}

// UpdateUser updates an existing user
func UpdateUser(id string, updatedUser models.User) error {
	var user models.User
	result := DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return errors.New("user not found")
	}
	updatedUser.ID = user.ID
	DB.Save(&updatedUser)
	return nil
}

// DeleteUser deletes a user by ID
func DeleteUser(id string) error {
	result := DB.Where("id = ?", id).Delete(&models.User{})
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
