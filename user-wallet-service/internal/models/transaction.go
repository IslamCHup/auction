package models

type TransactionType string

const (
	TransactionDeposit  TransactionType = "deposit"
	TransactionFreeze   TransactionType = "freeze"
	TransactionUnfreeze TransactionType = "unfreeze"
	TransactionCharge   TransactionType = "charge"
)

type Transaction struct {
	Base

	WalletID uint            `json:"wallet_id" gorm:"index;not null"`
	UserID   uint            `json:"user_id" gorm:"index;not null"`
	Type     TransactionType `json:"type" gorm:"type:varchar(32);not null"`
	Amount   int64           `json:"amount" gorm:"not null"`

	BalanceBefore int64 `json:"balance_before" gorm:"not null"`
	BalanceAfter  int64 `json:"balance_after" gorm:"not null"`

	FrozenBefore int64 `json:"frozen_before" gorm:"not null"`
	FrozenAfter  int64 `json:"frozen_after" gorm:"not null"`

	Description string `json:"description" gorm:"size:512"`
}
