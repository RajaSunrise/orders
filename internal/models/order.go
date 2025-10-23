package models

import (
	"time"

	"github.com/google/uuid"
)

// Order struct for message from kafka
type Order struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID      uuid.UUID `json:"user_id"`
	ProductID   uuid.UUID `json:"product_id"`
	Quantity    int       `json:"quantity"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
