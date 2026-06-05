package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	CategoryID uint     `gorm:"not null;index" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name       string   `gorm:"size:200;not null;index" json:"name"`
	Attribute  string   `gorm:"size:100" json:"attribute"`
	PackSize   float64  `gorm:"type:numeric(15,4);default:1;not null" json:"pack_size"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
