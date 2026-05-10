package postgres

import (
	"context"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
)

type achievementRepository struct {
	db *gorm.DB
}

func NewAchievementRepository(db *gorm.DB) domain.AchievementRepository {
	return &achievementRepository{db: db}
}

func (r *achievementRepository) Create(ctx context.Context, achievement *domain.Achievement) error {
	model := FromDomainAchievement(achievement)
	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return TranslateError(err)
	}

	achievement.ID = model.ID
	achievement.CreatedAt = model.CreatedAt
	achievement.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *achievementRepository) GetByTargetID(ctx context.Context, targetID string) ([]*domain.Achievement, error) {
	var models []AchievementModel

	err := r.db.WithContext(ctx).
		Where("target_id = ?", targetID).
		Order("closing_date desc").
		Find(&models).Error

	if err != nil {
		return nil, TranslateError(err)
	}

	var achievements []*domain.Achievement
	for _, m := range models {
		achievements = append(achievements, m.ToDomain())
	}

	return achievements, nil
}
