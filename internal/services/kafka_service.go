package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/rajasunsire/orders/internal/models"
	kafka "github.com/segmentio/kafka-go"
)

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers []string
}

// Global producers
var OrdersProducer *kafka.Writer
var ProcessedProducer *kafka.Writer

// InitKafkaProducers initializes Kafka producers
func InitKafkaProducers(brokers []string) {
	OrdersProducer = NewKafkaWriter(brokers, "orders-topic")
	ProcessedProducer = NewKafkaWriter(brokers, "processed-orders")
}

// SendOrderToKafka sends order to orders-topic
func SendOrderToKafka(order models.Order) error {
	msg, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return OrdersProducer.WriteMessages(context.Background(), kafka.Message{Value: msg})
}

// NewKafkaWriter creates a new Kafka writer
func NewKafkaWriter(brokers []string, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
}

// NewKafkaReader creates a new Kafka reader
func NewKafkaReader(brokers []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  1 * time.Second,
		Dialer:   &kafka.Dialer{Timeout: 10 * time.Second},
	})
}

// ProcessOrders consumes from orders-topic, validates, and produces to processed-orders
func ProcessOrders(consumer *kafka.Reader) {
	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Error unmarshaling order: %v", err)
			continue
		}

		// Validate order
		isValid := validateOrder(order)
		order.Status = "rejected"
		if isValid {
			order.Status = "approved"
		}

		log.Printf("Processed Order: ID=%s, User=%s, Product=%s, Status=%s | Time: %s",
			order.ID, order.UserID, order.Product, order.Status, time.Now().Format(time.RFC3339))

		// Send to processed-orders
		processedMsg, err := json.Marshal(order)
		if err != nil {
			log.Printf("Error marshaling processed order: %v", err)
			continue
		}
		if err := ProcessedProducer.WriteMessages(context.Background(), kafka.Message{Value: processedMsg}); err != nil {
			log.Printf("Error writing to processed-orders: %v", err)
		}
	}
}
