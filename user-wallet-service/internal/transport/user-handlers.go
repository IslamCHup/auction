package transport

import (
	"net/http"

	"user-service/internal/models"
	m "user-service/internal/models"
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	users services.UserService
	jwt   services.JWTService
}

func NewAuthHandler(users services.UserService, jwt services.JWTService) *AuthHandler {
	return &AuthHandler{users: users, jwt: jwt}
}

func toSimple(u *m.User) m.SimpleUser {
	return m.SimpleUser{ID: u.ID, Email: u.Email, Role: u.Role}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req m.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, token, err := h.users.Register(req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m.AuthResponse{Token: token, User: toSimple(u)})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, token, err := h.users.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m.AuthResponse{Token: token, User: toSimple(u)})
}

func (h *AuthHandler) Me(c *gin.Context) {
	uidAny, _ := c.Get("user_id")
	uid, _ := uidAny.(uint)
	u, err := h.users.GetByID(uid)
	if err != nil || u == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, toSimple(u))
}

func (h *AuthHandler) UpdateMe(c *gin.Context) {
	uidAny, _ := c.Get("user_id")
	uid, _ := uidAny.(uint)

	var req m.UpdateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.users.UpdateProfile(uid, req.FullName, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toSimple(u))
}
