package postgres

import (
	"context"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
)

// employeeRepository adalah implementasi dari domain.EmployeeRepository
type employeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository menginisialisasi repository employee
func NewEmployeeRepository(db *gorm.DB) domain.EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}

// Create menyimpan data employee baru ke database
func (r *employeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
	model := FromDomainEmployee(employee)

	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return TranslateError(err)
	}

	// Kembalikan ID dan waktu yang di-generate database ke domain object
	employee.ID = model.ID
	employee.CreatedAt = model.CreatedAt
	employee.UpdatedAt = model.UpdatedAt

	return nil
}

// Create menyimpan data employee baru ke database
func (r *employeeRepository) Update(ctx context.Context, employee *domain.Employee) error {
	model := FromDomainEmployee(employee)

	err := r.db.WithContext(ctx).Save(&model).Error
	if err != nil {
		return TranslateError(err)
	}

	employee.UpdatedAt = model.UpdatedAt

	return nil
}

// GetByID mengambil data employee berdasarkan UUID
func (r *employeeRepository) GetByID(ctx context.Context, id string) (*domain.Employee, error) {
	var model EmployeeModel

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	domainEmployee := model.ToDomain()
	return domainEmployee, nil
}

// GetAll mengambil seluruh data employee yang belum di-soft delete
func (r *employeeRepository) GetAll(ctx context.Context) ([]*domain.Employee, error) {
	var models []EmployeeModel

	err := r.db.WithContext(ctx).Find(&models).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	var employees []*domain.Employee
	for _, model := range models {
		employees = append(employees, model.ToDomain())
	}

	return employees, nil
}
