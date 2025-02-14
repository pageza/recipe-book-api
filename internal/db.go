// db.go
package internal

import (
	"fmt"
	"log"
	"os"


	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection
var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "youruser"),
		getEnv("DB_PASSWORD", "yourpassword"),
		getEnv("DB_NAME", "yourdb"),
		getEnv("DB_PORT", "5432"),
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection established")

	// Automatically migrate your models
	err = DB.AutoMigrate(&User{}, &Recipe{})
	if err != nil {
		log.Fatalf("Failed to migrate database schemas: %v", err)
	}
	log.Println("Database migration completed")
}

// getEnv retrieves environment variables or returns a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
