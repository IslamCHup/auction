package repository

import (
	"errors"

	models "user-service/internal/models"

	"gorm.io/gorm"
)

type WalletRepository interface {
	GetByUserID(userID uint) (*models.Wallet, error)
	CreateOrUpdate(wallet *models.Wallet) error
	CreateTransaction(tx *models.Transaction) error
	ListTransactions(userID uint, limit, offset int) ([]models.Transaction, error)
	Transaction(fn func(tx *gorm.DB) error) error
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *walletRepository) GetByUserID(userID uint) (*models.Wallet, error) {
	var w models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&w).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &w, nil
}

func (r *walletRepository) CreateOrUpdate(wallet *models.Wallet) error {
	// Save handles both create (when primary key is zero) and update
	return r.db.Save(wallet).Error
}

func (r *walletRepository) CreateTransaction(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *walletRepository) ListTransactions(userID uint, limit, offset int) ([]models.Transaction, error) {
	var txs []models.Transaction
	q := r.db.Where("user_id = ?", userID).Order("created_at desc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	if err := q.Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}
