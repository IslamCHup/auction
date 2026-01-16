package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	models "user-service/internal/models"
)

func InitDatabase() *gorm.DB {
	dsl := os.Getenv("DATABASE_URL")
	if dsl == "" {
		dsl = "host=localhost user=adamgowz password=9555 dbname=auction/wallet port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{}); err != nil {
		log.Fatal(err)
	}

	return db
}
