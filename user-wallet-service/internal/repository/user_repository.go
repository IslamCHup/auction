package repository

import (
	"errors"

	"log/slog"

	model "user-service/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	Update(user *model.User) error
}

type userRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUserRepository(db *gorm.DB, logger *slog.Logger) UserRepository {
	return &userRepository{db: db, logger: logger}
}

func (r *userRepository) Create(user *model.User) error {
	if r.logger != nil {
		r.logger.Info("db create user", "email", user.Email)
	}
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		if r.logger != nil {
			r.logger.Error("db find by email failed", "email", email, "err", err.Error())
		}
		return nil, err
	}
	if r.logger != nil {
		r.logger.Info("db found user by email", "email", email, "id", user.ID)
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var u model.User
	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		if r.logger != nil {
			r.logger.Error("db find by id failed", "id", id, "err", err.Error())
		}
		return nil, err
	}
	if r.logger != nil {
		r.logger.Info("db found user by id", "id", id, "email", u.Email)
	}
	return &u, nil
}

func (r *userRepository) Update(user *model.User) error {
	if r.logger != nil {
		r.logger.Info("db update user", "id", user.ID)
	}
	return r.db.Save(user).Error
}
