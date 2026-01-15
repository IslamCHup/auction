package main

import (
	"auction-service/internal/config"
	"auction-service/internal/repository"
	"auction-service/internal/services"
	"auction-service/internal/transport"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	db := config.InitDatabase()

	lotRepository := repository.NewLotRepository(db)
	bidRepository := repository.NewBidRepository(db)
	lotService := services.NewLotService(lotRepository)
	lotHandler := transport.NewLotHandler(lotService)
	bidService := services.NewBidService(bidRepository, lotRepository)
	bidHandler := transport.NewBidHandler(bidService)

	go startAuctionWorker(lotService)

	server.POST("/lots", lotHandler.CreateLot)
	server.GET("/lots", lotHandler.GetAllLots)
	server.GET("/lots/:id", lotHandler.GetLotByID)
	server.PUT("/lots/:id", lotHandler.UpdateLot)
	server.POST("/lots/:id/publish", lotHandler.PublishLot)
	server.POST("/lots/:id/bids", bidHandler.CreateBid)
	server.GET("/lots/:id/bids", bidHandler.GetAllBids)
	server.GET("/users/:id/lots", lotHandler.GetAllLotsByUser)
	server.GET("/users/:id/bids", bidHandler.GetAllBidsByUser)

	server.Run(":8081")
}

func startAuctionWorker(lotService *services.LotService) {
	interval := 5 * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := lotService.CompleteExpiredLots(); err != nil {
			continue
		}
	}
}
