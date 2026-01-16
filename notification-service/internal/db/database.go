package db

import (
	"log"
	"notification-service/internal/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	dsl := os.Getenv("DATABASE_URL")
	if dsl == "" {
		dsl = "host=localhost user=postgres password=12345 dbname=postgres port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Notification{})

	return db
}
