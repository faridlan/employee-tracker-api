package domain

import "time"

type Category struct {
	ID        string
	Name      string
	NameNorm  string
	CreatedAt time.Time
	UpdatedAt time.Time

	Products []Product
}

type CategoryRepository interface {
	Create(category *Category) error
	GetAll() ([]Category, error)
}

type CategoryUsecase interface {
	CreateCategory(category *Category) error
	GetAllCategories() ([]Category, error)
}
