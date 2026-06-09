// # SINGLE REASON: Define auth session entity.
package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `json:"-"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
