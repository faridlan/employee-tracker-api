package dto

import "time"

type ProductRequest struct {
	Name       string `json:"name" example:"KUR Bank" validate:"required"`
	CategoryID string `json:"category_id" example:"uuid" validate:"required,uuid"`
}

type ProductResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	NameNorm   string    `json:"name_norm"`
	CategoryID string    `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProductWithCategoryResponse struct {
	ProductResponse
	Category CategoryResponse `json:"category"`
}
