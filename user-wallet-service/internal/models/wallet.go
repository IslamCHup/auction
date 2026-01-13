package models

type Wallet struct {
	Base
	UserID        uint  `json:"user_id" gorm:"not null;index"`
	Balance       int64 `json:"balance" gorm:"not null;default:0"`
	FrozenBalance int64 `json:"frozen_balance" gorm:"not null;default:0"`
}