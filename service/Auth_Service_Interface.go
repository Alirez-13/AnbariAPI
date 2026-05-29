package service

import (
	"AnbariAPI/dto"
	"AnbariAPI/model"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
	Logout(sessionID uuid.UUID) error
	ValidateSession(sessionID uuid.UUID) (*model.User, error)
}
