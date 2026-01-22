package main

import (
	"gateway/internal/config"
	"gateway/internal/middleware"
	"gateway/internal/proxy"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logger := config.InitLogger()

	if err := godotenv.Load(); err != nil {
		logger.Warn("env file not found, using system env")
	}

	authProxy := proxy.NewReverseProxy(os.Getenv("AUTH_SERVICE_URL"), logger)
	auctionProxy := proxy.NewReverseProxy(os.Getenv("LOT_SERVICE_URL"), logger)
	walletProxy := proxy.NewReverseProxy(os.Getenv("WALLET_SERVICE_URL"), logger)
	notificationProxy := proxy.NewReverseProxy(os.Getenv("NOTIFICATION_SERVICE_URL"), logger)

	r := gin.Default()

	r.Use(cors.Default())
	r.Use(middleware.TimeoutMiddleware())

	r.Any("/api/auth/*path", proxy.MakeProxyHandler(authProxy))
	
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.UserRateLimitMiddleware())
	protected.Use(middleware.BidRateLimitMiddleware())
	protected.Any("/api/users/*path", proxy.MakeProxyHandler(authProxy))
	
	protected.Any("/api/lots", proxy.MakeProxyHandler(auctionProxy))
	protected.Any("/api/lots/*path", proxy.MakeProxyHandler(auctionProxy))

	protected.Any("/api/wallet/*path", proxy.MakeProxyHandler(walletProxy))
	protected.Any("/api/notifications/*path", proxy.MakeProxyHandler(notificationProxy))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info(
		"gateway started",
		"port", port,
		"auth", os.Getenv("AUTH_SERVICE_URL"),
		"lot", os.Getenv("LOT_SERVICE_URL"),
		"wallet", os.Getenv("WALLET_SERVICE_URL"),
		"notification", os.Getenv("NOTIFICATION_SERVICE_URL"),
	)

	if err := r.Run(":" + port); err != nil {
		logger.Error("failed to run gateway", "err", err.Error())
	}
}
