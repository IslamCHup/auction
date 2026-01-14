package services

import (
	models "user-service/internal/models"
	"user-service/internal/repository"
	"user-service/internal/utils"

	"gorm.io/gorm"
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
	db   *gorm.DB
}

func NewWalletService(repo repository.WalletRepository, db *gorm.DB) WalletService {
	return &walletService{repo: repo, db: db}
}

func (s *walletService) GetWallet(userID uint) (*models.Wallet, error) {
	// убедимся, что кошелек существует
	wallet, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		wallet = &models.Wallet{UserID: userID, Balance: 0, FrozenBalance: 0}
		if err := s.repo.CreateWallet(wallet); err != nil {
			return nil, err
		}
	}
	return wallet, nil
}

func (s *walletService) Deposit(userID uint, amount int64, description string) (*models.Wallet, error) {

	if amount <= 0 {
		return nil, utils.ErrAmountMustBePositive
	}

	var result *models.Wallet
	err := s.db.Transaction(func(tx *gorm.DB) error {
		walletRepo := s.repo.WithDB(tx)

		wallet, err := walletRepo.GetByUserIDWithLock(userID)
		if err != nil {
			return err
		}
		if wallet == nil {
			wallet = &models.Wallet{UserID: userID, Balance: 0, FrozenBalance: 0}
			if err := walletRepo.CreateWallet(wallet); err != nil {
				return err
			}
		}
		beforeBalance := wallet.Balance
		beforeFrozen := wallet.FrozenBalance
		wallet.Balance += amount
		if err := walletRepo.SaveWallet(wallet); err != nil {
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
		if err := walletRepo.CreateTransaction(&transaction); err != nil {
			return err
		}
		result = wallet
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *walletService) Freeze(userID uint, amount int64, description string) (*models.Wallet, error) {
	if amount <= 0 {
		return nil, utils.ErrAmountMustBePositive
	}

	var result *models.Wallet

	err := s.db.Transaction(func(tx *gorm.DB) error {
		walletRepo := s.repo.WithDB(tx)

		wallet, err := walletRepo.GetByUserIDWithLock(userID)
		if err != nil {
			return err
		}
		if wallet == nil {
			return gorm.ErrRecordNotFound
		}

		available := wallet.Balance - wallet.FrozenBalance
		if available < amount {
			return utils.ErrInsufficientAvailableBalance
		}

		beforeBalance := wallet.Balance
		beforeFrozen := wallet.FrozenBalance

		wallet.FrozenBalance += amount
		if err := walletRepo.SaveWallet(wallet); err != nil {
			return err
		}

		tran := models.Transaction{
			WalletID:      wallet.ID,
			UserID:        userID,
			Type:          models.TransactionFreeze,
			Amount:        amount,
			BalanceBefore: beforeBalance,
			BalanceAfter:  wallet.Balance,
			FrozenBefore:  beforeFrozen,
			FrozenAfter:   wallet.FrozenBalance,
			Description:   description,
		}

		if err := walletRepo.CreateTransaction(&tran); err != nil {
			return err
		}

		result = wallet
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *walletService) Unfreeze(userID uint, amount int64, description string) (*models.Wallet, error) {
	if amount <= 0 {
		return nil, utils.ErrAmountMustBePositive
	}
	var result *models.Wallet
	err := s.db.Transaction(func(txDB *gorm.DB) error {
		walletRepo := s.repo.WithDB(txDB)

		wallet, err := walletRepo.GetByUserIDWithLock(userID)
		if err != nil {
			return err
		}
		if wallet == nil {
			return gorm.ErrRecordNotFound
		}
		if wallet.FrozenBalance < amount {
			return utils.ErrInsufficientFrozenBalance
		}
		beforeBalance := wallet.Balance
		beforeFrozen := wallet.FrozenBalance
		wallet.FrozenBalance -= amount
		if err := walletRepo.SaveWallet(wallet); err != nil {
			return err
		}
		t := models.Transaction{
			WalletID:      wallet.ID,
			UserID:        userID,
			Type:          models.TransactionUnfreeze,
			Amount:        amount,
			BalanceBefore: beforeBalance,
			BalanceAfter:  wallet.Balance,
			FrozenBefore:  beforeFrozen,
			FrozenAfter:   wallet.FrozenBalance,
			Description:   description,
		}
		if err := walletRepo.CreateTransaction(&t); err != nil {
			return err
		}
		result = wallet
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *walletService) Charge(userID uint, amount int64, description string) (*models.Wallet, error) {
	if amount <= 0 {
		return nil, utils.ErrAmountMustBePositive
	}
	var result *models.Wallet
	err := s.db.Transaction(func(txDB *gorm.DB) error {
		walletRepo := s.repo.WithDB(txDB)

		wallet, err := walletRepo.GetByUserIDWithLock(userID)
		if err != nil {
			return err
		}
		if wallet == nil {
			return gorm.ErrRecordNotFound
		}
		if wallet.FrozenBalance < amount {
			return utils.ErrInsufficientFrozenBalance
		}

		beforeBalance := wallet.Balance
		beforeFrozen := wallet.FrozenBalance
		wallet.FrozenBalance -= amount
		wallet.Balance -= amount

		if wallet.Balance < 0 {
			return utils.ErrResultingBalanceNegative
		}

		if err := walletRepo.SaveWallet(wallet); err != nil {
			return err
		}

		t := models.Transaction{
			WalletID:      wallet.ID,
			UserID:        userID,
			Type:          models.TransactionCharge,
			Amount:        amount,
			BalanceBefore: beforeBalance,
			BalanceAfter:  wallet.Balance,
			FrozenBefore:  beforeFrozen,
			FrozenAfter:   wallet.FrozenBalance,
			Description:   description,
		}
		if err := walletRepo.CreateTransaction(&t); err != nil {
			return err
		}
		result = wallet
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *walletService) ListTransactions(userID uint, limit, offset int) ([]models.Transaction, error) {
	return s.repo.ListTransactions(userID, limit, offset)
}
