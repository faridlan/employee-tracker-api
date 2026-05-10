package postgres

import (
	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
)

type targetRepository struct {
	db *gorm.DB
}

func NewTargetRepository(db *gorm.DB) domain.TargetRepository {
	return &targetRepository{
		db: db,
	}
}

func (r *targetRepository) Create(target *domain.Target) error {
	model := FromDomainTarget(target)

	if err := r.db.Create(&model).Error; err != nil {
		return err
	}

	target.ID = model.ID
	target.CreatedAt = model.CreatedAt
	target.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *targetRepository) GetByID(id string) (*domain.Target, error) {
	var model TargetModel

	// Eager loading untuk Employee dan Product
	err := r.db.Preload("Employee").Preload("Product").Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	domainTarget := model.ToDomain()
	return &domainTarget, nil
}

func (r *targetRepository) GetByEmployeeAndPeriod(employeeID string, month int, year int) ([]domain.Target, error) {
	var models []TargetModel

	err := r.db.Preload("Product").
		Where("employee_id = ? AND month = ? AND year = ?", employeeID, month, year).
		Find(&models).Error

	if err != nil {
		return nil, TranslateError(err)
	}

	var targets []domain.Target
	for _, model := range models {
		targets = append(targets, model.ToDomain())
	}

	return targets, nil
}
