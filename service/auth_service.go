package service

import (
	"AnbariAPI/dto"
	"AnbariAPI/model"
	"AnbariAPI/repository"
	"AnbariAPI/utils"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const SessionDuration = 24 * time.Hour

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
	Logout(sessionID uuid.UUID) error
	ValidateSession(sessionID uuid.UUID) (*model.User, error)
}

type authService struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) AuthService {
	return &authService{
		userRepository:    userRepo,
		sessionRepository: sessionRepo,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	exists, err := s.userRepository.FindByEmail(req.Email)
	if err == nil && exists != nil {
		return nil, errors.New("email already registered")
	}

	exists, err = s.userRepository.FindByPhone(req.Phone)
	if err == nil && exists != nil {
		return nil, errors.New("phone already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &model.User{
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
	}

	if err := s.userRepository.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	session, err := s.createSession(user.ID)
	if err != nil {
		return nil, errors.New("failed to create session")
	}

	return s.buildAuthResponse(user, session), nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepository.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, errors.New("failed to find user")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	s.sessionRepository.DeactivateAllByUserID(user.ID)

	session, err := s.createSession(user.ID)
	if err != nil {
		return nil, errors.New("failed to create session")
	}

	return s.buildAuthResponse(user, session), nil
}

func (s *authService) Logout(sessionID uuid.UUID) error {
	return s.sessionRepository.Deactivate(sessionID)
}

func (s *authService) ValidateSession(sessionID uuid.UUID) (*model.User, error) {
	session, err := s.sessionRepository.FindByID(sessionID)
	if err != nil {
		return nil, errors.New("invalid or expired session")
	}

	if session.ExpiresAt.Before(time.Now()) {
		s.sessionRepository.Deactivate(sessionID)
		return nil, errors.New("session expired")
	}

	user, err := s.userRepository.FindByID(session.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *authService) createSession(userID uint) (*model.Session, error) {
	expiresAt := time.Now().Add(SessionDuration)
	session := &model.Session{
		UserID:    userID,
		ExpiresAt: expiresAt,
		IsActive:  true,
	}

	if err := s.sessionRepository.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *authService) buildAuthResponse(user *model.User, session *model.Session) *dto.AuthResponse {
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