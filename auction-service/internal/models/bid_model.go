package models

import "time"

type Bid struct {
	ID        uint      `gorm:"primaryKey"`
	Amount    float64   `json:"amount" gorm:"not null"`
	BidderID  uint64    `json:"bidder_id" gorm:"not null"`
	LotID     uint64    `json:"lot_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}
