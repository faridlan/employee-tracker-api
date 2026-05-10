package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"github.com/faridlan/employee-tracker-api/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAchievementUsecase() (*mocks.AchievementRepository, *mocks.TargetRepository, domain.AchievementUsecase) {
	mockAchievementRepo := new(mocks.AchievementRepository)
	mockTargetRepo := new(mocks.TargetRepository)
	usecase := NewAchievementUsecase(mockAchievementRepo, mockTargetRepo)
	return mockAchievementRepo, mockTargetRepo, usecase
}

func TestAchievementUsecase_RecordAchievement(t *testing.T) {
	mockAchievementRepo, mockTargetRepo, usecase := setupAchievementUsecase()

	t.Run("Success - With Provided Date", func(t *testing.T) {
		customDate := time.Date(2026, 5, 10, 14, 0, 0, 0, time.UTC)
		input := domain.RecordAchievementInput{
			TargetID:    "target-1",
			Nominal:     50000000,
			Description: "Closing Nasabah A",
			ClosingDate: customDate,
		}

		// Mock target exist
		mockTargetRepo.On("GetByID", mock.Anything, input.TargetID).Return(&domain.Target{ID: "target-1"}, nil).Once()
		// Mock create achievement
		mockAchievementRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Achievement")).Return(nil).Once()

		result, err := usecase.RecordAchievement(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, customDate, result.ClosingDate) // Pastikan tanggal sesuai input

		mockTargetRepo.AssertExpectations(t)
		mockAchievementRepo.AssertExpectations(t)
	})

	t.Run("Success - Auto Set Default Date", func(t *testing.T) {
		input := domain.RecordAchievementInput{
			TargetID: "target-1",
			Nominal:  10000000,
			// ClosingDate sengaja dikosongkan (Zero value)
		}

		mockTargetRepo.On("GetByID", mock.Anything, input.TargetID).Return(&domain.Target{ID: "target-1"}, nil).Once()
		mockAchievementRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Achievement")).Return(nil).Once()

		result, err := usecase.RecordAchievement(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.ClosingDate.IsZero()) // Pastikan sistem mengisi tanggal otomatis (tidak zero)

		mockTargetRepo.AssertExpectations(t)
		mockAchievementRepo.AssertExpectations(t)
	})

	t.Run("Error - Target Not Found", func(t *testing.T) {
		input := domain.RecordAchievementInput{TargetID: "target-unknown"}

		mockTargetRepo.On("GetByID", mock.Anything, input.TargetID).Return(nil, errors.New("not found")).Once()

		result, err := usecase.RecordAchievement(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Target tidak ditemukan")
	})

	t.Run("Error - Repository Create Failed", func(t *testing.T) {
		input := domain.RecordAchievementInput{TargetID: "target-1"}

		mockTargetRepo.On("GetByID", mock.Anything, input.TargetID).Return(&domain.Target{ID: "target-1"}, nil).Once()

		expectedErr := errors.New("database connection error")
		mockAchievementRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Achievement")).Return(expectedErr).Once()

		result, err := usecase.RecordAchievement(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestAchievementUsecase_GetAchievementsByTarget(t *testing.T) {
	mockAchievementRepo, mockTargetRepo, usecase := setupAchievementUsecase()

	t.Run("Success", func(t *testing.T) {
		targetID := "target-1"
		expectedAchievements := []*domain.Achievement{
			{ID: "ach-1", TargetID: targetID, Nominal: 10000000},
			{ID: "ach-2", TargetID: targetID, Nominal: 20000000},
		}

		// Mock target exist
		mockTargetRepo.On("GetByID", mock.Anything, targetID).Return(&domain.Target{ID: targetID}, nil).Once()
		// Mock get achievements
		mockAchievementRepo.On("GetByTargetID", mock.Anything, targetID).Return(expectedAchievements, nil).Once()

		result, err := usecase.GetAchievementsByTarget(context.Background(), targetID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockTargetRepo.AssertExpectations(t)
		mockAchievementRepo.AssertExpectations(t)
	})

	t.Run("Error - Target Not Found", func(t *testing.T) {
		targetID := "target-x"

		mockTargetRepo.On("GetByID", mock.Anything, targetID).Return(nil, errors.New("not found")).Once()

		result, err := usecase.GetAchievementsByTarget(context.Background(), targetID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Target tidak ditemukan")
	})

	t.Run("Error - Repository Get Failed", func(t *testing.T) {
		targetID := "target-1"

		mockTargetRepo.On("GetByID", mock.Anything, targetID).Return(&domain.Target{ID: targetID}, nil).Once()

		expectedErr := errors.New("db error")
		mockAchievementRepo.On("GetByTargetID", mock.Anything, targetID).Return(nil, expectedErr).Once()

		result, err := usecase.GetAchievementsByTarget(context.Background(), targetID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}
