package service

import (
	"AnbariAPI/Internal/auth/domain"
	"AnbariAPI/Internal/auth/dto"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
	Logout(sessionID uuid.UUID) error
	ValidateSession(sessionID uuid.UUID) (*domain.User, error)
}
