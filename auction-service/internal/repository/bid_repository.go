package repository

import (
	"auction-service/internal/models"

	"gorm.io/gorm"
)

type BidRepository interface {
	CreateBid(bidModel *models.Bid) error
	GetBidByID(id uint64) (*models.Bid, error)
	GetAllBids() ([]models.Bid, error)
	GetAllBidsByUser(userID uint64) ([]models.Bid, error)
	GetAllBidsByLot(lotID uint64) ([]models.Bid, error)
}

type bidRepository struct {
	db *gorm.DB
}

func NewBidRepository(db *gorm.DB) BidRepository {
	return &bidRepository{db: db}
}

func (r *bidRepository) CreateBid(bidModel *models.Bid) error {
	return r.db.Create(bidModel).Error
}

func (r *bidRepository) GetBidByID(id uint64) (*models.Bid, error) {
	var bidModel models.Bid
	if err := r.db.First(&bidModel, uint(id)).Error; err != nil {
		return nil, err
	}
	return &bidModel, nil
}

func (r *bidRepository) GetAllBids() ([]models.Bid, error) {
	var bidModels []models.Bid
	if err := r.db.Find(&bidModels).Error; err != nil {
		return nil, err
	}
	return bidModels, nil
}

func (r *bidRepository) GetAllBidsByUser(userID uint64) ([]models.Bid, error) {
	var bidModels []models.Bid
	if err := r.db.Where("user_id = ?", userID).Find(&bidModels).Error; err != nil {
		return nil, err
	}
	return bidModels, nil
}

func (r *bidRepository) GetAllBidsByLot(lotID uint64) ([]models.Bid, error) {
	var bidModels []models.Bid
	if err := r.db.Where("lot_model_id = ?", lotID).Order("created_at DESC").Find(&bidModels).Error; err != nil {
		return nil, err
	}
	return bidModels, nil
}
