package repository

import (
	"errors"

	models "user-service/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	GetByUserID(userID uint) (*models.Wallet, error)
	GetByUserIDWithLock(userID uint) (*models.Wallet, error)
	SaveWallet(wallet *models.Wallet) error
	CreateWallet(wallet *models.Wallet) error
	CreateTransaction(tx *models.Transaction) error
	ListTransactions(userID uint, limit, offset int) ([]models.Transaction, error)
	WithDB(db *gorm.DB) WalletRepository
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) WithDB(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetByUserID(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) GetByUserIDWithLock(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) SaveWallet(wallet *models.Wallet) error {
	return r.db.Save(wallet).Error
}

func (r *walletRepository) CreateWallet(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *walletRepository) CreateTransaction(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *walletRepository) ListTransactions(userID uint, limit, offset int) ([]models.Transaction, error) {
	var txs []models.Transaction
	query := r.db.Where("user_id = ?", userID).Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}
