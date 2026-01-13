package models

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     Role   `json:"role"`
}

type AuthResponse struct {
	Token string     `json:"token"`
	User  SimpleUser `json:"user"`
}

type SimpleUser struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  Role   `json:"role"`
}
