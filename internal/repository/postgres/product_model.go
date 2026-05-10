package postgres

import (
	"time"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
)

type ProductModel struct {
	ID         string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name       string         `gorm:"not null"`
	NameNorm   string         `gorm:"uniqueIndex;not null"`
	CategoryID string         `gorm:"type:uuid;not null"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	// Relasi ke CategoryModel
	Category *CategoryModel `gorm:"foreignKey:CategoryID"`
}

func (ProductModel) TableName() string {
	return "products"
}

func (m *ProductModel) ToDomain() domain.Product {
	product := domain.Product{
		ID:         m.ID,
		Name:       m.Name,
		NameNorm:   m.NameNorm,
		CategoryID: m.CategoryID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}

	// Cek apakah relasi Category ikut di-load oleh GORM (Preload)
	// Jika ada, mapping juga ke struct Domain
	if m.Category != nil {
		categoryDomain := m.Category.ToDomain()
		product.Category = &categoryDomain
	}

	return product
}

func FromDomainProduct(p *domain.Product) ProductModel {
	return ProductModel{
		ID:         p.ID,
		Name:       p.Name,
		NameNorm:   p.NameNorm,
		CategoryID: p.CategoryID,
	}
}
