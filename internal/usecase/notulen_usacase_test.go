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

func TestMeetingMinuteUsecase_CreateMeeting(t *testing.T) {
	mockRepo := new(mocks.MeetingMinuteRepository)
	usecase := NewMeetingMinuteUsecase(mockRepo)

	t.Run("Success - Create full meeting", func(t *testing.T) {
		// Arrange
		empID := "emp-1"
		var targetNominal int64 = 50000000
		now := time.Now()

		input := domain.CreateMeetingInput{
			Division:       "Marketing",
			Title:          "Rapat Q1",
			MeetingDate:    now,
			MeetingType:    "Offline",
			ParticipantIDs: []string{empID, "emp-2"},
			Results: []domain.CreateMeetingResultInput{
				{
					EmployeeID:        &empID,
					TargetDescription: "Capai target Q1",
					TargetNominal:     &targetNominal,
				},
			},
			ImageURLs: []string{"url1", "url2"},
		}

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.MeetingMinute")).Return(nil).Once()

		// Act
		result, err := usecase.CreateMeeting(context.Background(), input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, input.Title, result.Title)

		// Verifikasi logika usecase berjalan benar
		assert.Equal(t, 2, result.NumberOfParticipants) // Hitung array otomatis
		assert.Len(t, result.Participants, 2)
		assert.Len(t, result.Results, 1)
		assert.Len(t, result.Images, 2)
		assert.Equal(t, "To Do", result.Results[0].AchievementStatus) // Default status

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Create Failed", func(t *testing.T) {
		// Arrange
		input := domain.CreateMeetingInput{Title: "Rapat Error"}
		expectedErr := errors.New("db error")

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.MeetingMinute")).Return(expectedErr).Once()

		// Act
		result, err := usecase.CreateMeeting(context.Background(), input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestMeetingMinuteUsecase_UpdateMeeting(t *testing.T) {
	mockRepo := new(mocks.MeetingMinuteRepository)
	usecase := NewMeetingMinuteUsecase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		// Arrange
		input := domain.UpdateMeetingInput{
			ID:       "meet-1",
			Division: "IT",
			Title:    "Sprint Planning Updated",
		}

		existingMeeting := &domain.MeetingMinute{
			ID:       "meet-1",
			Division: "IT Lama",
			Title:    "Sprint Planning",
		}

		mockRepo.On("GetByID", mock.Anything, input.ID).Return(existingMeeting, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.MeetingMinute")).Return(nil).Once()

		// Act
		result, err := usecase.UpdateMeeting(context.Background(), input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, input.Title, result.Title) // Pastikan field berubah
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Meeting Not Found", func(t *testing.T) {
		// Arrange
		input := domain.UpdateMeetingInput{ID: "meet-999"}

		mockRepo.On("GetByID", mock.Anything, input.ID).Return(nil, errors.New("not found")).Once()

		// Act
		result, err := usecase.UpdateMeeting(context.Background(), input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "tidak ditemukan")
		mockRepo.AssertExpectations(t)
	})
}

func TestMeetingMinuteUsecase_GetMeetingDetails(t *testing.T) {
	mockRepo := new(mocks.MeetingMinuteRepository)
	usecase := NewMeetingMinuteUsecase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedMeeting := &domain.MeetingMinute{ID: "meet-1", Title: "Rapat Q1"}
		mockRepo.On("GetByID", mock.Anything, "meet-1").Return(expectedMeeting, nil).Once()

		result, err := usecase.GetMeetingDetails(context.Background(), "meet-1")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedMeeting.Title, result.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, "meet-unknown").Return(nil, errors.New("db err")).Once()

		result, err := usecase.GetMeetingDetails(context.Background(), "meet-unknown")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "tidak ditemukan")
		mockRepo.AssertExpectations(t)
	})
}

func TestMeetingMinuteUsecase_GetAllMeetings(t *testing.T) {
	mockRepo := new(mocks.MeetingMinuteRepository)
	usecase := NewMeetingMinuteUsecase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedMeetings := []*domain.MeetingMinute{
			{ID: "meet-1", Title: "Rapat 1"},
			{ID: "meet-2", Title: "Rapat 2"},
		}

		mockRepo.On("GetAll", mock.Anything).Return(expectedMeetings, nil).Once()

		result, err := usecase.GetAllMeetings(context.Background())

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Failed", func(t *testing.T) {
		expectedErr := errors.New("database error")
		mockRepo.On("GetAll", mock.Anything).Return(nil, expectedErr).Once()

		result, err := usecase.GetAllMeetings(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestMeetingMinuteUsecase_DeleteMeeting(t *testing.T) {
	mockRepo := new(mocks.MeetingMinuteRepository)
	usecase := NewMeetingMinuteUsecase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, "meet-1").Return(&domain.MeetingMinute{}, nil).Once()
		mockRepo.On("Delete", mock.Anything, "meet-1").Return(nil).Once()

		err := usecase.DeleteMeeting(context.Background(), "meet-1")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, "meet-999").Return(nil, errors.New("not found")).Once()

		err := usecase.DeleteMeeting(context.Background(), "meet-999")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tidak ditemukan")
		mockRepo.AssertExpectations(t)
	})
}

func TestMeetingMinuteUsecase_UpdateTaskStatus(t *testing.T) {
	mockRepo := new(mocks.MeetingMinuteRepository)
	usecase := NewMeetingMinuteUsecase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		input := domain.UpdateResultStatusInput{
			ResultID:          "res-1",
			AchievementStatus: "Done",
		}

		existingResult := &domain.MeetingResult{
			ID:                "res-1",
			AchievementStatus: "To Do",
		}

		mockRepo.On("GetResultByID", mock.Anything, input.ResultID).Return(existingResult, nil).Once()
		mockRepo.On("UpdateResult", mock.Anything, mock.AnythingOfType("*domain.MeetingResult")).Return(nil).Once()

		result, err := usecase.UpdateTaskStatus(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Done", result.AchievementStatus) // Status harus terupdate
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Task Not Found", func(t *testing.T) {
		input := domain.UpdateResultStatusInput{ResultID: "res-unknown"}

		mockRepo.On("GetResultByID", mock.Anything, input.ResultID).Return(nil, errors.New("not found")).Once()

		result, err := usecase.UpdateTaskStatus(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "tidak ditemukan")
		mockRepo.AssertExpectations(t)
	})
}
