package transport

import (
	"errors"
	"log/slog"
	"net/http"
	"notification-service/internal/models"
	"notification-service/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NotificationHandler struct {
	service services.NotificationService
	logger  *slog.Logger
}

func NewNotificationHandler(service services.NotificationService, logger *slog.Logger) *NotificationHandler {
	return &NotificationHandler{service: service, logger: logger}
}

func (h *NotificationHandler) RegisterRoutes(r *gin.Engine) {
	notifications := r.Group("/api/notifications")
	{
		notifications.POST("/", h.Create)
		notifications.PATCH("/:id/read", h.MarkAsRead)
		notifications.GET("/unread-count", h.CountUnread)
		notifications.GET("/", h.ListNotification)
	}
}

func (h *NotificationHandler) Create(c *gin.Context) {
	var req models.Notification
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind create", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(&req); err != nil {
		h.logger.Error("create notification", "err", err.Error(), "user_id", req.UserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create notification"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (h *NotificationHandler) ListNotification(c *gin.Context) {
	var filter models.FilterNotification
	if err := c.ShouldBindQuery(&filter); err != nil {
		h.logger.Error("bind query", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if v, ok := c.Get("user_id"); ok {
		switch t := v.(type) {
		case uint64:
			filter.UserID = &t
		case int:
			u := uint64(t)
			filter.UserID = &u
		case string:
			if parsed, err := strconv.ParseUint(t, 10, 64); err == nil {
				filter.UserID = &parsed
			}
		}
	}

	// Фоллбэк: берем user_id из заголовка, который проставляет gateway
	if filter.UserID == nil {
		if uidStr := c.GetHeader("X-User-Id"); uidStr != "" {
			if parsed, err := strconv.ParseUint(uidStr, 10, 64); err == nil {
				filter.UserID = &parsed
			}
		}
	}

	list, err := h.service.ListNotification(filter)
	if err != nil {
		h.logger.Error("list notifications", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list notifications"})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error("parse id", "err", err.Error(), "param", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.MarkAsRead(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
			return
		}
		h.logger.Error("mark as read", "err", err.Error(), "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark as read"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *NotificationHandler) CountUnread(c *gin.Context) {
	v, ok := c.Get("user_id")
	if !ok {
		// fallback: берем из заголовка
		if uidStr := c.GetHeader("X-User-Id"); uidStr != "" {
			if parsed, err := strconv.ParseUint(uidStr, 10, 64); err == nil {
				v = parsed
				ok = true
			}
		}
		if !ok {
			h.logger.Error("user_id missing in context and header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
	}

	userID, ok := v.(uint64)
	if !ok {
		// попытка преобразовать, если это строка/число
		switch t := v.(type) {
		case int:
			userID = uint64(t)
		case string:
			if parsed, err := strconv.ParseUint(t, 10, 64); err == nil {
				userID = parsed
			} else {
				h.logger.Error("invalid user_id type", slog.Any("value", v))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
				return
			}
		default:
			h.logger.Error("invalid user_id type", slog.Any("value", v))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
	}

	count, err := h.service.CountUnread(userID)
	if err != nil {
		h.logger.Error("failed to count unread", slog.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count unread"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
