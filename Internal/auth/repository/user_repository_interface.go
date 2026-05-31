package repository

import (
	"AnbariAPI/Internal/auth/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByPhone(phone string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	Update(user *models.User) error
}
