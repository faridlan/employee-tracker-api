package postgres

import (
	"time"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
)

// CategoryModel merepresentasikan tabel 'categories' di database
type CategoryModel struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name      string         `gorm:"not null"`
	NameNorm  string         `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (CategoryModel) TableName() string {
	return "categories"
}

// ToDomain mengubah model database menjadi entitas domain murni
func (m *CategoryModel) ToDomain() *domain.Category {
	return &domain.Category{
		ID:        m.ID,
		Name:      m.Name,
		NameNorm:  m.NameNorm,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromDomainCategory mengubah entitas domain menjadi model database
func FromDomainCategory(c *domain.Category) CategoryModel {
	return CategoryModel{
		ID:       c.ID,
		Name:     c.Name,
		NameNorm: c.NameNorm,
	}
}
