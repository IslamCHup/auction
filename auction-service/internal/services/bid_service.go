package services

import (
	"auction-service/internal/models"
	"auction-service/internal/repository"
)

type BidService struct {
	repository *repository.BidRepository
}

func NewBidService(repository *repository.BidRepository) *BidService {
	return &BidService{repository: repository}
}

func (s *BidService) CreateBid(bidModel *models.Bid) error {
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