package services

import (
	"auction-service/internal/kafka"
	"auction-service/internal/models"
	"auction-service/internal/repository"
	"errors"
	"fmt"
	"log"
	"time"
)

type LotService interface {
	CreateLot(lotModel *models.LotModel) error
	PublishLot(id uint64) error
	GetLotByID(id uint64) (*models.LotModel, error)
	GetAllLots(offset int, limit int) ([]models.LotModel, error)
	UpdateLot(lotModel *models.LotModel) error
	GetAllLotsByUser(userID uint64) ([]models.LotModel, error)
	CompleteExpiredLots() error
}

type lotService struct {
	repository    repository.LotRepository
	bidRepository repository.BidRepository
	kafkaProducer *kafka.Producer
}

func NewLotService(repository repository.LotRepository, bidRepository repository.BidRepository, kafkaProducer *kafka.Producer) LotService {
	return &lotService{
		repository:    repository,
		bidRepository: bidRepository,
		kafkaProducer: kafkaProducer,
	}
}

func (s *lotService) CreateLot(lotModel *models.LotModel) error {
	lotModel.Status = models.LotStatusDraft
	lotModel.StartDate = time.Now()
	lotModel.EndDate = time.Now().Add(24 * time.Hour)
	lotModel.CurrentPrice = lotModel.StartPrice
	return s.repository.CreateLot(lotModel)
}

func (s *lotService) PublishLot(id uint64) error {
	lotModel, err := s.repository.GetLotByID(id)
	if err != nil {
		return fmt.Errorf("failed to get lot: %w", err)
	}
	if lotModel.Status != models.LotStatusDraft {
		return errors.New("only draft lots can be published")
	}
	lotModel.Status = models.LotStatusActive
	// Убедимся, что CurrentPrice инициализирован
	if lotModel.CurrentPrice == 0 {
		lotModel.CurrentPrice = lotModel.StartPrice
	}
	if err := s.repository.UpdateLot(lotModel); err != nil {
		return fmt.Errorf("failed to publish lot: %w", err)
	}

	return nil
}

func (s *lotService) GetLotByID(id uint64) (*models.LotModel, error) {
	return s.repository.GetLotByID(uint64(id))
}

func (s *lotService) GetAllLots(offset int, limit int) ([]models.LotModel, error) {
	return s.repository.GetAllLots(offset, limit)
}

func (s *lotService) UpdateLot(lotModel *models.LotModel) error {
	if lotModel.Status != models.LotStatusDraft {
		return errors.New("only draft lots can be updated")
	}
	if lotModel.StartDate.Before(time.Now()) {
		return errors.New("lot start date cannot be in the past")
	}
	if lotModel.EndDate.Before(time.Now()) {
		return errors.New("lot end date cannot be in the past")
	}
	if lotModel.StartPrice <= 0 {
		return errors.New("start price must be greater than 0")
	}
	if lotModel.MinStep <= 0 {
		return errors.New("min step must be greater than 0")
	}

	return s.repository.UpdateLot(lotModel)
}

func (s *lotService) GetAllLotsByUser(userID uint64) ([]models.LotModel, error) {
	return s.repository.GetAllLotsByUser(userID)
}

func (s *lotService) CompleteExpiredLots() error {
	expiredLots, err := s.repository.GetExpiredActiveLots()
	if err != nil {
		return err
	}

	for _, lot := range expiredLots {
		lot.Status = models.LotStatusCompleted
		// Установить WinnerID из текущей лучшей ставки
		if lot.CurrentBidID != 0 {
			bid, err := s.bidRepository.GetBidByID(lot.CurrentBidID)
			if err == nil && bid != nil {
				lot.WinnerID = uint64(bid.UserID)
			}
		}
		if err := s.repository.UpdateLot(&lot); err != nil {
			return err
		}

		// Отправка события в Kafka о завершении лота
		if s.kafkaProducer != nil {
			event := map[string]interface{}{
				"lot_id":      lot.ID,
				"winner_id":   lot.WinnerID,
				"final_price": lot.CurrentPrice,
			}
			if err := s.kafkaProducer.SendMessage("auction.lot.completed", fmt.Sprintf("%d", lot.ID), event); err != nil {
				log.Printf("WARNING: failed to send auction.lot.completed event to Kafka: %v", err)
			}
		}
	}
	return nil
}
