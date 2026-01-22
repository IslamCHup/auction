package transport

import (
	"net/http"
	"strconv"

	"log/slog"
	"user-service/internal/models"
	m "user-service/internal/models"
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	users  services.UserService
	jwt    services.JWTService
	logger *slog.Logger
}

func NewAuthHandler(users services.UserService, jwt services.JWTService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{users: users, jwt: jwt, logger: logger}
}

func toSimple(u *m.User) m.SimpleUser {
	return m.SimpleUser{ID: u.ID, FullName: u.FullName, Email: u.Email, Role: u.Role}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req m.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("register bad request", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("register attempt", "email", req.Email, "role", string(req.Role))

	u, token, err := h.users.Register(req.Email, req.Password, req.Role)
	if err != nil {
		h.logger.Error("register failed", "email", req.Email, "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("user registered", "user_id", u.ID, "email", u.Email)
	c.JSON(http.StatusCreated, m.AuthResponse{Token: token, User: toSimple(u)})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("login bad request", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("login attempt", "email", req.Email)

	u, token, err := h.users.Login(req.Email, req.Password)
	if err != nil {
		h.logger.Warn("login failed", "email", req.Email, "err", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("login success", "user_id", u.ID, "email", u.Email)
	c.JSON(http.StatusOK, m.AuthResponse{Token: token, User: toSimple(u)})
}

func (h *AuthHandler) Me(c *gin.Context) {
	var uid uint
	
	// Попытка получить user_id из контекста (для прямых запросов)
	uidAny, exists := c.Get("user_id")
	if exists {
		uid, _ = uidAny.(uint)
	}
	
	// Если не найдено в контексте, попытаться из заголовка (для проксированных запросов)
	if uid == 0 {
		uidInt, _ := strconv.Atoi(c.GetHeader("X-User-Id"))
		uid = uint(uidInt)
	}
	
	if uid == 0 {
		h.logger.Warn("me unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	
	h.logger.Info("me request", "user_id", uid)

	u, err := h.users.GetByID(uid)
	if err != nil || u == nil {
		h.logger.Warn("me not found", "user_id", uid)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, toSimple(u))
}

func (h *AuthHandler) UpdateMe(c *gin.Context) {
	var uid uint
	
	// Попытка получить user_id из контекста (для прямых запросов)
	uidAny, exists := c.Get("user_id")
	if exists {
		uid, _ = uidAny.(uint)
	}
	
	// Если не найдено в контексте, попытаться из заголовка (для проксированных запросов)
	if uid == 0 {
		uidInt, _ := strconv.Atoi(c.GetHeader("X-User-Id"))
		uid = uint(uidInt)
	}
	
	if uid == 0 {
		h.logger.Warn("update me unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req m.UpdateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("update me bad request", "user_id", uid, "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("update profile attempt", "user_id", uid, "email", req.Email)

	u, err := h.users.UpdateProfile(uid, req.FullName, req.Email)
	if err != nil {
		h.logger.Warn("update profile failed", "user_id", uid, "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("profile updated", "user_id", u.ID)
	c.JSON(http.StatusOK, toSimple(u))
}
