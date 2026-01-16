package services

import (
	"fmt"
	"log/slog"
	"notification-service/internal/models"
	"notification-service/internal/repository"
)

type NotificationService interface {
	CreateWinnerLoserNotification(event *models.LotCompletedEvent) error
	CreateBidPlacedNotification(event *models.BidPlacedEvent) error
	ListNotification(filter models.FilterNotification) ([]models.Notification, error)
	MarkAsRead(id uint64) error
	CountUnread(userID uint64) (int64, error)
}

type notificationService struct {
	repo   repository.NotificationRepository
	logger *slog.Logger
}

func NewNotificationService(repo repository.NotificationRepository, logger *slog.Logger) NotificationService {
	return &notificationService{repo: repo, logger: logger}
}

func (s *notificationService) CreateWinnerLoserNotification(event *models.LotCompletedEvent) error {

	winnerNotif := &models.Notification{
		UserID:  event.WinnerID,
		LotID:   event.LotID,
		Type:    models.NotificationTypeAuctionWon,
		Title:   "Аукцион выигран",
		Message: fmt.Sprintf("Поздравляем Вы победили, ваш ставка %d", event.FinalPrice),
	}
	if err := s.repo.CreateNotification(winnerNotif); err != nil {
		s.logger.Error("create notification failed", "err", err.Error(), "user_id", winnerNotif.UserID)
		return err
	}

	for _, losers := range event.LoserIDs {
		loserNotif := &models.Notification{
			UserID:  losers,
			LotID:   event.LotID,
			Type:    models.NotificationTypeAuctionLost,
			Title:   "Аукцион проигран",
			Message: fmt.Sprintf("Вы проиграли аукцион"),
		}
		if err := s.repo.CreateNotification(loserNotif); err != nil {
			s.logger.Error("create notification failed", "err", err.Error(), "user_id", loserNotif.UserID)
			continue
		}
		s.logger.Info("loser notification created", "id", loserNotif.ID, "user_id", loserNotif.UserID)
	}
	s.logger.Info("winner notification created", "id", winnerNotif.ID, "user_id", winnerNotif.UserID)
	return nil
}

func (s *notificationService) CreateBidPlacedNotification(event *models.BidPlacedEvent) error {
	bidPlaced := models.Notification{
		UserID:  event.PreviousLeaderID,
		LotID:   event.LotID,
		Type:    models.NotificationTypeBidOutbid,
		Title:   "Ставка перебита",
		Message: fmt.Sprintf("Ваша ставка на аукционе перебита. Сумма последней ставки %d", event.NewBidAmount),
	}

	if err := s.repo.CreateNotification(&bidPlaced); err != nil {
		s.logger.Error("create bid placed notification failed", "err", err, "user_id", bidPlaced.UserID, "lot_id", bidPlaced.LotID)
		return err
	}
	return nil
}

func (s *notificationService) ListNotification(filter models.FilterNotification) ([]models.Notification, error) {
	list, err := s.repo.ListNotification(filter)
	if err != nil {
		s.logger.Error("list notifications failed", "err", err.Error(), "user_id", filter.UserID)
		return nil, err
	}
	s.logger.Debug("notifications listed", "count", len(list), "limit", filter.Limit, "offset", filter.Offset)
	return list, nil
}

func (s *notificationService) MarkAsRead(id uint64) error {
	if err := s.repo.MarkAsRead(id); err != nil {
		s.logger.Error("mark as read failed", "err", err.Error(), "id", id)
		return err
	}
	s.logger.Info("notification marked as read", "id", id)
	return nil
}

func (s *notificationService) CountUnread(userID uint64) (int64, error) {
	count, err := s.repo.CountUnread(userID)
	if err != nil {
		s.logger.Error("count unread failed", "err", err.Error(), "user_id", userID)
		return 0, err
	}
	s.logger.Debug("unread count", "user_id", userID, "count", count)
	return count, nil
}
