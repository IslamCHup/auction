package services

import (
	"errors"

	models "user-service/internal/models"
	"user-service/internal/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletService interface {
	GetWallet(userID uint) (*models.Wallet, error)
	Deposit(userID uint, amount int64, description string) (*models.Wallet, error)
	Freeze(userID uint, amount int64, description string) (*models.Wallet, error)
	Unfreeze(userID uint, amount int64, description string) (*models.Wallet, error)
	Charge(userID uint, amount int64, description string) (*models.Wallet, error)
	ListTransactions(userID uint, limit, offset int) ([]models.Transaction, error)
}

type walletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) GetWallet(userID uint) (*models.Wallet, error) {
	// убедимся, что кошелек существует
	wallet, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		wallet = &models.Wallet{UserID: userID, Balance: 0, FrozenBalance: 0}
		if err := s.repo.CreateOrUpdate(wallet); err != nil {
			return nil, err
		}
	}
	return wallet, nil
}

func (s *walletService) Deposit(userID uint, amount int64, description string) (*models.Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	var result *models.Wallet
	err := s.repo.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				wallet = models.Wallet{UserID: userID, Balance: 0, FrozenBalance: 0}
				if err := tx.Create(&wallet).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		beforeBalance := wallet.Balance
		beforeFrozen := wallet.FrozenBalance
		wallet.Balance += amount
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}
		transaction := models.Transaction{
			WalletID:      wallet.ID,
			UserID:        userID,
			Type:          models.TransactionDeposit,
			Amount:        amount,
			BalanceBefore: beforeBalance,
			BalanceAfter:  wallet.Balance,
			FrozenBefore:  beforeFrozen,
			FrozenAfter:   wallet.FrozenBalance,
			Description:   description,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		result = &wallet
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *walletService) Freeze(userID uint, amount int64, description string) (*models.Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	var res *models.Wallet
	err := s.repo.Transaction(func(txDB *gorm.DB) error {
		var w models.Wallet
		if err := txDB.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&w).Error; err != nil {
			return err
		}
		available := w.Balance - w.FrozenBalance
		if available < amount {
			return errors.New("insufficient available balance")
		}
		beforeBal := w.Balance
		beforeFrozen := w.FrozenBalance
		w.FrozenBalance += amount
		if err := txDB.Save(&w).Error; err != nil {
			return err
		}
		t := models.Transaction{
			WalletID:      w.ID,
			UserID:        userID,
			Type:          models.TransactionFreeze,
			Amount:        amount,
			BalanceBefore: beforeBal,
			BalanceAfter:  w.Balance,
			FrozenBefore:  beforeFrozen,
			FrozenAfter:   w.FrozenBalance,
			Description:   description,
		}
		if err := txDB.Create(&t).Error; err != nil {
			return err
		}
		res = &w
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *walletService) Unfreeze(userID uint, amount int64, description string) (*models.Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	var res *models.Wallet
	err := s.repo.Transaction(func(txDB *gorm.DB) error {
		var w models.Wallet
		if err := txDB.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&w).Error; err != nil {
			return err
		}
		if w.FrozenBalance < amount {
			return errors.New("insufficient frozen balance")
		}
		beforeBal := w.Balance
		beforeFrozen := w.FrozenBalance
		w.FrozenBalance -= amount
		if err := txDB.Save(&w).Error; err != nil {
			return err
		}
		t := models.Transaction{
			WalletID:      w.ID,
			UserID:        userID,
			Type:          models.TransactionUnfreeze,
			Amount:        amount,
			BalanceBefore: beforeBal,
			BalanceAfter:  w.Balance,
			FrozenBefore:  beforeFrozen,
			FrozenAfter:   w.FrozenBalance,
			Description:   description,
		}
		if err := txDB.Create(&t).Error; err != nil {
			return err
		}
		res = &w
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *walletService) Charge(userID uint, amount int64, description string) (*models.Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	var res *models.Wallet
	err := s.repo.Transaction(func(txDB *gorm.DB) error {
		var w models.Wallet
		if err := txDB.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&w).Error; err != nil {
			return err
		}
		if w.FrozenBalance < amount {
			return errors.New("insufficient frozen balance")
		}
		beforeBal := w.Balance
		beforeFrozen := w.FrozenBalance
		w.FrozenBalance -= amount
		w.Balance -= amount
		if w.Balance < 0 {
			return errors.New("resulting balance negative")
		}
		if err := txDB.Save(&w).Error; err != nil {
			return err
		}
		t := models.Transaction{
			WalletID:      w.ID,
			UserID:        userID,
			Type:          models.TransactionCharge,
			Amount:        amount,
			BalanceBefore: beforeBal,
			BalanceAfter:  w.Balance,
			FrozenBefore:  beforeFrozen,
			FrozenAfter:   w.FrozenBalance,
			Description:   description,
		}
		if err := txDB.Create(&t).Error; err != nil {
			return err
		}
		res = &w
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *walletService) ListTransactions(userID uint, limit, offset int) ([]models.Transaction, error) {
	return s.repo.ListTransactions(userID, limit, offset)
}
