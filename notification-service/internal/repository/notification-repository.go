package repository

import (
	"log/slog"
	"notification-service/internal/models"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	CreateNotification(req *models.Notification) error
	ListNotification(filter models.FilterNotification) ([]models.Notification, error)
	MarkAsRead(id uint64) error
	CountUnread(userID uint64) (int64, error)
}

type notificationRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewNotificationRepository(db *gorm.DB, logger *slog.Logger) NotificationRepository {
	return &notificationRepository{db: db, logger: logger}
}

func (r *notificationRepository) CreateNotification(req *models.Notification) error {
	if err := r.db.Create(req).Error; err != nil {
		r.logger.Error("failed to create notification", "err", err.Error(), "user_id", req.UserID)
		return err
	}
	r.logger.Debug("notification created", "id", req.ID, "user_id", req.UserID)
	return nil
}

func (r *notificationRepository) ListNotification(filter models.FilterNotification) ([]models.Notification, error) {
	var notifications []models.Notification

	query := r.db.Model(&models.Notification{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	if filter.IsRead != nil {
		query = query.Where("is_read = ?", *filter.IsRead)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}

	order := filter.Order
	if order == "" {
		order = "created_at desc"
	}

	if err := query.Order(order).Limit(limit).Offset(filter.Offset).Find(&notifications).Error; err != nil {
		r.logger.Error("failed to list notifications", "err", err.Error(), "limit", limit, "offset", filter.Offset)
		return nil, err
	}

	r.logger.Debug("notifications listed", "count", len(notifications), "limit", limit, "offset", filter.Offset)
	return notifications, nil
}

func (r *notificationRepository) MarkAsRead(id uint64) error {
	result := r.db.Model(&models.Notification{}).Where("id = ?", id).Update("is_read", true)
	if result.Error != nil {
		r.logger.Error("failed to mark notification as read", "err", result.Error.Error(), "id", id)
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.logger.Debug("notification not found to mark as read", "id", id)
		return gorm.ErrRecordNotFound
	}
	r.logger.Debug("notification marked as read", "id", id, "rows", result.RowsAffected)
	return nil
}

func (r *notificationRepository) CountUnread(userID uint64) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error; err != nil {
		r.logger.Error("failed to count unread notifications", "err", err.Error(), "user_id", userID)
		return 0, err
	}
	r.logger.Debug("unread notifications counted", "user_id", userID, "count", count)
	return count, nil
}
