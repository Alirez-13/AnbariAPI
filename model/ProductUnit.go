package model

import "github.com/shopspring/decimal"

type ProductUnit struct {
	ID         uint            `gorm:"primaryKey" json:"id"`
	ProductID  uint            `gorm:"not null;index" json:"product_id"`
	UnitName   string          `gorm:"size:50;not null" json:"unit_name"`             // مثلا: "عدد"، "کارتن ۱۲ تایی"، "پالت ۱۲۰ تایی"
	Multiplier decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"multiplier"` // ضریب تبدیل به واحد پایه. برای واحد پایه این مقدار $1$ است.
	IsBaseUnit bool            `gorm:"default:false" json:"is_base_unit"`             // مشخص کننده واحد پایه (کوچکترین جزء)
}
