package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"notification-service/internal/models"
	"notification-service/internal/services"

	"github.com/segmentio/kafka-go"
)

func RunConsumerLotCompleted(
	ctx context.Context,
	logger *slog.Logger,
	service services.NotificationService,
) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    "lot_completed",
		GroupID:  "notifications-group",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	defer reader.Close()

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Info("consumer stopped")
				return
			}

			logger.Error("failed to fetch message", "err", err)
			continue
		}

		var event models.LotCompletedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			logger.Error("invalid message format", "err", err)

			_ = reader.CommitMessages(ctx, msg)
			continue
		}

		err = service.CreateWinnerLoserNotification(&event)
		if err != nil {
			logger.Error("failed to process event", "err", err)
			continue
		}
		if err := reader.CommitMessages(ctx, msg); err != nil {
			logger.Error("failed to commit message", "err", err)
		}
	}
}

func RunConsumerBidPlaced(
	ctx context.Context,
	logger *slog.Logger,
	service services.NotificationService,
) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    "bid_placed",
		GroupID:  "notifications-group",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	defer reader.Close()

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Info("consumer stopped")
				return
			}

			logger.Error("failed to fetch message", "err", err)
			continue
		}

		var event models.BidPlacedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			logger.Error("invalid message format", "err", err)

			_ = reader.CommitMessages(ctx, msg)
			continue
		}

		err = service.CreateBidPlacedNotification(&event)
		if err != nil {
			logger.Error("failed to process event", "err", err)
			continue
		}
		if err := reader.CommitMessages(ctx, msg); err != nil {
			logger.Error("failed to commit message", "err", err)
		}
	}
}

/*
Минимально необходимое в БД:
CREATE UNIQUE INDEX uniq_bid_outbid
ON notifications (user_id, lot_id, type);


И в сервисе:

if errors.Is(err, gorm.ErrDuplicatedKey) {
	return nil // считаем событие успешно обработанным
}


Без этого дубликаты неизбежны, и это не баг consumer’а.
*/
