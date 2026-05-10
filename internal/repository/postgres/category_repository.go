package postgres

import (
	"context"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	model := FromDomainCategory(category)

	// Pastikan menggunakan WithContext(ctx)
	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return TranslateError(err)
	}

	category.ID = model.ID
	category.CreatedAt = model.CreatedAt
	category.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	model := FromDomainCategory(category)

	err := r.db.WithContext(ctx).Save(&model).Error
	if err != nil {
		return TranslateError(err)
	}

	category.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	var model CategoryModel

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	return model.ToDomain(), nil
}

func (r *categoryRepository) GetAll(ctx context.Context) ([]*domain.Category, error) {
	var models []CategoryModel

	err := r.db.WithContext(ctx).Find(&models).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	var categories []*domain.Category
	for _, model := range models {
		categories = append(categories, model.ToDomain())
	}

	return categories, nil
}
