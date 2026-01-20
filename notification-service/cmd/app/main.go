package main

import (
	"log"
	"notification-service/internal/config"
	"notification-service/internal/db"
	"notification-service/internal/repository"
	"notification-service/internal/services"
	"notification-service/internal/transport"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("env file not found, using system env")
	}

	logger := config.InitLogger()

	dbConn := db.InitDatabase()
	if dbConn == nil {
		logger.Error("db is nil")
		os.Exit(1)
	}

	notificationRepo := repository.NewNotificationRepository(dbConn, logger)

	notificationService := services.NewNotificationService(notificationRepo, logger)

	notificationHandler := transport.NewNotificationHandler(notificationService, logger)

	r := gin.Default()

	notificationHandler.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	if err := r.Run(":" + port); err != nil {
		logger.Error("failed to run server", "err", err)
	}

}
