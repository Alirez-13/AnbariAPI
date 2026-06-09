// # SINGLE REASON: Persist auth users with GORM.
package infrastructure

import (
	"errors"
	"fmt"

	"AnbariAPI/internal/auth/domain"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("database error: %w", result.Error)
	}
	return &user, nil
}

func (r *GormUserRepository) FindByPhone(phone string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *GormUserRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *GormUserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}
