package postgres

import (
	"context"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
)

type targetRepository struct {
	db *gorm.DB
}

func NewTargetRepository(db *gorm.DB) domain.TargetRepository {
	return &targetRepository{db: db}
}

func (r *targetRepository) Create(ctx context.Context, target *domain.Target) error {
	model := FromDomainTarget(target)
	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return TranslateError(err)
	}

	target.ID = model.ID
	target.CreatedAt = model.CreatedAt
	target.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *targetRepository) GetByID(ctx context.Context, id string) (*domain.Target, error) {
	var model TargetModel

	err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("Product").
		Preload("Achievements").
		Where("id = ?", id).First(&model).Error

	if err != nil {
		return nil, TranslateError(err)
	}

	return model.ToDomain(), nil
}

func (r *targetRepository) GetAll(ctx context.Context, filter domain.TargetFilter) ([]*domain.Target, error) {
	var models []TargetModel

	query := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("Product").
		Preload("Achievements")

	// 1. Terapkan Filter Kondisional
	if filter.Month > 0 {
		query = query.Where("month = ?", filter.Month)
	}
	if filter.Year > 0 {
		query = query.Where("year = ?", filter.Year)
	}
	if filter.ProductID != "" {
		query = query.Where("product_id = ?", filter.ProductID)
	}

	// 2. Terapkan Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	// 3. Eksekusi
	err := query.Find(&models).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	var targets []*domain.Target
	for _, m := range models {
		targets = append(targets, m.ToDomain())
	}

	return targets, nil
}

func (r *targetRepository) GetByEmployeeAndPeriod(ctx context.Context, employeeID string, month int, year int) ([]*domain.Target, error) {
	var models []TargetModel

	// Inisiasi query dasar
	query := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("Product").
		Preload("Achievements").
		Where("employee_id = ?", employeeID)

	// Filter dinamis: Jika > 0, tambahkan kondisi WHERE
	if month > 0 {
		query = query.Where("month = ?", month)
	}
	if year > 0 {
		query = query.Where("year = ?", year)
	}

	// Eksekusi query
	err := query.Find(&models).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	var targets []*domain.Target
	for _, m := range models {
		targets = append(targets, m.ToDomain())
	}

	return targets, nil
}

func (r *targetRepository) Update(ctx context.Context, target *domain.Target) error {
	model := FromDomainTarget(target)

	// Menggunakan Save akan memperbaharui seluruh field sesuai ID
	err := r.db.WithContext(ctx).Save(&model).Error
	if err != nil {
		return TranslateError(err)
	}

	target.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *targetRepository) Delete(ctx context.Context, id string) error {
	// GORM otomatis melakukan Soft Delete karena ada field DeletedAt
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&TargetModel{}).Error
	if err != nil {
		return TranslateError(err)
	}
	return nil
}
