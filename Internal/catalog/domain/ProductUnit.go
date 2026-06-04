package domain

import "github.com/shopspring/decimal"

type ProductUnit struct {
	ID         uint            `gorm:"primaryKey" json:"id"`
	ProductID  uint            `gorm:"not null;index" json:"product_id"`
	UnitName   string          `gorm:"size:50;not null" json:"unit_name"`
	Multiplier decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"multiplier"`
	IsBaseUnit bool            `gorm:"default:false" json:"is_base_unit"`
}
