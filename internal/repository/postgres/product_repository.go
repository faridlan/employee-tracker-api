package postgres

import (
	"context"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	model := FromDomainProduct(product)
	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return TranslateError(err)
	}

	product.ID = model.ID
	product.CreatedAt = model.CreatedAt
	product.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	model := FromDomainProduct(product)
	err := r.db.WithContext(ctx).Save(&model).Error
	if err != nil {
		return TranslateError(err)
	}

	product.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	var model ProductModel
	err := r.db.WithContext(ctx).Preload("Category").Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, TranslateError(err)
	}
	return model.ToDomain(), nil
}

func (r *productRepository) GetAll(ctx context.Context) ([]*domain.Product, error) {
	var models []ProductModel
	err := r.db.WithContext(ctx).Preload("Category").Find(&models).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	var products []*domain.Product
	for _, m := range models {
		products = append(products, m.ToDomain())
	}
	return products, nil
}

func (r *productRepository) GetByCategoryID(ctx context.Context, categoryID string) ([]*domain.Product, error) {
	var models []ProductModel
	err := r.db.WithContext(ctx).Preload("Category").Where("category_id = ?", categoryID).Find(&models).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	var products []*domain.Product
	for _, m := range models {
		products = append(products, m.ToDomain())
	}
	return products, nil
}
