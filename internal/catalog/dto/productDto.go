// # SINGLE REASON: Define product HTTP DTOs.
package dto

type ProductCreateRequest struct {
	CategoryID uint    `json:"category_id" binding:"required"`
	Name       string  `json:"name" binding:"required,max=200"`
	Attribute  string  `json:"attribute,max=100"`
	PackSize   float64 `json:"pack_size"`
}

type ProductResponse struct {
	ID         uint             `json:"id"`
	CategoryID uint             `json:"category_id"`
	Category   CategoryResponse `json:"category,omitempty"`
	Name       string           `json:"name"`
	Attribute  string           `json:"attribute"`
	PackSize   float64          `json:"pack_size"`
	CreatedAt  string           `json:"created_at"`
	UpdatedAt  string           `json:"updated_at"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
}
