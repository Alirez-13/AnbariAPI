package service

import (
	"AnbariAPI/Internal/auth/dto"
	"AnbariAPI/Internal/auth/models"
	"AnbariAPI/Internal/auth/repository"
	"AnbariAPI/Internal/auth/service/utils"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const SessionDuration = 24 * time.Hour

// Sentinel errors — define once, check anywhere with errors.Is
var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrPhoneAlreadyRegistered = errors.New("phone already registered")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrSessionNotFound        = errors.New("session not found")
	ErrSessionExpired         = errors.New("session expired")
	ErrUserNotFound           = errors.New("user not found")
)

type authService struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	logger            *slog.Logger
}

func NewAuthService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	logger *slog.Logger,
) AuthService {
	return &authService{
		userRepository:    userRepo,
		sessionRepository: sessionRepo,
		logger:            logger,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	exists, err := s.userRepository.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("register: check email %s: %w", req.Email, err)
	}
	if exists != nil {
		return nil, ErrEmailAlreadyRegistered
	}

	exists, err = s.userRepository.FindByPhone(req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("register: check phone %s: %w", req.Phone, err)
	}
	if exists != nil {
		return nil, ErrPhoneAlreadyRegistered
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("register: hash password: %w", err)
	}

	user := &models.User{
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
	}

	if err := s.userRepository.Create(user); err != nil {
		return nil, fmt.Errorf("register: create user: %w", err)
	}

	session, err := s.createSession(user.ID)
	if err != nil {
		return nil, fmt.Errorf("register: create session for user %d: %w", user.ID, err)
	}

	return s.buildAuthResponse(user, session), nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepository.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("login: find user by email %s: %w", req.Email, err)
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	if err := s.sessionRepository.DeactivateAllByUserID(user.ID); err != nil {
		s.logger.Warn("login: failed to deactivate old sessions",
			"userID", user.ID,
			"error", err,
		)
	}

	session, err := s.createSession(user.ID)
	if err != nil {
		return nil, fmt.Errorf("login: create session for user %d: %w", user.ID, err)
	}

	return s.buildAuthResponse(user, session), nil
}

func (s *authService) Logout(sessionID uuid.UUID) error {
	if err := s.sessionRepository.Deactivate(sessionID); err != nil {
		return fmt.Errorf("logout: deactivate session %s: %w", sessionID, err)
	}
	return nil
}

func (s *authService) ValidateSession(sessionID uuid.UUID) (*models.User, error) {
	session, err := s.sessionRepository.FindByID(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("ValidateSession: session %s: %w", sessionID, ErrSessionNotFound)
		}
		return nil, fmt.Errorf("ValidateSession: find session %s: %w", sessionID, err)
	}

	if session.ExpiresAt.Before(time.Now()) {
		if err := s.sessionRepository.Deactivate(sessionID); err != nil {
			s.logger.Warn("ValidateSession: failed to deactivate expired session",
				"sessionID", sessionID,
				"error", err,
			)
		}
		return nil, fmt.Errorf("ValidateSession: session %s: %w", sessionID, ErrSessionExpired)
	}

	user, err := s.userRepository.FindByID(session.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("ValidateSession: user %d: %w", session.UserID, ErrUserNotFound)
		}
		return nil, fmt.Errorf("ValidateSession: find user %d: %w", session.UserID, err)
	}

	return user, nil
}

func (s *authService) createSession(userID uint) (*models.Session, error) {
	session := &models.Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(SessionDuration),
		IsActive:  true,
	}

	if err := s.sessionRepository.Create(session); err != nil {
		return nil, fmt.Errorf("createSession: user %d: %w", userID, err)
	}

	return session, nil
}

func (s *authService) buildAuthResponse(user *models.User, session *models.Session) *dto.AuthResponse {
	return &dto.AuthResponse{
		User: dto.UserDTO{
			ID:        user.ID,
			Email:     user.Email,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
		},
		SessionID: session.ID.String(),
		ExpiresAt: session.ExpiresAt,
	}
}
