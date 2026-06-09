// # SINGLE REASON: Define auth application input and output models.
package application

import (
	"time"

	"AnbariAPI/internal/auth/domain"
)

type RegisterInput struct {
	Email    string
	Phone    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	User      *domain.User
	SessionID string
	ExpiresAt time.Time
}
