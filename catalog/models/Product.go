package models

import (
	models2 "AnbariAPI/Internal/Inventory/models"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	CategoryID   uint            `gorm:"not null;index" json:"category_id"`
	Category     Category        `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name         string          `gorm:"size:200;not null;index" json:"name"`
	BaseUnit     string          `gorm:"size:50;not null" json:"base_unit"`
	CurrentStock decimal.Decimal `gorm:"type:numeric(15,4);default:0;not null" json:"current_stock"`

	Units   []ProductUnit            `gorm:"foreignKey:ProductID" json:"units,omitempty"`
	Batches []models2.InventoryBatch `gorm:"foreignKey:ProductID" json:"batches,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
