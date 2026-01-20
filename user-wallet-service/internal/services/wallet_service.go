package services

import (
	"log/slog"

	models "user-service/internal/models"
	"user-service/internal/repository"
	"user-service/internal/utils"

	"gorm.io/gorm"
)

type WalletService interface {
	GetWallet(userID uint) (*models.Wallet, error)
	Deposit(userID uint, amount int64, description string) (*models.Wallet, *models.Transaction, error)
	Freeze(userID uint, amount int64, description string) (*models.Wallet, *models.Transaction, error)
	Unfreeze(userID uint, amount int64, description string) (*models.Wallet, *models.Transaction, error)
	Charge(userID uint, amount int64, description string) (*models.Wallet, *models.Transaction, error)
	ListTransactions(userID uint, limit, offset int) ([]models.Transaction, error)
}

type walletService struct {
	repo   repository.WalletRepository
	db     *gorm.DB
	logger *slog.Logger
}

func NewWalletService(repo repository.WalletRepository, db *gorm.DB, logger *slog.Logger) WalletService {
	return &walletService{repo: repo, db: db, logger: logger}
}

func (s *walletService) GetWallet(userID uint) (*models.Wallet, error) {
	s.logger.Info("service get wallet", "user_id", userID)
	// убедимся, что кошелек существует
	wallet, err := s.repo.GetByUserID(userID)
	if err != nil {
		s.logger.Error("service get wallet failed", "user_id", userID, "err", err.Error())
		return nil, err
	}
	if wallet == nil {
		wallet = &models.Wallet{UserID: userID, Balance: 0, FrozenBalance: 0}
		if err := s.repo.CreateWallet(wallet); err != nil {
			s.logger.Error("service create wallet failed", "user_id", userID, "err", err.Error())
			return nil, err
		}
		s.logger.Info("service wallet created", "user_id", userID, "wallet_id", wallet.ID)
	}
	return wallet, nil
}

func (s *walletService) Deposit(userID uint, amount int64, description string) (*models.Wallet, *models.Transaction, error) {

	s.logger.Info("service deposit attempt", "user_id", userID, "amount", amount)

	var result *models.Wallet
	var tran *models.Transaction

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
		tran = &transaction
		return nil
	})
	if err != nil {
		s.logger.Error("service deposit failed", "user_id", userID, "err", err.Error())
		return nil, nil, err
	}
	s.logger.Info("service deposit success", "user_id", userID, "transaction_id", tran.ID)
	return result, tran, nil
}

func (s *walletService) Freeze(userID uint, amount int64, description string) (*models.Wallet, *models.Transaction, error) {

	s.logger.Info("service freeze attempt", "user_id", userID, "amount", amount)

	var result *models.Wallet
	var tran *models.Transaction

	err := s.db.Transaction(func(tx *gorm.DB) error {
		walletRepo := s.repo.WithDB(tx)

		wallet, err := walletRepo.GetByUserIDWithLock(userID)
		if err != nil {
			return err
		}
		if wallet == nil {
			return utils.ErrWalletNotFound
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

		transaction := models.Transaction{
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

		if err := walletRepo.CreateTransaction(&transaction); err != nil {
			return err
		}

		result = wallet
		tran = &transaction
		return nil
	})
	if err != nil {
		s.logger.Error("service freeze failed", "user_id", userID, "err", err.Error())
		return nil, nil, err
	}
	s.logger.Info("service freeze success", "user_id", userID, "transaction_id", tran.ID)
	return result, tran, nil
}

func (s *walletService) Unfreeze(userID uint, amount int64, description string) (*models.Wallet, *models.Transaction, error) {

	s.logger.Info("service unfreeze attempt", "user_id", userID, "amount", amount)

	var result *models.Wallet
	var tran *models.Transaction
	err := s.db.Transaction(func(txDB *gorm.DB) error {
		walletRepo := s.repo.WithDB(txDB)

		wallet, err := walletRepo.GetByUserIDWithLock(userID)
		if err != nil {
			return err
		}
		if wallet == nil {
			return utils.ErrWalletNotFound
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
		transaction := models.Transaction{
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
		if err := walletRepo.CreateTransaction(&transaction); err != nil {
			return err
		}
		result = wallet
		tran = &transaction
		return nil
	})
	if err != nil {
		s.logger.Error("service unfreeze failed", "user_id", userID, "err", err.Error())
		return nil, nil, err
	}
	s.logger.Info("service unfreeze success", "user_id", userID, "transaction_id", tran.ID)
	return result, tran, nil
}

func (s *walletService) Charge(userID uint, amount int64, description string) (*models.Wallet, *models.Transaction, error) {

	s.logger.Info("service charge attempt", "user_id", userID, "amount", amount)

	var result *models.Wallet
	var tran *models.Transaction
	err := s.db.Transaction(func(txDB *gorm.DB) error {
		walletRepo := s.repo.WithDB(txDB)

		wallet, err := walletRepo.GetByUserIDWithLock(userID)
		if err != nil {
			return err
		}
		if wallet == nil {
			return utils.ErrWalletNotFound
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

		transaction := models.Transaction{
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
		if err := walletRepo.CreateTransaction(&transaction); err != nil {
			return err
		}
		result = wallet
		tran = &transaction
		return nil
	})
	if err != nil {
		s.logger.Error("service charge failed", "user_id", userID, "err", err.Error())
		return nil, nil, err
	}
	s.logger.Info("service charge success", "user_id", userID, "transaction_id", tran.ID)
	return result, tran, nil
}

func (s *walletService) ListTransactions(userID uint, limit, offset int) ([]models.Transaction, error) {
	s.logger.Info("service list transactions", "user_id", userID, "limit", limit, "offset", offset)
	return s.repo.ListTransactions(userID, limit, offset)
}
