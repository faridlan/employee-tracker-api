package postgres

import (
	"time"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
)

type TargetModel struct {
	ID         string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	EmployeeID string         `gorm:"type:uuid;not null"`
	ProductID  string         `gorm:"type:uuid;not null"`
	Nominal    int64          `gorm:"not null"`
	Month      int            `gorm:"not null"`
	Year       int            `gorm:"not null"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	Employee     *EmployeeModel     `gorm:"foreignKey:EmployeeID"`
	Product      *ProductModel      `gorm:"foreignKey:ProductID"`
	Achievements []AchievementModel `gorm:"foreignKey:TargetID"`
}

func (TargetModel) TableName() string {
	return "targets"
}

func (m *TargetModel) ToDomain() *domain.Target {
	target := &domain.Target{
		ID:         m.ID,
		EmployeeID: m.EmployeeID,
		ProductID:  m.ProductID,
		Nominal:    m.Nominal,
		Month:      m.Month,
		Year:       m.Year,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}

	if m.Employee != nil {
		target.Employee = m.Employee.ToDomain()
	}

	if m.Product != nil {
		target.Product = m.Product.ToDomain()
	}

	// Mapping achievements jika di-preload
	if len(m.Achievements) > 0 {
		for _, ach := range m.Achievements {
			target.Achievements = append(target.Achievements, *ach.ToDomain())
		}
	}

	return target
}

func FromDomainTarget(t *domain.Target) TargetModel {
	return TargetModel{
		ID:         t.ID,
		EmployeeID: t.EmployeeID,
		ProductID:  t.ProductID,
		Nominal:    t.Nominal,
		Month:      t.Month,
		Year:       t.Year,
	}
}
