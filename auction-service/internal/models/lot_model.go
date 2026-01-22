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

	// Даты можно не передавать — сервис подставит значения по умолчанию.
	StartDate time.Time `json:"start_date,omitempty" binding:"omitempty" gorm:"not null"`
	EndDate   time.Time `json:"end_date,omitempty" binding:"omitempty" gorm:"not null"`

	StartPrice   int64 `json:"start_price" binding:"required,gte=1" gorm:"not null"`
	CurrentPrice int64 `json:"current_price" gorm:"not null"`
	MinStep      int64 `json:"min_step" binding:"required,gte=1" gorm:"not null"`

	// Status устанавливается в сервисе, поэтому binding не нужен.
	Status LotStatus `json:"status,omitempty" gorm:"not null"`

	SellerID uint64 `json:"seller_id" binding:"required" gorm:"not null;index"`
	WinnerID uint64 `json:"winner_id" gorm:"default:0"`

	CurrentBidID uint64 `json:"current_bid_id" gorm:"default:0"`
	Bids         []Bid  `gorm:"foreignKey:LotModelID"`
}

// UpdateLotRequest используется для обновления существующего лота
// Все поля опциональны, чтобы можно было обновлять только нужные поля
type UpdateLotRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=1,max=255"`
	Description *string    `json:"description" binding:"omitempty,min=1"`
	StartPrice  *int64     `json:"start_price" binding:"omitempty,gte=1"`
	MinStep     *int64     `json:"min_step" binding:"omitempty,gte=1"`
	EndDate     *time.Time `json:"end_date" binding:"omitempty"`
}
