package postgres

import (
	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(product *domain.Product) error {
	model := FromDomainProduct(product)

	err := r.db.Create(&model).Error
	if err != nil {
		return TranslateError(err)
	}

	product.ID = model.ID
	product.CreatedAt = model.CreatedAt
	product.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *productRepository) GetByID(id string) (*domain.Product, error) {
	var model ProductModel

	// Menggunakan Preload untuk mengambil data relasi Category
	err := r.db.Preload("Category").Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	domainProduct := model.ToDomain()
	return &domainProduct, nil
}

func (r *productRepository) GetByCategoryID(categoryID string) ([]domain.Product, error) {
	var models []ProductModel

	// Untuk list product berdasarkan kategori, opsional apakah butuh Preload atau tidak
	err := r.db.Where("category_id = ?", categoryID).Find(&models).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	var products []domain.Product
	for _, model := range models {
		products = append(products, model.ToDomain())
	}

	return products, nil
}
