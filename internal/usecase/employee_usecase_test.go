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

func TestEmployeeUsecase_RegisterEmployee(t *testing.T) {
	mockRepo := new(mocks.EmployeeRepository)
	usecase := NewEmployeeUsecase(mockRepo)

	t.Run("Success - EntryDate provided", func(t *testing.T) {
		// Arrange
		input := domain.CreateEmployeeInput{
			Name:           "Budi",
			Position:       "Sales",
			OfficeLocation: "Jakarta",
			EntryDate:      time.Now(),
		}

		// Kita expect fungsi Create dipanggil dengan context dan pointer domain.Employee
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Employee")).Return(nil).Once()

		// Act
		result, err := usecase.RegisterEmployee(context.Background(), input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, input.Name, result.Name)
		assert.Equal(t, input.EntryDate, result.EntryDate)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Create Failed", func(t *testing.T) {
		// Arrange
		input := domain.CreateEmployeeInput{
			Name: "Andi",
		}

		expectedErr := errors.New("database error")
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Employee")).Return(expectedErr).Once()

		// Act
		result, err := usecase.RegisterEmployee(context.Background(), input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeUsecase_UpdateEmployee(t *testing.T) {
	mockRepo := new(mocks.EmployeeRepository)
	usecase := NewEmployeeUsecase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		// Arrange
		input := domain.UpdateEmployeeInput{
			ID:             "emp-1",
			Name:           "Budi Updated",
			Position:       "Manager",
			OfficeLocation: "Bandung",
			EntryDate:      time.Now(),
		}

		existingEmployee := &domain.Employee{
			ID:   "emp-1",
			Name: "Budi Lama",
		}

		// Mocking GetByID berhasil
		mockRepo.On("GetByID", mock.Anything, input.ID).Return(existingEmployee, nil).Once()
		// Mocking Update berhasil
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Employee")).Return(nil).Once()

		// Act
		result, err := usecase.UpdateEmployee(context.Background(), input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, input.Name, result.Name) // Pastikan field berubah
		assert.Equal(t, input.Position, result.Position)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Employee Not Found", func(t *testing.T) {
		// Arrange
		input := domain.UpdateEmployeeInput{ID: "emp-999"}

		// Mocking GetByID me-return error
		mockRepo.On("GetByID", mock.Anything, input.ID).Return(nil, errors.New("not found")).Once()

		// Act
		result, err := usecase.UpdateEmployee(context.Background(), input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Karyawan dengan ID tersebut tidak ditemukan")
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeUsecase_GetEmployeeDetails(t *testing.T) {
	mockRepo := new(mocks.EmployeeRepository)
	usecase := NewEmployeeUsecase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedEmployee := &domain.Employee{ID: "emp-1", Name: "Budi"}
		mockRepo.On("GetByID", mock.Anything, "emp-1").Return(expectedEmployee, nil).Once()

		result, err := usecase.GetEmployeeDetails(context.Background(), "emp-1")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedEmployee.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, "emp-unknown").Return(nil, errors.New("db err")).Once()

		result, err := usecase.GetEmployeeDetails(context.Background(), "emp-unknown")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "tidak ditemukan")
		mockRepo.AssertExpectations(t)
	})
}

func TestEmployeeUsecase_GetAllEmployees(t *testing.T) {
	mockRepo := new(mocks.EmployeeRepository)
	usecase := NewEmployeeUsecase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedEmployees := []*domain.Employee{
			{ID: "emp-1", Name: "Budi"},
			{ID: "emp-2", Name: "Andi"},
		}

		mockRepo.On("GetAll", mock.Anything).Return(expectedEmployees, nil).Once()

		result, err := usecase.GetAllEmployees(context.Background())

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Repository Failed", func(t *testing.T) {
		expectedErr := errors.New("database timeout")
		mockRepo.On("GetAll", mock.Anything).Return(nil, expectedErr).Once()

		result, err := usecase.GetAllEmployees(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}
