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

func (r *targetRepository) GetByEmployeeAndPeriod(ctx context.Context, employeeID string, month int, year int) ([]*domain.Target, error) {
	var models []TargetModel

	err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("Product").
		Preload("Achievements").
		Where("employee_id = ? AND month = ? AND year = ?", employeeID, month, year).
		Find(&models).Error

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
