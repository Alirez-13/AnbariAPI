package model

import (
	"github.com/shopspring/decimal"
)

type TransactionDetail struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	TransactionID uint   `gorm:"not null;index" json:"transaction_id"`
	ProductID     uint   `gorm:"not null;index" json:"product_id"`
	UnitName      string `gorm:"size:50;not null" json:"unit_name"` // نام واحد در زمان تراکنش (مثلا کارتن)

	// --- اطلاعات ورودی کاربر ---
	InputQuantity  decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"input_quantity"`   // تعداد وارد شده (مثلا $2$ کارتن)
	InputUnitPrice decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"input_unit_price"` // قیمت هر واحد ورودی (مثلا $1200$ برای هر کارتن)

	// --- اطلاعات پایه و محاسباتی (Snapshot) ---
	Multiplier    decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"multiplier"`      // ضریب در لحظه تراکنش (مثلا $12$)
	BaseQuantity  decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"base_quantity"`   // مقدار پایه: $InputQuantity \times Multiplier$ (مثلا $24$ عدد)
	BaseUnitPrice decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"base_unit_price"` // قیمت هر عدد: $InputUnitPrice / Multiplier$ (مثلا $100$ برای هر عدد)

	TotalPrice decimal.Decimal `gorm:"type:numeric(15,4);not null" json:"total_price"` // قیمت کل خط: $BaseQuantity \times BaseUnitPrice$
}
