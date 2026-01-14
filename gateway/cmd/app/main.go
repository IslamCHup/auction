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
		logger.Error("env не найдено", "err", err.Error())
	}

	auctionProxy := proxy.NewReverseProxy(os.Getenv("AUCTION_SERVICE_URL"), logger)
	userProxy := proxy.NewReverseProxy(os.Getenv("USER_SERVICE_URL"), logger)
	walletProxy := proxy.NewReverseProxy(os.Getenv("WALLET_SERVICE_URL"), logger)

	r := gin.Default()

	r.Use(cors.Default())
	r.Use(middleware.RateLimitMiddleware())
	r.Use(middleware.TimeoutMiddleware())

	r.Any("/auction/*path", proxy.MakeProxyHandler(auctionProxy))
	r.Any("/user/*path", proxy.MakeProxyHandler(userProxy))
	r.Any("/wallet/*path", proxy.MakeProxyHandler(walletProxy))

	logger.Info(
		"starting gateway",
		"addr", ":8080",
		"auction", os.Getenv("AUCTION_SERVICE_URL"),
		"user", os.Getenv("USER_SERVICE_URL"),
		"wallet", os.Getenv("WALLET_SERVICE_URL"),
	)

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		logger.Error("failed to run server", "err", err.Error())
	}
	
}
