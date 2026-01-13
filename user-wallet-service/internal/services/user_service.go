package services

import (
	"errors"
	"strings"
	"time"

	model "user-service/internal/models"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(email, password string, role model.Role) (*model.User, string, error)
	Login(email, password string) (*model.User, string, error)
	GetByID(id uint) (*model.User, error)
	UpdateProfile(id uint, fullName, email string) (*model.User, error)
}

type userService struct {
	repo      repository.UserRepository
	jwt       JWTService
	tokenTTL  time.Duration
	minPassLn int
}

func NewUserService(repo repository.UserRepository, jwt JWTService) UserService {
	return &userService{repo: repo, jwt: jwt, tokenTTL: 24 * time.Hour, minPassLn: 6}
}

func (s *userService) Register(email, password string, role model.Role) (*model.User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || len(password) < s.minPassLn {
		return nil, "", errors.New("invalid email or password too short")
	}
	if role == "" {
		role = model.RoleBuyer
	}
	if role != model.RoleBuyer && role != model.RoleSeller && role != model.RoleAdmin {
		return nil, "", errors.New("invalid role")
	}
	exists, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, "", err
	}
	if exists != nil {
		return nil, "", errors.New("email already registered")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	u := &model.User{Email: email, PasswordHash: string(hash), Role: role}
	if err := s.repo.Create(u); err != nil {
		return nil, "", err
	}
	token, err := s.jwt.GenerateToken(u, s.tokenTTL)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

func (s *userService) Login(email, password string) (*model.User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, "", err
	}
	if u == nil {
		return nil, "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}
	token, err := s.jwt.GenerateToken(u, s.tokenTTL)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

func (s *userService) GetByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) UpdateProfile(id uint, fullName, email string) (*model.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, errors.New("invalid email")
	}
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	if u.Email != email {
		existed, err := s.repo.FindByEmail(email)
		if err != nil {
			return nil, err
		}
		if existed != nil && existed.ID != u.ID {
			return nil, errors.New("email already in use")
		}
	}
	u.FullName = strings.TrimSpace(fullName)
	u.Email = email
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}
