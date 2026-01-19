package repository

import (
	"errors"

	"log/slog"

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
	db     *gorm.DB
	logger *slog.Logger
}

func NewWalletRepository(db *gorm.DB, logger *slog.Logger) WalletRepository {
	return &walletRepository{db: db, logger: logger}
}

func (r *walletRepository) WithDB(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db, logger: r.logger}
}

func (r *walletRepository) GetByUserID(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		if r.logger != nil {
			r.logger.Error("db get wallet failed", "user_id", userID, "err", err.Error())
		}
		return nil, err
	}
	if r.logger != nil {
		r.logger.Info("db got wallet", "user_id", userID, "wallet_id", wallet.ID)
	}
	return &wallet, nil
}

func (r *walletRepository) GetByUserIDWithLock(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		if r.logger != nil {
			r.logger.Error("db get wallet lock failed", "user_id", userID, "err", err.Error())
		}
		return nil, err
	}
	if r.logger != nil {
		r.logger.Info("db got wallet with lock", "user_id", userID, "wallet_id", wallet.ID)
	}
	return &wallet, nil
}

func (r *walletRepository) SaveWallet(wallet *models.Wallet) error {
	if r.logger != nil {
		r.logger.Info("db save wallet", "wallet_id", wallet.ID)
	}
	return r.db.Save(wallet).Error
}

func (r *walletRepository) CreateWallet(wallet *models.Wallet) error {
	if r.logger != nil {
		r.logger.Info("db create wallet", "user_id", wallet.UserID)
	}
	return r.db.Create(wallet).Error
}

func (r *walletRepository) CreateTransaction(tx *models.Transaction) error {
	if r.logger != nil {
		r.logger.Info("db create transaction", "wallet_id", tx.WalletID, "user_id", tx.UserID, "type", tx.Type)
	}
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
		if r.logger != nil {
			r.logger.Error("db list transactions failed", "user_id", userID, "err", err.Error())
		}
		return nil, err
	}
	if r.logger != nil {
		r.logger.Info("db list transactions", "user_id", userID, "count", len(txs))
	}
	return txs, nil
}
