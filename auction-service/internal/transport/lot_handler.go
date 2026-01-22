package transport

import (
	"auction-service/internal/models"
	"auction-service/internal/repository"
	"auction-service/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LotHandler struct {
	service services.LotService
}

func NewLotHandler(service services.LotService) *LotHandler {
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

func paginationParams(c *gin.Context) (page int, limit int) {
	page = 1
	limit = 10

	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
		if limit > 100 {
			limit = 100
		}
	}

	return page, limit
}

func parseFilters(c *gin.Context) *repository.LotFilters {
	filters := &repository.LotFilters{}

	// Фильтр по статусу
	if statusStr := c.Query("status"); statusStr != "" {
		status := models.LotStatus(statusStr)
		if status == models.LotStatusDraft || status == models.LotStatusActive || status == models.LotStatusCompleted {
			filters.Status = &status
		}
	}

	// Фильтр по минимальной цене
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseInt(minPriceStr, 10, 64); err == nil {
			if minPrice > 0 {
				filters.MinPrice = &minPrice
			}
			// minPrice <= 0 трактуем как "фильтр не задан", поэтому просто не устанавливаем его
		}
	}

	// Фильтр по максимальной цене
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseInt(maxPriceStr, 10, 64); err == nil {
			if maxPrice > 0 {
				filters.MaxPrice = &maxPrice
			}
			// maxPrice <= 0 трактуем как "фильтр не задан"
		}
	}

	// Фильтр по минимальной дате окончания
	if minEndDateStr := c.Query("min_end_date"); minEndDateStr != "" {
		if minEndDate, err := time.Parse(time.RFC3339, minEndDateStr); err == nil {
			filters.MinEndDate = &minEndDate
		}
	}

	// Фильтр по максимальной дате окончания
	if maxEndDateStr := c.Query("max_end_date"); maxEndDateStr != "" {
		if maxEndDate, err := time.Parse(time.RFC3339, maxEndDateStr); err == nil {
			filters.MaxEndDate = &maxEndDate
		}
	}

	// Если фильтры не заданы, возвращаем nil (будет использован фильтр по умолчанию в service - только active)
	if filters.Status == nil && filters.MinPrice == nil && filters.MaxPrice == nil && filters.MinEndDate == nil && filters.MaxEndDate == nil {
		return nil
	}

	return filters
}

func (h *LotHandler) GetAllLots(c *gin.Context) {
	page, limit := paginationParams(c)
	offset := (page - 1) * limit
	filters := parseFilters(c)
	lots, err := h.service.GetAllLots(offset, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lots)
}

func (h *LotHandler) UpdateLot(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lot id"})
		return
	}

	var updateReq models.UpdateLotRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получить существующий лот
	lot, err := h.service.GetLotByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lot not found"})
		return
	}

	// Обновить только предоставленные поля
	if updateReq.Title != nil {
		lot.Title = *updateReq.Title
	}
	if updateReq.Description != nil {
		lot.Description = *updateReq.Description
	}
	if updateReq.StartPrice != nil {
		lot.StartPrice = *updateReq.StartPrice
	}
	if updateReq.MinStep != nil {
		lot.MinStep = *updateReq.MinStep
	}
	if updateReq.EndDate != nil {
		lot.EndDate = *updateReq.EndDate
	}

	if err := h.service.UpdateLot(lot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lot updated successfully"})
}

func (h *LotHandler) PublishLot(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.PublishLot(idUint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
