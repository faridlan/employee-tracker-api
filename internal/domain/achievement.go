package domain

import (
	"context"
	"time"
)

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

type RecordAchievementInput struct {
	TargetID    string
	Nominal     int64
	Description string
	ClosingDate time.Time
}

type AchievementRepository interface {
	Create(ctx context.Context, achievement *Achievement) error
	GetByTargetID(ctx context.Context, targetID string) ([]*Achievement, error)
}

type AchievementUsecase interface {
	RecordAchievement(ctx context.Context, input RecordAchievementInput) (*Achievement, error)
	GetAchievementsByTarget(ctx context.Context, targetID string) ([]*Achievement, error)
}
