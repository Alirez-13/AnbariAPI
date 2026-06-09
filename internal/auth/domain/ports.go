// # SINGLE REASON: Define auth application ports.
package domain

import "github.com/google/uuid"

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByPhone(phone string) (*User, error)
	FindByID(id uint) (*User, error)
	Update(user *User) error
}

type SessionRepository interface {
	Create(session *Session) error
	FindByID(id uuid.UUID) (*Session, error)
	FindByUserID(userID uint) ([]Session, error)
	Deactivate(id uuid.UUID) error
	DeactivateAllByUserID(userID uint) error
	DeleteExpired() error
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type Logger interface {
	Warn(msg string, args ...any)
}
