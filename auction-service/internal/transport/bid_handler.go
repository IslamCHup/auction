package transport

import (
	"auction-service/internal/models"
	"auction-service/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BidHandler struct {
	service *services.BidService
}

func NewBidHandler(service *services.BidService) *BidHandler {
	return &BidHandler{service: service}
}

func (h *BidHandler) CreateBid(c *gin.Context) {
	var bidModel models.Bid
	if err := c.ShouldBindJSON(&bidModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.service.CreateBid(&bidModel)
	c.JSON(http.StatusCreated, bidModel)
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
	bids, err := h.service.GetAllBids()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bids)
}

func (h *BidHandler) GetAllBidsByUser(c *gin.Context) {
	userID := c.Param("id")
	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bids, err := h.service.GetAllBidsByUser(userIDUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bids)
}