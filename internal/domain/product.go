package domain

import "time"

type Product struct {
	ID         string
	Name       string
	NameNorm   string
	CategoryID string
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Category *Category
	Targets  []Target
}

type ProductRepository interface {
	Create(product *Product) error
	GetByID(id string) (*Product, error)
	GetByCategoryID(categoryID string) ([]Product, error)
}

type ProductUsecase interface {
	CreateProduct(product *Product) error
	GetProductsByCategory(categoryID string) ([]Product, error)
}
