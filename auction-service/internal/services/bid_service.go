package services

import (
	"auction-service/internal/kafka"
	"auction-service/internal/models"
	"auction-service/internal/repository"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
)

type BidService interface {
	CreateBid(bidModel *models.Bid) error
	GetBidByID(id uint64) (*models.Bid, error)
	GetAllBids() ([]models.Bid, error)
	GetAllBidsByUser(userID uint64) ([]models.Bid, error)
	GetAllBidsByLot(lotID uint64) ([]models.Bid, error)
}

type bidService struct {
	repository    repository.BidRepository
	lotRepository repository.LotRepository
	kafkaProducer *kafka.Producer
}

func NewBidService(repository repository.BidRepository, lotRepository repository.LotRepository, kafkaProducer *kafka.Producer) BidService {
	return &bidService{
		repository:    repository,
		lotRepository: lotRepository,
		kafkaProducer: kafkaProducer,
	}
}

func (s *bidService) freezeWallet(userID uint, amount int64) error {
	base := os.Getenv("WALLET_SERVICE_URL")
	if base == "" {
		return errors.New("wallet service url is not configured")
	}
	url := fmt.Sprintf("%s/api/wallet/freeze", base)

	jsonData, err := json.Marshal(map[string]interface{}{
		"amount":      amount,
		"description": "", // опционально, сервис подставит дефолт
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", userID))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call wallet service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("wallet service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *bidService) unfreezeWallet(userID uint, amount int64) error {
	base := os.Getenv("WALLET_SERVICE_URL")
	if base == "" {
		return errors.New("wallet service url is not configured")
	}
	url := fmt.Sprintf("%s/api/wallet/unfreeze", base)

	jsonData, err := json.Marshal(map[string]interface{}{
		"amount":      amount,
		"description": "",
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", userID))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call wallet service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("wallet service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *bidService) CreateBid(bidModel *models.Bid) error {
	lotModel, err := s.lotRepository.GetLotByID(uint64(bidModel.LotModelID))
	if err != nil {
		return errors.New("failed to get lot")
	}

	if lotModel.Status != models.LotStatusActive {
		return errors.New("bids can only be placed on active lots")
	}

	// Используем серверное время в UTC для валидации и фиксации времени ставки
	now := time.Now().UTC()
	bidModel.CreatedAt = now

	lotStart := lotModel.StartDate.UTC()
	lotEnd := lotModel.EndDate.UTC()

	if now.Before(lotStart) {
		return errors.New("bid cannot be created before the lot start date")
	}
	if now.After(lotEnd) {
		return errors.New("bid cannot be created after the lot end date")
	}

	minRequiredAmount := lotModel.CurrentPrice + lotModel.MinStep
	if bidModel.Amount < minRequiredAmount {
		return fmt.Errorf("bid amount must be at least %d (current price %d + min step %d)",
			minRequiredAmount, lotModel.CurrentPrice, lotModel.MinStep)
	}

	// Сохранить ID предыдущей ставки для разморозки
	previousBidID := lotModel.CurrentBidID
	var previousBid *models.Bid
	if previousBidID != 0 {
		var err error
		previousBid, err = s.repository.GetBidByID(previousBidID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("WARNING: failed to get previous bid %d: %v", previousBidID, err)
				previousBid = nil
			}
		}
	}

	// 1. Заморозить средства нового участника
	err = s.freezeWallet(bidModel.UserID, bidModel.Amount)
	if err != nil {
		return fmt.Errorf("failed to freeze wallet: %w", err)
	}

	// 2. Сохранить ставку
	err = s.repository.CreateBid(bidModel)
	if err != nil {
		// Компенсация: откатить заморозку средств текущего пользователя при ошибке сохранения
		s.unfreezeWallet(bidModel.UserID, bidModel.Amount)
		return fmt.Errorf("failed to create bid: %w", err)
	}

	// Проверить, что GORM установил ID
	if bidModel.ID == 0 {
		// Компенсация: откатить заморозку средств текущего пользователя
		s.unfreezeWallet(bidModel.UserID, bidModel.Amount)
		return errors.New("failed to create bid: ID not set")
	}

	// 3. Обновить лот: текущая цена и ID текущей ставки
	lotModel.CurrentPrice = bidModel.Amount
	lotModel.CurrentBidID = uint64(bidModel.ID)
	err = s.lotRepository.UpdateLot(lotModel)
	if err != nil {
		// Компенсация: откатить заморозку средств текущего пользователя при ошибке обновления
		s.unfreezeWallet(bidModel.UserID, bidModel.Amount)
		return fmt.Errorf("failed to update lot: %w", err)
	}

	if previousBid != nil {
		err = s.unfreezeWallet(previousBid.UserID, previousBid.Amount)
		if err != nil {
			log.Printf("WARNING: failed to unfreeze wallet for previous bid %d: %v", previousBid.ID, err)
		}
	}

	// Отправка события в Kafka о создании ставки
	if s.kafkaProducer != nil {
		previousLeaderID := uint64(0)
		if previousBid != nil {
			previousLeaderID = uint64(previousBid.UserID)
		}

		// Формируем событие в формате, который ожидает notification-service.
		event := kafka.BidPlacedEvent{
			LotID:            uint64(bidModel.LotModelID),
			PreviousLeaderID: previousLeaderID,
			NewBidAmount:     bidModel.Amount,
		}

		if err := s.kafkaProducer.SendMessage("bid_placed", fmt.Sprintf("%d", bidModel.LotModelID), event); err != nil {
			log.Printf("WARNING: failed to send bid_placed event to Kafka: %v", err)
		}
	}

	return nil
}

func (s *bidService) GetBidByID(id uint64) (*models.Bid, error) {
	return s.repository.GetBidByID(uint64(id))
}

func (s *bidService) GetAllBids() ([]models.Bid, error) {
	return s.repository.GetAllBids()
}

func (s *bidService) GetAllBidsByUser(userID uint64) ([]models.Bid, error) {
	return s.repository.GetAllBidsByUser(userID)
}

func (s *bidService) GetAllBidsByLot(lotID uint64) ([]models.Bid, error) {
	return s.repository.GetAllBidsByLot(lotID)
}
