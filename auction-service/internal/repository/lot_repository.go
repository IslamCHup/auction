package repository

import (
	"auction-service/internal/models"

	"gorm.io/gorm"
)

type LotRepository struct {
	db *gorm.DB
}

func NewLotRepository(db *gorm.DB) *LotRepository {
	return &LotRepository{db: db}
}

func (r *LotRepository) CreateLot(lotModel *models.LotModel) error {
	return r.db.Create(lotModel).Error
}

func (r *LotRepository) GetLotByID(id uint64) (*models.LotModel, error) {
	var lotModel models.LotModel
	if err := r.db.First(&lotModel, uint(id)).Error; err != nil {
		return nil, err
	}

	return &lotModel, nil
}

func (r *LotRepository) GetAllLots() ([]models.LotModel, error) {
	var lots []models.LotModel
	if err := r.db.Find(&lots).Error; err != nil {
		return nil, err
	}
	return lots, nil
}

func (r *LotRepository) GetAllLotsByUser(userID uint64) ([]models.LotModel, error) {
	var lots []models.LotModel
	if err := r.db.Where("user_id = ?", userID).Find(&lots).Error; err != nil {
		return nil, err
	}
	return lots, nil
}