// # SINGLE REASON: Orchestrate user login.
package application

import (
	"errors"
	"fmt"

	"AnbariAPI/internal/auth/domain"
)

func (s *AuthService) Login(input LoginInput) (*AuthResult, error) {
	user, err := s.users.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("login: find user by email %s: %w", input.Email, err)
	}

	if !s.hasher.CheckPasswordHash(input.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	if err := s.sessions.DeactivateAllByUserID(user.ID); err != nil && s.logger != nil {
		s.logger.Warn("login: failed to deactivate old sessions", "userID", user.ID, "error", err)
	}

	session, err := s.createSession(user.ID)
	if err != nil {
		return nil, fmt.Errorf("login: create session for user %d: %w", user.ID, err)
	}

	return &AuthResult{User: user, SessionID: session.ID.String(), ExpiresAt: session.ExpiresAt}, nil
}
