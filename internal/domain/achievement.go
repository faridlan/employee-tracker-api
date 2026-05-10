package domain

import "time"

type Achievement struct {
	ID          string
	TargetID    string
	Nominal     int64
	Description string
	ClosingDate time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Target *Target
}

type AchievementRepository interface {
	Create(achievement *Achievement) error
	GetByTargetID(targetID string) ([]Achievement, error)
}

type AchievementUsecase interface {
	// Menambah riwayat pencapaian baru (Ledger entry)
	RecordAchievement(targetID string, nominal int64, description string, closingDate time.Time) error
}
