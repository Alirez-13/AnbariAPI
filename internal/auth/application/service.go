// # SINGLE REASON: Define auth application service dependencies.
package application

import (
	"time"

	"AnbariAPI/internal/auth/domain"
)

const SessionDuration = 24 * time.Hour

type AuthService struct {
	users    domain.UserRepository
	sessions domain.SessionRepository
	hasher   domain.PasswordHasher
	logger   domain.Logger
}

func NewAuthService(users domain.UserRepository, sessions domain.SessionRepository, hasher domain.PasswordHasher, logger domain.Logger) *AuthService {
	return &AuthService{users: users, sessions: sessions, hasher: hasher, logger: logger}
}
