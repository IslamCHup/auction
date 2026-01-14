package repository

import (
	"auction-service/internal/models"

	"gorm.io/gorm"
)

type BidRepository struct {
	db *gorm.DB
}

func NewBidRepository(db *gorm.DB) *BidRepository {
	return &BidRepository{db: db}
}

func (r *BidRepository) CreateBid(bidModel *models.Bid) error {
	return r.db.Create(bidModel).Error
}

func (r *BidRepository) GetBidByID(id uint64) (*models.Bid, error) {
	var bidModel models.Bid
	if err := r.db.First(&bidModel, uint(id)).Error; err != nil {
		return nil, err
	}
	return &bidModel, nil
}

func (r *BidRepository) GetAllBids() ([]models.Bid, error) {
	var bidModels []models.Bid
	if err := r.db.Find(&bidModels).Error; err != nil {
		return nil, err
	}
	return bidModels, nil
}

func (r *BidRepository) GetAllBidsByUser(userID uint64) ([]models.Bid, error) {
	var bidModels []models.Bid
	if err := r.db.Where("user_id = ?", userID).Find(&bidModels).Error; err != nil {
		return nil, err
	}
	return bidModels, nil
}
