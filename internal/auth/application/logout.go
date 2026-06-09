// # SINGLE REASON: Orchestrate user logout.
package application

import (
	"fmt"

	"github.com/google/uuid"
)

func (s *AuthService) Logout(sessionID uuid.UUID) error {
	if err := s.sessions.Deactivate(sessionID); err != nil {
		return fmt.Errorf("logout: deactivate session %s: %w", sessionID, err)
	}
	return nil
}
