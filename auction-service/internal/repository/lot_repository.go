package repository

import (
	"auction-service/internal/models"
	"time"

	"gorm.io/gorm"
)

type LotFilters struct {
	Status     *models.LotStatus
	MinPrice   *int64
	MaxPrice   *int64
	MinEndDate *time.Time
	MaxEndDate *time.Time
}

type LotRepository interface {
	CreateLot(lotModel *models.LotModel) error
	UpdateLot(lotModel *models.LotModel) error
	GetLotByID(id uint64) (*models.LotModel, error)
	GetAllLots(offset int, limit int, filters *LotFilters) ([]models.LotModel, error)
	GetAllLotsByUser(userID uint64) ([]models.LotModel, error)
	GetExpiredActiveLots() ([]models.LotModel, error)
}

type lotRepository struct {
	db *gorm.DB
}

func NewLotRepository(db *gorm.DB) LotRepository {
	return &lotRepository{db: db}
}

func (r *lotRepository) CreateLot(lotModel *models.LotModel) error {
	return r.db.Create(lotModel).Error
}

func (r *lotRepository) UpdateLot(lotModel *models.LotModel) error {
	return r.db.Save(lotModel).Error
}

func (r *lotRepository) GetLotByID(id uint64) (*models.LotModel, error) {
	var lotModel models.LotModel
	if err := r.db.Preload("Bids").First(&lotModel, uint(id)).Error; err != nil {
		return nil, err
	}
	if lotModel.Bids == nil {
		lotModel.Bids = []models.Bid{}
	}

	return &lotModel, nil
}

func (r *lotRepository) GetAllLots(offset int, limit int, filters *LotFilters) ([]models.LotModel, error) {
	var lots []models.LotModel
	query := r.db

	if filters != nil {
		if filters.Status != nil {
			query = query.Where("status = ?", *filters.Status)
		}
		if filters.MinPrice != nil {
			query = query.Where("current_price >= ?", *filters.MinPrice)
		}
		if filters.MaxPrice != nil {
			query = query.Where("current_price <= ?", *filters.MaxPrice)
		}
		if filters.MinEndDate != nil {
			query = query.Where("end_date >= ?", *filters.MinEndDate)
		}
		if filters.MaxEndDate != nil {
			query = query.Where("end_date <= ?", *filters.MaxEndDate)
		}
	}

	if err := query.Preload("Bids").Offset(offset).Limit(limit).Order("created_at DESC").Find(&lots).Error; err != nil {
		return nil, err
	}
	for i := range lots {
		if lots[i].Bids == nil {
			lots[i].Bids = []models.Bid{}
		}
	}
	return lots, nil
}

func (r *lotRepository) GetAllLotsByUser(userID uint64) ([]models.LotModel, error) {
	var lots []models.LotModel
	if err := r.db.Preload("Bids").Where("seller_id = ?", userID).Find(&lots).Error; err != nil {
		return nil, err
	}
	for i := range lots {
		if lots[i].Bids == nil {
			lots[i].Bids = []models.Bid{}
		}
	}
	return lots, nil
}

func (r *lotRepository) GetExpiredActiveLots() ([]models.LotModel, error) {
	var lots []models.LotModel
	if err := r.db.Where("status = ? AND end_date < ?", models.LotStatusActive, time.Now()).Find(&lots).Error; err != nil {
		return nil, err
	}
	return lots, nil
}
