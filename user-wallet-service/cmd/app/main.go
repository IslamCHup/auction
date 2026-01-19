package main

import (
	"os"

	"user-service/internal/config"
	"user-service/internal/repository"
	"user-service/internal/services"
	"user-service/internal/transport"
)

func main() {
	logger := config.InitLogger()

	// load database
	db := config.InitDatabase()

	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)

	jwt := services.NewJWTService()
	userSvc := services.NewUserService(userRepo, jwt)
	walletSvc := services.NewWalletService(walletRepo, db)

	authHandler := transport.NewAuthHandler(userSvc, jwt)
	walletHandler := transport.NewWalletHandler(userSvc, walletSvc)

	r := transport.SetupRouter(logger, authHandler, jwt, walletHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	logger.Info("service started", "port", port)

	if err := r.Run(":" + port); err != nil {
		logger.Error("failed to run service", "err", err.Error())
	}
}
