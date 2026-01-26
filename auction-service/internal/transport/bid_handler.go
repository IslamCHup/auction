package transport

import (
	"auction-service/internal/models"
	"auction-service/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BidHandler struct {
	service services.BidService
}

func NewBidHandler(service services.BidService) *BidHandler {
	return &BidHandler{service: service}
}

func (h *BidHandler) CreateBid(c *gin.Context) {
	lotID := c.Param("id")
	lotIDUint, err := strconv.ParseUint(lotID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lot id"})
		return
	}

	var bidModel models.Bid
	if err := c.ShouldBindJSON(&bidModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bidModel.LotModelID = uint(lotIDUint)
	if bidModel.UserID == 0 {
		if uidStr := c.GetHeader("X-User-Id"); uidStr != "" {
			if uidParsed, convErr := strconv.Atoi(uidStr); convErr == nil && uidParsed > 0 {
				bidModel.UserID = uint(uidParsed)
			}
		}
	}
	bidModel.CreatedAt = time.Now().UTC()

	err = h.service.CreateBid(&bidModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Bid created successfully"})
}

func (h *BidHandler) GetBidByID(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bidModel, err := h.service.GetBidByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bidModel)
}

func (h *BidHandler) GetAllBids(c *gin.Context) {
	lotID := c.Param("id")
	lotIDUint, err := strconv.ParseUint(lotID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lot id"})
		return
	}
	bids, err := h.service.GetAllBidsByLot(lotIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bids)
}

func (h *BidHandler) GetAllBidsByUser(c *gin.Context) {
	var userIDUint uint64
	if uidStr := c.GetHeader("X-User-Id"); uidStr != "" {
		if parsed, err := strconv.ParseUint(uidStr, 10, 64); err == nil && parsed > 0 {
			userIDUint = parsed
		}
	}
	if userIDUint == 0 {
		userID := c.Param("id")
		parsed, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userIDUint = parsed
	}
	if userIDUint == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	bids, err := h.service.GetAllBidsByUser(userIDUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bids)
}
