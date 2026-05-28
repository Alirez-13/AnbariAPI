package dto

type ProductCreateRequest struct {
	CategoryID uint    `json:"category_id" binding:"required"`
	Name       string  `json:"name" binding:"required,max=200"`
	Attribute  string  `json:"attribute,max=100"`
	Unit       string  `json:"unit" binding:"required,max=50"`
	PackSize   float64 `json:"pack_size"`
	IsPackable bool    `json:"is_packable"`
	BaseUnit   string  `json:"base_unit,max=50"`
	CurrentStock float64 `json:"current_stock"`
}

type ProductResponse struct {
	ID           uint    `json:"id"`
	CategoryID   uint    `json:"category_id"`
	Category     CategoryResponse `json:"category,omitempty"`
	Name         string  `json:"name"`
	Attribute    string  `json:"attribute"`
	Unit         string  `json:"unit"`
	PackSize     float64 `json:"pack_size"`
	IsPackable   bool    `json:"is_packable"`
	BaseUnit     string  `json:"base_unit"`
	CurrentStock float64 `json:"current_stock"`
	DisplayStock float64 `json:"display_stock"`
	DisplayUnit  string  `json:"display_unit"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
}