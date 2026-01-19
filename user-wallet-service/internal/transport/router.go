package transport

import (
	"log/slog"
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	logger *slog.Logger,
	authHandler *AuthHandler, jwt services.JWTService, walletHandler *WalletHandler,
) *gin.Engine {
	r := gin.New()


	r.Use(gin.Recovery())

	r.Use(LoggingMiddleware(logger))

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		users := api.Group("/users")
		users.Use(AuthMiddleware(jwt))
		{
			users.GET("/me", authHandler.Me)
			users.PUT("/me", authHandler.UpdateMe)
		}

		wallet := api.Group("/wallet")
		wallet.Use(AuthMiddleware(jwt))
		{
			wallet.GET("/", walletHandler.GetWallet)
			wallet.POST("/deposit", walletHandler.WalletDeposit)
			wallet.POST("/freeze", walletHandler.WalletFreeze)
			wallet.POST("/unfreeze", walletHandler.WalletUnfreeze)
			wallet.POST("/charge", walletHandler.WalletCharge)
			wallet.GET("/transactions", walletHandler.ListTransactions)
		}
	}

	return r
}
