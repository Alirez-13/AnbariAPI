// # SINGLE REASON: Preserve legacy session repository API.
package repository

import (
	"AnbariAPI/internal/auth/domain"
	"AnbariAPI/internal/auth/infrastructure"

	"gorm.io/gorm"
)

type SessionRepository = domain.SessionRepository

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return infrastructure.NewGormSessionRepository(db)
}
