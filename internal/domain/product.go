package domain

import (
	"context"
	"time"
)

type Product struct {
	ID         string
	Name       string
	NameNorm   string
	CategoryID string
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Category *Category // Relasi ke Category
}

type CreateProductInput struct {
	Name       string
	CategoryID string
}

type UpdateProductInput struct {
	ID         string
	Name       string
	CategoryID string
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	GetAll(ctx context.Context) ([]*Product, error)
	GetByCategoryID(ctx context.Context, categoryID string) ([]*Product, error)
}

type ProductUsecase interface {
	CreateProduct(ctx context.Context, input CreateProductInput) (*Product, error)
	UpdateProduct(ctx context.Context, input UpdateProductInput) (*Product, error)
	GetProductByID(ctx context.Context, id string) (*Product, error)
	GetAllProducts(ctx context.Context) ([]*Product, error)
	GetProductsByCategoryID(ctx context.Context, categoryID string) ([]*Product, error)
}
