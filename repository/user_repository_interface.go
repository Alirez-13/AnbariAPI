package repository

import "AnbariAPI/model"

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByPhone(phone string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	Update(user *model.User) error
}