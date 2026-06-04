package dto

import "time"

type CategoryCreateRequest struct {
	Name string `json:"name" binding:"required,max=100"`
}

type CategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryListResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Total      int64              `json:"total"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
