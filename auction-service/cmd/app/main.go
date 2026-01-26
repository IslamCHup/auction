package main

import (
	"auction-service/internal/config"
	"auction-service/internal/kafka"
	"auction-service/internal/repository"
	"auction-service/internal/services"
	"auction-service/internal/transport"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	db := config.InitDatabase()

	kafkaProducer, err := kafka.NewProducer()
	if err != nil {
		log.Printf("WARNING: Failed to initialize Kafka producer: %v. Continuing without Kafka.", err)
		kafkaProducer = nil
	} else {
		defer kafkaProducer.Close()
	}

	lotRepository := repository.NewLotRepository(db)
	bidRepository := repository.NewBidRepository(db)
	lotService := services.NewLotService(lotRepository, bidRepository, kafkaProducer)
	lotHandler := transport.NewLotHandler(lotService)
	bidService := services.NewBidService(bidRepository, lotRepository, kafkaProducer)
	bidHandler := transport.NewBidHandler(bidService)

	go startAuctionWorker(lotService)

	tutu := server.Group("/api")

	tutu.POST("/lots", lotHandler.CreateLot)
	tutu.GET("/lots", lotHandler.GetAllLots)
	tutu.GET("/lots/:id", lotHandler.GetLotByID)
	tutu.PUT("/lots/:id", lotHandler.UpdateLot)
	tutu.POST("/lots/:id/publish", lotHandler.PublishLot)
	tutu.POST("/lots/complete-expired", lotHandler.CompleteExpired)
	tutu.POST("/lots/:id/force-complete", lotHandler.ForceComplete)
	tutu.POST("/lots/:id/bids", bidHandler.CreateBid)
	tutu.GET("/lots/:id/bids", bidHandler.GetAllBids)
	tutu.GET("/users/:id/lots", lotHandler.GetAllLotsByUser)
	tutu.GET("/users/:id/bids", bidHandler.GetAllBidsByUser)

	server.Run(":8081")
}

func startAuctionWorker(lotService services.LotService) {
	interval := 5 * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := lotService.CompleteExpiredLots(); err != nil {
			continue
		}
	}
}
