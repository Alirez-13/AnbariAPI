package repository

import (
	"AnbariAPI/Internal/auth/domain"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("database error: %w", result.Error)
	}
	return &user, nil
}

func (r *userRepository) FindByPhone(phone string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}
