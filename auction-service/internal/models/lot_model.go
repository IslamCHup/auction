package models

import "time"

type LotStatus string

const (
	LotStatusDraft     LotStatus = "draft"
	LotStatusActive    LotStatus = "active"
	LotStatusCompleted LotStatus = "completed"
)

type LotModel struct {
	Base
	Title       string `json:"title" binding:"required,min=1,max=255" gorm:"not null"`
	Description string `json:"description" binding:"required,min=1" gorm:"not null"`

	StartDate time.Time `json:"start_date" gorm:"not null"`
	EndDate   time.Time `json:"end_date" gorm:"not null"`

	StartPrice   int64 `json:"start_price" binding:"required,gte=1" gorm:"not null"`
	CurrentPrice int64 `json:"current_price" gorm:"not null"`
	MinStep      int64 `json:"min_step" binding:"required,gte=1" gorm:"not null"`

	Status LotStatus `json:"status" gorm:"not null"`

	SellerID uint64 `json:"seller_id" binding:"required" gorm:"not null;index"`
	WinnerID uint64 `json:"winner_id" gorm:"default:0"`

	CurrentBidID uint64 `json:"current_bid_id" gorm:"default:0"`
	Bids         []Bid  `gorm:"foreignKey:LotModelID"`
}
