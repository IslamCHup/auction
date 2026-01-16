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
	walletRepo := repository.NewWalletRepository(db)

	jwt := services.NewJWTService()
	userSvc := services.NewUserService(userRepo, jwt)
	walletSvc := services.NewWalletService(walletRepo, db)

	authHandler := transport.NewAuthHandler(userSvc, jwt)
	walletHandler := transport.NewWalletHandler(userSvc, walletSvc)

	r := transport.SetupRouter(authHandler, jwt, walletHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
