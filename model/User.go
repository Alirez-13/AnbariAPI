package model

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Phone     string    `gorm:"uniqueIndex;not null" json:"phone"`
	Password  string    `gorm:"not null" json:"-"` // "-" prevents password from being sent in JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Sessions  []Session `json:"-"`
}
