package models

type Role string

const (
	RoleBuyer  Role = "buyer"
	RoleSeller Role = "seller"
	RoleAdmin  Role = "admin"
)

// User — участник платформы. Может быть покупателем, продавцом или администратором.
// Поле PasswordHash хранит хэш пароля (не сам пароль).
type User struct {
	Base
	FullName     string `gorm:"size:255" json:"full_name"`
	Email        string `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	Role         Role   `gorm:"type:varchar(16);not null;default:buyer" json:"role"`
}

type UpdateMeRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email" binding:"required,email"`
}