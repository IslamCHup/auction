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
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`

	StartDate time.Time `json:"start_date" gorm:"not null"`
	EndDate   time.Time `json:"end_date" gorm:"not null"`

	StartPrice   int64
	CurrentPrice int64 `json:"current_price" gorm:"not null"`
	MinStep      int64 `json:"min_step" gorm:"not null"`

	Status LotStatus `json:"status" gorm:"not null"`

	WinnerID uint64 `json:"winner_id" gorm:"not null"`

	CurrentBidID uint64 `json:"current_bid_id" gorm:"not null"`
	Bids         []Bid  `gorm:"foreignKey:LotModelID"`
}
