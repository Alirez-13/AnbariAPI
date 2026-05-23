package model

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

	// فیلدهای جدید برای مدیریت کسری
	Unit       string  `gorm:"size:50;not null" json:"unit"`     // "کیلوگرم", "لیتر", "عدد"
	PackSize   float64 `gorm:"default:0" json:"pack_size"`       // اندازه بسته: ۲۵ برای سطل ۲۵ کیلویی
	IsPackable bool    `gorm:"default:false" json:"is_packable"` // آیا قابل تقسیم به واحدهای کوچکتر هست؟
	BaseUnit   string  `gorm:"size:50" json:"base_unit"`         // واحد پایه: "کیلوگرم" برای سطل ۲۵ کیلویی

	CurrentStock float64 `gorm:"default:0;not null" json:"current_stock"` // همیشه بر اساس واحد پایه (کیلوگرم)
	DisplayStock float64 `gorm:"-" json:"display_stock"`                  // محاسباتی: نمایش به کاربر (مثلاً ۰.۶ سطل)
	DisplayUnit  string  `gorm:"-" json:"display_unit"`                   // محاسباتی: واحد نمایش

	Transactions []TransactionDetail `gorm:"foreignKey:ProductID" json:"transactions,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	DeletedAt    gorm.DeletedAt      `gorm:"index" json:"-"`
}

// Hook برای محاسبه نمایش موجودی
func (p *Product) AfterFind(_ *gorm.DB) error {
	if p.IsPackable && p.PackSize > 0 {
		// موجودی رو به واحد بسته تبدیل کن
		p.DisplayStock = p.CurrentStock / p.PackSize
		p.DisplayUnit = p.Unit // "سطل", "گالن", etc.
	} else {
		p.DisplayStock = p.CurrentStock
		p.DisplayUnit = p.BaseUnit
	}
	return nil
}
