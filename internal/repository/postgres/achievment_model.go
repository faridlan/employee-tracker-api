package postgres

import (
	"time"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
)

type AchievementModel struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	TargetID    string         `gorm:"type:uuid;not null"`
	Nominal     int64          `gorm:"not null"`
	Description string         `gorm:"type:text"`
	ClosingDate time.Time      `gorm:"not null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relasi kembali ke Target
	Target *TargetModel `gorm:"foreignKey:TargetID"`
}

func (AchievementModel) TableName() string {
	return "achievements"
}

func (m *AchievementModel) ToDomain() domain.Achievement {
	achievement := domain.Achievement{
		ID:          m.ID,
		TargetID:    m.TargetID,
		Nominal:     m.Nominal,
		Description: m.Description,
		ClosingDate: m.ClosingDate,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}

	// Cek apakah relasi Target ikut di-load (Preload)
	if m.Target != nil {
		targetDomain := m.Target.ToDomain()
		achievement.Target = &targetDomain
	}

	return achievement
}

func FromDomainAchievement(a *domain.Achievement) AchievementModel {
	return AchievementModel{
		ID:          a.ID,
		TargetID:    a.TargetID,
		Nominal:     a.Nominal,
		Description: a.Description,
		ClosingDate: a.ClosingDate,
	}
}
