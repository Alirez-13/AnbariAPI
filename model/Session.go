package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `json:"-"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (s *Session) BeforeCreate(_ *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
