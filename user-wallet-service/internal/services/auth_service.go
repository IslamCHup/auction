package services

import (
	"errors"
	"os"
	"time"

	model "user-service/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(u *model.User, ttl time.Duration) (string, error)
	ParseToken(tokenStr string) (*jwt.RegisteredClaims, model.Role, uint, error)
}

type jwtService struct {
	secret string
}

func NewJWTService() JWTService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}
	return &jwtService{secret: secret}
}

type userClaims struct {
	Role string `json:"role"`
	UID  uint   `json:"uid"`
	jwt.RegisteredClaims
}

func (s *jwtService) GenerateToken(u *model.User, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := &userClaims{
		Role: string(u.Role),
		UID:  u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Subject:   "user_auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *jwtService) ParseToken(tokenStr string) (*jwt.RegisteredClaims, model.Role, uint, error) {
	parsed, err := jwt.ParseWithClaims(tokenStr, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, "", 0, errors.New("invalid token")
	}
	cl, ok := parsed.Claims.(*userClaims)
	if !ok {
		return nil, "", 0, errors.New("invalid claims")
	}
	return &cl.RegisteredClaims, model.Role(cl.Role), cl.UID, nil
}
