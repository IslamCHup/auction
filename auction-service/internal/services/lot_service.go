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
)

type LotService interface {
	CreateLot(lotModel *models.LotModel) error
	PublishLot(id uint64) error
	GetLotByID(id uint64) (*models.LotModel, error)
	GetAllLots(offset int, limit int, filters *repository.LotFilters) ([]models.LotModel, error)
	UpdateLot(lotModel *models.LotModel) error
	GetAllLotsByUser(userID uint64) ([]models.LotModel, error)
	CompleteExpiredLots() error
	ForceCompleteLot(id uint64) error
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
	nowUTC := time.Now().UTC()
	if lotModel.StartDate.IsZero() {
		lotModel.StartDate = nowUTC
	} else {
		lotModel.StartDate = lotModel.StartDate.UTC()
	}
	if lotModel.EndDate.IsZero() {
		lotModel.EndDate = lotModel.StartDate.Add(24 * time.Hour)
	} else {
		lotModel.EndDate = lotModel.EndDate.UTC()
	}

	if lotModel.EndDate.Before(lotModel.StartDate) {
		return errors.New("end date must be after start date")
	}
	if lotModel.StartDate.Before(nowUTC) {
		return errors.New("start date cannot be in the past")
	}

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
	if filters == nil {
		status := models.LotStatusActive
		filters = &repository.LotFilters{Status: &status}
	} else if filters.Status == nil {
		status := models.LotStatusActive
		filters.Status = &status
	}
	return s.repository.GetAllLots(offset, limit, filters)
}

func (s *lotService) UpdateLot(lotModel *models.LotModel) error {
	existingLot, err := s.repository.GetLotByID(uint64(lotModel.ID))
	if err != nil {
		return fmt.Errorf("failed to get lot: %w", err)
	}

	if existingLot.Status != models.LotStatusDraft {
		return errors.New("only draft lots can be updated")
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
		if lot.CurrentBidID != 0 {
			bid, err := s.bidRepository.GetBidByID(lot.CurrentBidID)
			if err == nil && bid != nil {
				lot.WinnerID = uint64(bid.UserID)
			}
		}
		if err := s.repository.UpdateLot(&lot); err != nil {
			return err
		}

		if lot.WinnerID != 0 && lot.CurrentPrice > 0 {
			if err := s.chargeWallet(uint(lot.WinnerID), lot.CurrentPrice, fmt.Sprintf("Auction payment for lot #%d", lot.ID)); err != nil {
				log.Printf("WARNING: failed to charge winner wallet for lot %d: %v", lot.ID, err)
			}
		}

		if s.kafkaProducer != nil {
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

func (s *lotService) chargeWallet(userID uint, amount int64, description string) error {
	base := os.Getenv("WALLET_SERVICE_URL")
	if base == "" {
		return errors.New("wallet service url is not configured")
	}
	url := fmt.Sprintf("%s/api/wallet/charge", base)

	payload := map[string]any{
		"amount":      amount,
		"description": description,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal charge payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to build charge request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", fmt.Sprintf("%d", userID))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call wallet charge: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("wallet charge returned %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (s *lotService) ForceCompleteLot(id uint64) error {
	lot, err := s.repository.GetLotByID(id)
	if err != nil {
		return err
	}

	lot.Status = models.LotStatusCompleted
	if lot.CurrentBidID != 0 {
		bid, err := s.bidRepository.GetBidByID(lot.CurrentBidID)
		if err == nil && bid != nil {
			lot.WinnerID = uint64(bid.UserID)
		}
	}
	if err := s.repository.UpdateLot(lot); err != nil {
		return err
	}

	if lot.WinnerID != 0 && lot.CurrentPrice > 0 {
		if err := s.chargeWallet(uint(lot.WinnerID), lot.CurrentPrice, fmt.Sprintf("Auction payment for lot #%d", lot.ID)); err != nil {
			log.Printf("WARNING: failed to charge winner wallet for lot %d: %v", lot.ID, err)
		}
	}

	if s.kafkaProducer != nil {
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
	return nil
}
