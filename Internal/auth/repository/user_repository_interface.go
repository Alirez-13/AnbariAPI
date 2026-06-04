package repository

import (
	"AnbariAPI/Internal/auth/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByPhone(phone string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	Update(user *domain.User) error
}
