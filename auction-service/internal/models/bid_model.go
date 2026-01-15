package models

type Bid struct {
	Base
	Amount int64 `json:"amount" gorm:"not null"`
	UserID uint  `json:"user_id" gorm:"not null"`

	LotModelID uint     `json:"lot_model_id" gorm:"not null"`
	LotModel   LotModel `gorm:"foreignKey:LotModelID"`
}
