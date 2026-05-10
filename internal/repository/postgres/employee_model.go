package postgres

import (
	"time"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
)

// EmployeeModel merepresentasikan tabel 'employees' di database
type EmployeeModel struct {
	ID             string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name           string         `gorm:"not null"`
	Position       string         `gorm:"not null"`
	OfficeLocation string         `gorm:"not null"`
	EntryDate      time.Time      `gorm:"not null"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// TableName meng-override nama tabel default GORM
func (EmployeeModel) TableName() string {
	return "employees"
}

// ToDomain mengubah model database menjadi entitas domain murni
func (m *EmployeeModel) ToDomain() *domain.Employee {
	return &domain.Employee{
		ID:             m.ID,
		Name:           m.Name,
		Position:       m.Position,
		OfficeLocation: m.OfficeLocation,
		EntryDate:      m.EntryDate,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// FromDomainEmployee mengubah entitas domain menjadi model database
func FromDomainEmployee(e *domain.Employee) EmployeeModel {
	return EmployeeModel{
		ID:             e.ID,
		Name:           e.Name,
		Position:       e.Position,
		OfficeLocation: e.OfficeLocation,
		EntryDate:      e.EntryDate,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
