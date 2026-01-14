package models

import "time"

type LotModel struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `json:"title" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
	StartDate time.Time `json:"start_date" gorm:"not null"`
	EndDate   time.Time `json:"end_date" gorm:"not null"`
	Status    string `json:"status" gorm:"not null"`
}