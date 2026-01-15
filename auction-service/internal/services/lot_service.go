package services

import (
	"auction-service/internal/models"
	"auction-service/internal/repository"
	"errors"
	"fmt"
	"time"
)

type LotServiceInterface interface {
	CreateLot(lotModel *models.LotModel) error
	PublishLot(id uint64) error
	GetLotByID(id uint64) (*models.LotModel, error)
	GetAllLots() ([]models.LotModel, error)
}

type LotService struct {
	repository *repository.LotRepository
}

func NewLotService(repository *repository.LotRepository) *LotService {
	return &LotService{repository: repository}
}

func (s *LotService) CreateLot(lotModel *models.LotModel) error {
	lotModel.Status = "draft"
	lotModel.StartDate = time.Now()
	lotModel.EndDate = time.Now().Add(24 * time.Hour)
	return s.repository.CreateLot(lotModel)
}

func (s *LotService) PublishLot(id uint64) error {
	lotModel, err := s.repository.GetLotByID(id)
	if err != nil {
		return fmt.Errorf("failed to get lot: %w", err)
	}
	lotModel.Status = "active"
	if err := s.repository.CreateLot(lotModel); err != nil {
		return fmt.Errorf("failed to publish lot: %w", err)
	}

	return nil
}

func (s *LotService) GetLotByID(id uint64) (*models.LotModel, error) {
	return s.repository.GetLotByID(uint64(id))
}

func (s *LotService) GetAllLots() ([]models.LotModel, error) {
	return s.repository.GetAllLots()
}

func (s *LotService) UpdateLot(lotModel *models.LotModel) error {
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

func (s *LotService) GetAllLotsByUser(userID uint64) ([]models.LotModel, error) {
	return s.repository.GetAllLotsByUser(userID)
}

func (s *LotService) CompleteExpiredLots() error {
	expiredLots, err := s.repository.GetExpiredActiveLots()
	if err != nil {
		return err
	}

	for _, lot := range expiredLots {
		lot.Status = models.LotStatusCompleted
		if err := s.repository.UpdateLot(&lot); err != nil {
			return err
		}
	}
	return nil
}
