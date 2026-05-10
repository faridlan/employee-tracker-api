package dto

import "time"

// CategoryRequest adalah payload untuk membuat atau mengubah kategori
type CategoryRequest struct {
	Name string `json:"name" example:"Kredit Mikro" validate:"required"`
}

// CategoryResponse adalah format balikan data kategori
type CategoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	NameNorm  string    `json:"name_norm"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
