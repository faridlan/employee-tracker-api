package postgres

import (
	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
	// "github.com/faridlan/employee-tracker/domain"
)

type achievementRepository struct {
	db *gorm.DB
}

func NewAchievementRepository(db *gorm.DB) domain.AchievementRepository {
	return &achievementRepository{
		db: db,
	}
}

func (r *achievementRepository) Create(achievement *domain.Achievement) error {
	model := FromDomainAchievement(achievement)

	if err := r.db.Create(&model).Error; err != nil {
		return err
	}

	achievement.ID = model.ID
	achievement.CreatedAt = model.CreatedAt
	achievement.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *achievementRepository) GetByTargetID(targetID string) ([]domain.Achievement, error) {
	var models []AchievementModel

	// Mengambil semua riwayat untuk satu target tertentu, diurutkan dari transaksi terbaru
	err := r.db.Where("target_id = ?", targetID).
		Order("closing_date desc").
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	var achievements []domain.Achievement
	for _, model := range models {
		achievements = append(achievements, model.ToDomain())
	}

	return achievements, nil
}
