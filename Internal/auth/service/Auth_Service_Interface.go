package service

import (
	"AnbariAPI/Internal/auth/dto"
	"AnbariAPI/Internal/auth/models"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
	Logout(sessionID uuid.UUID) error
	ValidateSession(sessionID uuid.UUID) (*models.User, error)
}
