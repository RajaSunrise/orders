package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	kafka "github.com/segmentio/kafka-go"
)

// Order struct untuk pesan dari Kafka
type Order struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Product     string  `json:"product"`
	Quantity    int     `json:"quantity"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"` // Diupdate setelah validasi
}

// Data referensi untuk validasi
var validUsers = map[string]bool{
	"user123": true,
	"user456": true,
}
var validProducts = map[string]bool{
	"laptop":     true,
	"smartphone": true,
}

func main() {
	// Setup Fiber untuk API
	app := fiber.New()

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Order Processing Service is running",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Jalankan Fiber di goroutine
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("Failed to start Fiber: %v", err)
		}
	}()

	// Setup Kafka consumer (connect ke Kafka KRaft di Podman)
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:29092"}, // Sesuai EXTERNAL listener di docker-compose.yml
		Topic:     "orders-topic",
		GroupID:   "order-processor-group",
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   1 * time.Second,
		Dialer:    &kafka.Dialer{Timeout: 10 * time.Second},
	})

	// Setup Kafka producer untuk hasil validasi
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:29092"},
		Topic:   "processed-orders",
	})

	// Handle graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Consumer loop
	go func() {
		for {
			// Baca pesan
			msg, err := consumer.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			// Parse JSON ke Order
			var order Order
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				log.Printf("Error unmarshaling order: %v", err)
				continue
			}

			// Validasi order
			isValid := validateOrder(order)
			order.Status = "rejected"
			if isValid {
				order.Status = "approved"
			}

			// Log hasil
			log.Printf("Processed Order: ID=%s, User=%s, Product=%s, Status=%s | Time: %s",
				order.ID, order.UserID, order.Product, order.Status, time.Now().Format(time.RFC3339))

			// Kirim hasil ke topic processed-orders
			processedMsg, err := json.Marshal(order)
			if err != nil {
				log.Printf("Error marshaling processed order: %v", err)
				continue
			}
			if err := producer.WriteMessages(context.Background(), kafka.Message{Value: processedMsg}); err != nil {
				log.Printf("Error writing to processed-orders: %v", err)
			}
		}
	}()

	// Tunggu sinyal shutdown
	<-signals
	log.Println("Shutting down...")

	// Tutup consumer dan producer
	if err := consumer.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}
	if err := producer.Close(); err != nil {
		log.Printf("Error closing producer: %v", err)
	}

	// Shutdown Fiber
	if err := app.Shutdown(); err != nil {
		log.Printf("Error shutting down Fiber: %v", err)
	}
}

// Validasi order
func validateOrder(order Order) bool {
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