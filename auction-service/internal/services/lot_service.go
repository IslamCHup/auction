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
	GetAllLots(offset int, limit int, filters *repository.LotFilters) ([]models.LotModel, error)
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
	// Валидация базовых полей (проверка на обязательность уже в JSON тегах)
	if lotModel.StartPrice <= 0 {
		return errors.New("start price must be greater than 0")
	}
	if lotModel.MinStep <= 0 {
		return errors.New("min step must be greater than 0")
	}

	// Установка значений по умолчанию
	lotModel.Status = models.LotStatusDraft
	if lotModel.StartDate.IsZero() {
		lotModel.StartDate = time.Now()
	}
	if lotModel.EndDate.IsZero() {
		lotModel.EndDate = lotModel.StartDate.Add(24 * time.Hour)
	}

	// Валидация бизнес-правил после установки значений по умолчанию
	if lotModel.EndDate.Before(lotModel.StartDate) || lotModel.EndDate.Equal(lotModel.StartDate) {
		return errors.New("end date must be after start date")
	}
	if lotModel.StartDate.Before(time.Now()) {
		return errors.New("start date cannot be in the past")
	}

	// Инициализация полей
	lotModel.CurrentPrice = lotModel.StartPrice
	lotModel.WinnerID = 0
	lotModel.CurrentBidID = 0

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

func (s *lotService) GetAllLots(offset int, limit int, filters *repository.LotFilters) ([]models.LotModel, error) {
	// Бизнес-правило: по умолчанию показываем только активные лоты
	// Это валидация/правило на уровне service, а не repository
	if filters == nil {
		filters = &repository.LotFilters{
			Status: func() *models.LotStatus {
				status := models.LotStatusActive
				return &status
			}(),
		}
	} else if filters.Status == nil {
		// Если фильтры заданы, но статус не указан, по умолчанию показываем активные
		status := models.LotStatusActive
		filters.Status = &status
	}
	return s.repository.GetAllLots(offset, limit, filters)
}

func (s *lotService) UpdateLot(lotModel *models.LotModel) error {
	// Получить существующий лот для проверки статуса
	existingLot, err := s.repository.GetLotByID(uint64(lotModel.ID))
	if err != nil {
		return fmt.Errorf("failed to get lot: %w", err)
	}

	// Бизнес-правило: редактировать можно только draft лоты
	if existingLot.Status != models.LotStatusDraft {
		return errors.New("only draft lots can be updated")
	}

	// Валидация бизнес-правил
	if lotModel.StartPrice <= 0 {
		return errors.New("start price must be greater than 0")
	}
	if lotModel.MinStep <= 0 {
		return errors.New("min step must be greater than 0")
	}
	if lotModel.EndDate.Before(lotModel.StartDate) {
		return errors.New("end date must be after start date")
	}
	if lotModel.StartDate.Before(time.Now()) {
		return errors.New("lot start date cannot be in the past")
	}
	if lotModel.EndDate.Before(time.Now()) {
		return errors.New("lot end date cannot be in the past")
	}

	// Обновить текущую цену если стартовая цена изменилась
	if lotModel.CurrentPrice == 0 || lotModel.CurrentPrice < lotModel.StartPrice {
		lotModel.CurrentPrice = lotModel.StartPrice
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
			// Пока у нас нет списка всех участников аукциона, LoserIDs оставляем пустым.
			event := kafka.LotCompletedEvent{
				LotID:      uint64(lot.ID),
				Winner:     lot.WinnerID,
				FinalPrice: lot.CurrentPrice,
				LoserIDs:   nil,
			}

			if err := s.kafkaProducer.SendMessage("lot_completed", fmt.Sprintf("%d", lot.ID), event); err != nil {
				log.Printf("WARNING: failed to send lot_completed event to Kafka: %v", err)
			}
		}
	}
	return nil
}
