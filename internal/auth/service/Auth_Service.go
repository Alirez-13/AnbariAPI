// # SINGLE REASON: Preserve legacy auth service constructor and errors.
package service

import (
	"log/slog"

	"AnbariAPI/internal/auth/application"
	"AnbariAPI/internal/auth/domain"
	"AnbariAPI/internal/auth/dto"
	"AnbariAPI/internal/auth/infrastructure"
	"AnbariAPI/internal/auth/repository"

	"github.com/google/uuid"
)

const SessionDuration = application.SessionDuration

var (
	ErrEmailAlreadyRegistered = application.ErrEmailAlreadyRegistered
	ErrPhoneAlreadyRegistered = application.ErrPhoneAlreadyRegistered
	ErrInvalidCredentials     = application.ErrInvalidCredentials
	ErrSessionNotFound        = application.ErrSessionNotFound
	ErrSessionExpired         = application.ErrSessionExpired
	ErrUserNotFound           = application.ErrUserNotFound
)

type authService struct {
	inner *application.AuthService
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, logger *slog.Logger) AuthService {
	if logger == nil {
		logger = slog.Default()
	}
	return &authService{inner: application.NewAuthService(userRepo, sessionRepo, infrastructure.NewBcryptPasswordHasher(), logger)}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	result, err := s.inner.Register(application.RegisterInput(req))
	if err != nil {
		return nil, err
	}
	return authResponse(result), nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	result, err := s.inner.Login(application.LoginInput(req))
	if err != nil {
		return nil, err
	}
	return authResponse(result), nil
}

func (s *authService) Logout(sessionID uuid.UUID) error {
	return s.inner.Logout(sessionID)
}

func (s *authService) ValidateSession(sessionID uuid.UUID) (*domain.User, error) {
	return s.inner.ValidateSession(sessionID)
}

func authResponse(result *application.AuthResult) *dto.AuthResponse {
	return &dto.AuthResponse{
		User:      dto.UserDTO{ID: result.User.ID, Email: result.User.Email, Phone: result.User.Phone, CreatedAt: result.User.CreatedAt},
		SessionID: result.SessionID,
		ExpiresAt: result.ExpiresAt,
	}
}
