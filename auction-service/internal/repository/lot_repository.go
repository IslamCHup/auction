package repository

import (
	"auction-service/internal/models"
	"time"

	"gorm.io/gorm"
)

type LotRepository interface {
	CreateLot(lotModel *models.LotModel) error
	UpdateLot(lotModel *models.LotModel) error
	GetLotByID(id uint64) (*models.LotModel, error)
	GetAllLots(offset int, limit int) ([]models.LotModel, error)
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
	if err := r.db.First(&lotModel, uint(id)).Error; err != nil {
		return nil, err
	}

	return &lotModel, nil
}

func (r *lotRepository) GetAllLots(offset int, limit int) ([]models.LotModel, error) {
	var lots []models.LotModel
	if err := r.db.Offset(offset).Limit(limit).Find(&lots).Error; err != nil {
		return nil, err
	}
	return lots, nil
}

func (r *lotRepository) GetAllLotsByUser(userID uint64) ([]models.LotModel, error) {
	var lots []models.LotModel
	if err := r.db.Where("user_id = ?", userID).Find(&lots).Error; err != nil {
		return nil, err
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
