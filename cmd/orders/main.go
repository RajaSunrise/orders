package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rajasunsire/orders/internal/handlers"
	"github.com/rajasunsire/orders/internal/services"
)

func main() {
	// Initialize database
	services.InitDB()

	// Initialize Kafka producers
	brokers := []string{"localhost:29092"}
	services.InitKafkaProducers(brokers)

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

	// CRUD endpoints for orders
	app.Get("/orders", handlers.GetOrders)
	app.Get("/orders/:id", handlers.GetOrder)
	app.Post("/orders", handlers.CreateOrder)
	app.Put("/orders/:id", handlers.UpdateOrder)
	app.Delete("/orders/:id", handlers.DeleteOrder)

	// Jalankan Fiber di goroutine
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("Failed to start Fiber: %v", err)
		}
	}()

	// Setup Kafka consumers
	consumer := services.NewKafkaReader(brokers, "orders-topic", "order-processor-group")
	processedConsumer := services.NewKafkaReader(brokers, "processed-orders", "processed-orders-consumer-group")

	// Handle graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Start order processing service
	go services.ProcessOrders(consumer)

	// Start processed orders consumer service
	go services.ProcessProcessedOrders(processedConsumer)

	// Tunggu sinyal shutdown
	<-signals
	log.Println("Shutting down...")

	// Tutup consumers dan producers
	if err := consumer.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}
	if err := services.OrdersProducer.Close(); err != nil {
		log.Printf("Error closing orders producer: %v", err)
	}
	if err := services.ProcessedProducer.Close(); err != nil {
		log.Printf("Error closing processed producer: %v", err)
	}
	if err := processedConsumer.Close(); err != nil {
		log.Printf("Error closing processed consumer: %v", err)
	}

	// Shutdown Fiber
	if err := app.Shutdown(); err != nil {
		log.Printf("Error shutting down Fiber: %v", err)
	}
}
