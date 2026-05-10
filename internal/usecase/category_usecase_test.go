package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"github.com/faridlan/employee-tracker-api/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryUsecase_CreateCategory(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	usecase := NewCategoryUsecase(mockRepo)

	t.Run("Success - Generate NameNorm Correctly", func(t *testing.T) {
		// Arrange
		input := domain.CreateCategoryInput{
			Name: "Kredit Mikro",
		}

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil).Once()

		// Act
		result, err := usecase.CreateCategory(context.Background(), input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Kredit Mikro", result.Name)
		assert.Equal(t, "kredit-mikro", result.NameNorm) // Memastikan logic NameNorm berjalan benar
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Create Failed", func(t *testing.T) {
		// Arrange
		input := domain.CreateCategoryInput{
			Name: "Tabungan",
		}

		expectedErr := errors.New("db connection lost")
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(expectedErr).Once()

		// Act
		result, err := usecase.CreateCategory(context.Background(), input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryUsecase_UpdateCategory(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	usecase := NewCategoryUsecase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		// Arrange
		input := domain.UpdateCategoryInput{
			ID:   "cat-1",
			Name: "Kredit Konsumtif",
		}

		existingCat := &domain.Category{
			ID:       "cat-1",
			Name:     "Kredit Lama",
			NameNorm: "kredit-lama",
		}

		// Mock GetByID
		mockRepo.On("GetByID", mock.Anything, input.ID).Return(existingCat, nil).Once()
		// Mock Update
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil).Once()

		// Act
		result, err := usecase.UpdateCategory(context.Background(), input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Kredit Konsumtif", result.Name)
		assert.Equal(t, "kredit-konsumtif", result.NameNorm) // Memastikan NameNorm ikut terupdate
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Category Not Found", func(t *testing.T) {
		// Arrange
		input := domain.UpdateCategoryInput{ID: "cat-unknown"}

		mockRepo.On("GetByID", mock.Anything, input.ID).Return(nil, errors.New("record not found")).Once()

		// Act
		result, err := usecase.UpdateCategory(context.Background(), input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Kategori dengan ID tersebut tidak ditemukan")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Update Failed", func(t *testing.T) {
		// Arrange
		input := domain.UpdateCategoryInput{ID: "cat-1", Name: "New Name"}
		existingCat := &domain.Category{ID: "cat-1", Name: "Old Name"}

		mockRepo.On("GetByID", mock.Anything, input.ID).Return(existingCat, nil).Once()

		expectedErr := errors.New("failed to save")
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(expectedErr).Once()

		// Act
		result, err := usecase.UpdateCategory(context.Background(), input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryUsecase_GetCategoryByID(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	usecase := NewCategoryUsecase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedCat := &domain.Category{ID: "cat-1", Name: "Deposito"}
		mockRepo.On("GetByID", mock.Anything, "cat-1").Return(expectedCat, nil).Once()

		result, err := usecase.GetCategoryByID(context.Background(), "cat-1")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedCat.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, "invalid-id").Return(nil, errors.New("not found")).Once()

		result, err := usecase.GetCategoryByID(context.Background(), "invalid-id")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "tidak ditemukan")
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoryUsecase_GetAllCategories(t *testing.T) {
	mockRepo := new(mocks.CategoryRepository)
	usecase := NewCategoryUsecase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedCats := []*domain.Category{
			{ID: "cat-1", Name: "Kredit"},
			{ID: "cat-2", Name: "Tabungan"},
		}

		mockRepo.On("GetAll", mock.Anything).Return(expectedCats, nil).Once()

		result, err := usecase.GetAllCategories(context.Background())

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Failed", func(t *testing.T) {
		expectedErr := errors.New("db error")
		mockRepo.On("GetAll", mock.Anything).Return(nil, expectedErr).Once()

		result, err := usecase.GetAllCategories(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}
