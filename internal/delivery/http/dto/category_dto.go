package dto

import "time"

// CreateCategoryRequest adalah payload untuk membuat kategori baru
type CreateCategoryRequest struct {
	Name string `json:"name" example:"Kredit" validate:"required"`
}

// CategoryResponse adalah format balikan data kategori
type CategoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	NameNorm  string    `json:"name_norm"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
