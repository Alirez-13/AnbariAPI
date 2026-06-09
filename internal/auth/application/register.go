// # SINGLE REASON: Orchestrate user registration.
package application

import (
	"errors"
	"fmt"
	"time"

	"AnbariAPI/internal/auth/domain"
	"github.com/google/uuid"
)

func (s *AuthService) Register(input RegisterInput) (*AuthResult, error) {
	exists, err := s.users.FindByEmail(input.Email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("register: check email %s: %w", input.Email, err)
	}
	if exists != nil {
		return nil, ErrEmailAlreadyRegistered
	}

	exists, err = s.users.FindByPhone(input.Phone)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("register: check phone %s: %w", input.Phone, err)
	}
	if exists != nil {
		return nil, ErrPhoneAlreadyRegistered
	}

	hashedPassword, err := s.hasher.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("register: hash password: %w", err)
	}

	user := &domain.User{Email: input.Email, Phone: input.Phone, Password: hashedPassword}
	if err := s.users.Create(user); err != nil {
		return nil, fmt.Errorf("register: create user: %w", err)
	}

	session, err := s.createSession(user.ID)
	if err != nil {
		return nil, fmt.Errorf("register: create session for user %d: %w", user.ID, err)
	}

	return &AuthResult{User: user, SessionID: session.ID.String(), ExpiresAt: session.ExpiresAt}, nil
}

func (s *AuthService) createSession(userID uint) (*domain.Session, error) {
	session := &domain.Session{ID: uuid.New(), UserID: userID, ExpiresAt: timeNow().Add(SessionDuration), IsActive: true}
	if err := s.sessions.Create(session); err != nil {
		return nil, fmt.Errorf("createSession: user %d: %w", userID, err)
	}
	return session, nil
}

func timeNow() time.Time {
	return time.Now()
}
