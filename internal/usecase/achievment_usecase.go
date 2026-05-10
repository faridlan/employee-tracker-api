package usecase

import (
	"context"
	"time"

	"github.com/faridlan/employee-tracker-api/internal/domain"
)

type achievementUsecase struct {
	achievementRepo domain.AchievementRepository
	targetRepo      domain.TargetRepository // Di-inject untuk validasi
}

func NewAchievementUsecase(aRepo domain.AchievementRepository, tRepo domain.TargetRepository) domain.AchievementUsecase {
	return &achievementUsecase{
		achievementRepo: aRepo,
		targetRepo:      tRepo,
	}
}

func (u *achievementUsecase) RecordAchievement(ctx context.Context, input domain.RecordAchievementInput) (*domain.Achievement, error) {
	// 1. Validasi Eksistensi Target
	if _, err := u.targetRepo.GetByID(ctx, input.TargetID); err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Target tidak ditemukan")
	}

	// 2. Set default date jika tidak dikirim dari request
	if input.ClosingDate.IsZero() {
		input.ClosingDate = time.Now()
	}

	achievement := &domain.Achievement{
		TargetID:    input.TargetID,
		Nominal:     input.Nominal,
		Description: input.Description,
		ClosingDate: input.ClosingDate,
	}

	// 3. Simpan ke database
	if err := u.achievementRepo.Create(ctx, achievement); err != nil {
		return nil, err
	}

	return achievement, nil
}

func (u *achievementUsecase) GetAchievementsByTarget(ctx context.Context, targetID string) ([]*domain.Achievement, error) {
	// Validasi target ada atau tidak
	if _, err := u.targetRepo.GetByID(ctx, targetID); err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Target tidak ditemukan")
	}

	return u.achievementRepo.GetByTargetID(ctx, targetID)
}
