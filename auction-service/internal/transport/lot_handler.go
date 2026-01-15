package transport

import (
	"auction-service/internal/models"
	"auction-service/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LotHandler struct {
	service *services.LotService
}

func NewLotHandler(service *services.LotService) *LotHandler {
	return &LotHandler{service: service}
}

func (h *LotHandler) CreateLot(c *gin.Context) {
	var lotModel models.LotModel
	if err := c.ShouldBindJSON(&lotModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateLot(&lotModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lot created successfully"})
}

func (h *LotHandler) GetAllLots(c *gin.Context) {
	lots, err := h.service.GetAllLots()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lots)
}

func (h *LotHandler) UpdateLot(c *gin.Context) {
	var lotModel models.LotModel
	if err := c.ShouldBindJSON(&lotModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateLot(&lotModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lot created successfully"})
}

func (h *LotHandler) PublishLot(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.service.PublishLot(idUint)
	c.JSON(http.StatusOK, gin.H{"message": "Lot published successfully"})
}

func (h *LotHandler) GetLotByID(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lotModel, err := h.service.GetLotByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lotModel)
}

func (h *LotHandler) GetAllLotsByUser(c *gin.Context) {
	userID := c.Param("id")
	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lots, err := h.service.GetAllLotsByUser(userIDUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lots)
}
