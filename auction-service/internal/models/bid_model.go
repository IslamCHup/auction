package models

type Bid struct {
	Base
	Amount int64 `json:"amount" binding:"required,gte=1" gorm:"not null"`
	UserID uint  `json:"user_id" binding:"required" gorm:"not null"`

	LotModelID uint     `json:"lot_model_id" gorm:"not null"`
	LotModel   LotModel `gorm:"foreignKey:LotModelID"`
}
