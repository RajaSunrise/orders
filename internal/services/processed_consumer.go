package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/rajasunsire/orders/internal/models"
	kafka "github.com/segmentio/kafka-go"
)

// ProcessProcessedOrders consumes from processed-orders topic and saves to DB
func ProcessProcessedOrders(consumer *kafka.Reader) {
	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading processed message: %v", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Error unmarshaling processed order: %v", err)
			continue
		}

		// Log the received processed order
		log.Printf("Received Processed Order: ID=%s, User=%s, Product=%s, Status=%s | Time: %s",
			order.ID.String(), order.UserID.String(), order.ProductID.String(), order.Status, time.Now().Format(time.RFC3339))

		// Save or update in DB
		if err := DB.Where("id = ?", order.ID).Assign(order).FirstOrCreate(&order).Error; err != nil {
			log.Printf("Error saving processed order to DB: %v", err)
		}
	}
}
