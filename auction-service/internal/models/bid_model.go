package models

type Bid struct {
	Base
	Amount int64 `json:"amount" binding:"required,gte=1" gorm:"not null"`
	UserID uint  `json:"user_id" binding:"required" gorm:"not null"`

	// LotModelID заполняется из URL, не ожидается во входящем JSON
	LotModelID uint `json:"-" gorm:"not null"`
	// Не валидируем и не принимаем из JSON вложенную модель лота при создании ставки
	LotModel *LotModel `json:"-" binding:"-" gorm:"foreignKey:LotModelID"`
}
