package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/rajasunsire/orders/internal/routes"
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
	
	// setup Routes
	routes.Routes(app)

	// Running Fiber di goroutine
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

	// A Wait signal shutdown
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
