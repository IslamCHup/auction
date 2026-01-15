package utils

import "errors"

var (
	ErrAmountMustBePositive         = errors.New("amount must be positive")
	ErrInsufficientAvailableBalance = errors.New("insufficient available balance")
	ErrResultingBalanceNegative     = errors.New("resulting balance negative")
	ErrInsufficientFrozenBalance    = errors.New("insufficient frozen balance")
	ErrWalletNotFound               = errors.New("wallet not found")
)
