package models

import (
	"gorm.io/gorm"
)

const (
	NotificationTypeBidOutbid    = "bid_outbid"
	NotificationTypeAuctionWon   = "auction_won"
	NotificationTypeAuctionLost  = "auction_lost"
	NotificationTypeAuctionEnded = "auction_ended"
)

type Notification struct {
	gorm.Model

	UserID  uint64 `gorm:"index;not null" json:"user_id"`
	LotID   uint64 `gorm:"index;not null" json:"lot_id"`
	Type    string `gorm:"type:varchar(32);not null;index" json:"type"`
	Title   string `gorm:"type:varchar(255);not null" json:"title"`
	Message string `gorm:"type:text;not null" json:"message"`
	IsRead  bool   `gorm:"default:false;index" json:"is_read"`
}

type BidPlacedEvent struct {
	LotID            uint64 `json:"lot_id"`
	PreviousLeaderID uint64 `json:"previous_leader_id"`
	NewBidAmount     int64  `json:"new_bid_amount"`
}

type LotCompletedEvent struct {
	LotID      uint64   `json:"lot_id"`
	WinnerID   uint64   `json:"winner"`
	FinalPrice int64    `json:"final_price"`
	LoserIDs   []uint64 `json:"loser_ids"`
}

type FilterNotification struct {
	UserID *uint64 `form:"-"`
	IsRead *bool   `form:"is_read"`
	Limit  int     `form:"limit" default:"20"`
	Offset int     `form:"offset" default:"0"`
	Order  string  `form:"order" default:"desc"`
}
