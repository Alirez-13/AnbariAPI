package model

//
//import (
//	"errors"
//
//	"gorm.io/gorm"
//)
//
//// UnitConversion handles conversion between pack units and base units
//type UnitConversion struct {
//	ProductID uint
//	FromUnit  string // "سطل" یا "کیلوگرم"
//	ToUnit    string // "کیلوگرم" یا "سطل"
//	Quantity  float64
//	Result    float64
//}
//
//// ConvertUnit converts between pack and base units
//func ConvertUnit(db *gorm.DB, productID uint, quantity float64, fromUnit string) (*UnitConversion, error) {
//	var product Product
//	if err := db.First(&product, productID).Error; err != nil {
//		return nil, err
//	}
//
//	conversion := &UnitConversion{
//		ProductID: productID,
//		Quantity:  quantity,
//	}
//
//	if !product.IsPackable {
//		// محصول قابل تقسیم نیست، واحد همون base unit هست
//		conversion.FromUnit = product.BaseUnit
//		conversion.ToUnit = product.BaseUnit
//		conversion.Result = quantity
//		return conversion, nil
//	}
//
//	if fromUnit == product.Unit {
//		// تبدیل از بسته به واحد پایه: مثلاً ۲ سطل × ۲۵ کیلو = ۵۰ کیلو
//		conversion.FromUnit = product.Unit
//		conversion.ToUnit = product.BaseUnit
//		conversion.Result = quantity * product.PackSize
//	} else if fromUnit == product.BaseUnit {
//		// تبدیل از واحد پایه به بسته: مثلاً ۳۰ کیلو = ۱.۲ سطل
//		conversion.FromUnit = product.BaseUnit
//		conversion.ToUnit = product.Unit
//		conversion.Result = quantity / product.PackSize
//	} else {
//		return nil, errors.New("invalid unit")
//	}
//
//	return conversion, nil
//}
