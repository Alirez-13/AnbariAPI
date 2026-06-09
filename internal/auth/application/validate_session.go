// # SINGLE REASON: Orchestrate session validation.
package application

import (
	"errors"
	"fmt"
	"time"

	"AnbariAPI/internal/auth/domain"
	"github.com/google/uuid"
)

func (s *AuthService) ValidateSession(sessionID uuid.UUID) (*domain.User, error) {
	session, err := s.sessions.FindByID(sessionID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("ValidateSession: session %s: %w", sessionID, ErrSessionNotFound)
		}
		return nil, fmt.Errorf("ValidateSession: find session %s: %w", sessionID, err)
	}

	if session.ExpiresAt.Before(time.Now()) {
		if err := s.sessions.Deactivate(sessionID); err != nil && s.logger != nil {
			s.logger.Warn("ValidateSession: failed to deactivate expired session", "sessionID", sessionID, "error", err)
		}
		return nil, fmt.Errorf("ValidateSession: session %s: %w", sessionID, ErrSessionExpired)
	}

	user, err := s.users.FindByID(session.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("ValidateSession: user %d: %w", session.UserID, ErrUserNotFound)
		}
		return nil, fmt.Errorf("ValidateSession: find user %d: %w", session.UserID, err)
	}

	return user, nil
}
