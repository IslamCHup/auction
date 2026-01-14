package main

import (
	"log"
	"os"

	"user-service/internal/config"
	"user-service/internal/repository"
	"user-service/internal/services"
	"user-service/internal/transport"
)

func main() {
	db := config.InitDatabase()

	userRepo := repository.NewUserRepository(db)
	jwt := services.NewJWTService()
	userSvc := services.NewUserService(userRepo, jwt)
	authHandler := transport.NewAuthHandler(userSvc, jwt)

	r := transport.SetupRouter(authHandler, jwt)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
