package models

import "time"

// Order struct untuk pesan dari Kafka
type Order struct {
	ID          string  `json:"id" gorm:"primaryKey"`
	UserID      string  `json:"user_id"`
	Product     string  `json:"product"`
	Quantity    int     `json:"quantity"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
	CreatedAt   time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}
