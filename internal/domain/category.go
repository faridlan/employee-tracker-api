package domain

import (
	"context"
	"time"
)

type Category struct {
	ID        string
	Name      string
	NameNorm  string
	CreatedAt time.Time
	UpdatedAt time.Time

	Products []Product
}

type CreateCategoryInput struct {
	Name string
}

type UpdateCategoryInput struct {
	ID   string
	Name string
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id string) (*Category, error)
	GetAll(ctx context.Context) ([]*Category, error)
}

type CategoryUsecase interface {
	CreateCategory(ctx context.Context, input CreateCategoryInput) (*Category, error)
	UpdateCategory(ctx context.Context, input UpdateCategoryInput) (*Category, error)
	GetCategoryByID(ctx context.Context, id string) (*Category, error)
	GetAllCategories(ctx context.Context) ([]*Category, error)
}
