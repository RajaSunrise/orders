# Order Service with Kafka and CRUD

A Go microservice for processing orders with Kafka integration, CRUD API, and PostgreSQL persistence.

## Features
- REST API for order CRUD operations
- Kafka consumer for order processing
- Separate Kafka consumer for processed orders
- PostgreSQL database with GORM
- Docker Compose setup for Kafka and PostgreSQL

## Setup

1. Start services:
```bash
docker-compose up -d
```

2. Create Kafka topics:
```bash
# orders-topic
docker exec -it kafka kafka-topics.sh --create --topic orders-topic --bootstrap-server kafka:9092

# processed-orders topic
docker exec -it kafka kafka-topics.sh --create --topic processed-orders --bootstrap-server kafka:9092
```

3. Run the application:
```bash
go run cmd/orders/main.go
```

## API Endpoints

- `GET /health` - Health check
- `GET /orders` - Get all orders
- `GET /orders/:id` - Get order by ID
- `POST /orders` - Create new order
- `PUT /orders/:id` - Update order
- `DELETE /orders/:id` - Delete order

## Kafka Topics

- `orders-topic` - Input for order processing
- `processed-orders` - Output with validated orders
### Test producer
```bash
# orders-topic producer
docker exec -it kafka opt/kafka/bin/kafka-console-producer.sh --topic orders-topic --bootstrap-server kafka:9092

# processed-order
docker exec -it kafka opt/kafka/bin/kafka-console-producer.sh --topic orders-topic --bootstrap-server kafka:9092
```
### Testing
Succes response or Approve Status
```json
{
  "user_id":"user123",
  "product":"laptop",
  "quantity":1,
  "total_amount":15000.00
}

{
  "user_id":"user456",
  "product":"smartphone",
  "quantity":2,
  "total_amount":15090.00
}

{
  "user_id":"user456",
  "product":"laptop",
  "quantity":3,
  "total_amount":150980.00
}

{
  "user_id":"user123",
  "product":"smartphone",
  "quantity":4,
  "total_amount":15200.00
}
```
Failed response or Rejected Status
```json
{
  "user_id":"user456",
  "product":"pillow",
  "quantity":3,
  "total_amount":1200.00
}

{
  "user_id":"user921",
  "product":"laptop",
  "quantity":3,
  "total_amount":1200.00
}

{
  "user_id":"user900",
  "product":"laptop","quantity":3,
  "total_amount":1200.00
}

{
  "user_id":"user456",
  "product":"ice cream",
  "quantity":3,
  "total_amount":1200.00
}
```