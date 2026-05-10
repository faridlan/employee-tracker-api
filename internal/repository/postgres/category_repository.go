package postgres

import (
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

func (r *categoryRepository) Create(category *domain.Category) error {
	model := FromDomainCategory(category)

	err := r.db.Create(&model).Error
	if err != nil {
		return TranslateError(err)
	}

	category.ID = model.ID
	category.CreatedAt = model.CreatedAt
	category.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *categoryRepository) GetAll() ([]domain.Category, error) {
	var models []CategoryModel

	err := r.db.Find(&models).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	var categories []domain.Category
	for _, model := range models {
		categories = append(categories, model.ToDomain())
	}

	return categories, nil
}
