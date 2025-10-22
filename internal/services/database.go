package services

import (
	"log"

	"github.com/rajasunsire/orders/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	dsn := "host=localhost user=user password=password dbname=ordersdb port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	DB.AutoMigrate(
		&models.Order{},
	)
	log.Println("Database connected and migrated")
}
