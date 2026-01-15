package services

import (
	"auction-service/internal/models"
	"auction-service/internal/repository"
	"errors"
)

type BidService struct {
	repository    *repository.BidRepository
	lotRepository *repository.LotRepository
}

func NewBidService(repository *repository.BidRepository, lotRepository *repository.LotRepository) *BidService {
	return &BidService{repository: repository, lotRepository: lotRepository}
}

func (s *BidService) CreateBid(bidModel *models.Bid) error {
	lotModel, err := s.lotRepository.GetLotByID(uint64(bidModel.LotModelID))
	if err != nil {
		return errors.New("failed to get lot")
	}

	if bidModel.CreatedAt.Before(lotModel.StartDate) {
		return errors.New("bid cannot be created before the lot start date")
	}
	if bidModel.CreatedAt.After(lotModel.EndDate) {
		return errors.New("bid cannot be created after the lot end date")
	}
	if bidModel.Amount <= lotModel.StartPrice {
		return errors.New("bid amount must be greater than the lot start price")
	}
	if bidModel.Amount%lotModel.MinStep != 0 {
		return errors.New("bid amount must be a multiple of the lot min step")
	}
	return s.repository.CreateBid(bidModel)
}

func (s *BidService) GetBidByID(id uint64) (*models.Bid, error) {
	return s.repository.GetBidByID(uint64(id))
}

func (s *BidService) GetAllBids() ([]models.Bid, error) {
	return s.repository.GetAllBids()
}

func (s *BidService) GetAllBidsByUser(userID uint64) ([]models.Bid, error) {
	return s.repository.GetAllBidsByUser(userID)
}
