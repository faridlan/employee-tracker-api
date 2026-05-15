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

func setupTargetUsecase() (*mocks.TargetRepository, *mocks.EmployeeRepository, *mocks.ProductRepository, domain.TargetUsecase) {
	mockTargetRepo := new(mocks.TargetRepository)
	mockEmployeeRepo := new(mocks.EmployeeRepository)
	mockProductRepo := new(mocks.ProductRepository)
	usecase := NewTargetUsecase(mockTargetRepo, mockEmployeeRepo, mockProductRepo)
	return mockTargetRepo, mockEmployeeRepo, mockProductRepo, usecase
}

func TestTargetUsecase_AssignTargetToEmployee(t *testing.T) {
	mockTargetRepo, mockEmployeeRepo, mockProductRepo, usecase := setupTargetUsecase()

	t.Run("Success", func(t *testing.T) {
		input := domain.AssignTargetInput{
			EmployeeID: "emp-1",
			ProductID:  "prod-1",
			Nominal:    5000000,
			Month:      5,
			Year:       2026,
		}

		// 1. Mock Employee Exist
		mockEmployeeRepo.On("GetByID", mock.Anything, input.EmployeeID).Return(&domain.Employee{ID: "emp-1"}, nil).Once()
		// 2. Mock Product Exist
		mockProductRepo.On("GetByID", mock.Anything, input.ProductID).Return(&domain.Product{ID: "prod-1"}, nil).Once()
		// 3. Mock Check Duplicate (Return empty slice = no duplicate)
		mockTargetRepo.On("GetByEmployeeAndPeriod", mock.Anything, input.EmployeeID, input.Month, input.Year).Return([]*domain.Target{}, nil).Once()
		// 4. Mock Create
		mockTargetRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Target")).Return(nil).Once()
		// 5. Mock Reload
		expectedTarget := &domain.Target{ID: "target-1", EmployeeID: "emp-1", ProductID: "prod-1"}
		mockTargetRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(expectedTarget, nil).Once()

		result, err := usecase.AssignTargetToEmployee(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "target-1", result.ID)

		mockEmployeeRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
		mockTargetRepo.AssertExpectations(t)
	})

	t.Run("Error - Employee Not Found", func(t *testing.T) {
		input := domain.AssignTargetInput{EmployeeID: "emp-unknown"}

		mockEmployeeRepo.On("GetByID", mock.Anything, input.EmployeeID).Return(nil, errors.New("not found")).Once()

		result, err := usecase.AssignTargetToEmployee(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Karyawan tidak ditemukan")
	})

	t.Run("Error - Conflict Duplicate Target", func(t *testing.T) {
		input := domain.AssignTargetInput{EmployeeID: "emp-1", ProductID: "prod-1", Month: 5, Year: 2026}

		mockEmployeeRepo.On("GetByID", mock.Anything, input.EmployeeID).Return(&domain.Employee{ID: "emp-1"}, nil).Once()
		mockProductRepo.On("GetByID", mock.Anything, input.ProductID).Return(&domain.Product{ID: "prod-1"}, nil).Once()

		// Return array yang berisi target dengan product_id yang SAMA
		existingTargets := []*domain.Target{
			{ID: "target-old", ProductID: "prod-1"},
		}
		mockTargetRepo.On("GetByEmployeeAndPeriod", mock.Anything, input.EmployeeID, input.Month, input.Year).Return(existingTargets, nil).Once()

		result, err := usecase.AssignTargetToEmployee(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Karyawan sudah memiliki target untuk produk ini")
	})
}

func TestTargetUsecase_CalculateEmployeePerformance(t *testing.T) {
	mockTargetRepo, mockEmployeeRepo, _, usecase := setupTargetUsecase()

	t.Run("Success - Calculate Percentage Correctly", func(t *testing.T) {
		employeeID := "emp-1"
		month := 5
		year := 2026

		mockEmployeeRepo.On("GetByID", mock.Anything, employeeID).Return(&domain.Employee{ID: employeeID}, nil).Once()

		// Skenario: 2 Target.
		// Target 1 = 100jt (Achieved: 50jt + 30jt)
		// Target 2 = 200jt (Achieved: 10jt)
		// Total Target = 300jt, Total Achieved = 90jt. Persentase = 30%
		targets := []*domain.Target{
			{
				Nominal: 100000000,
				Achievements: []domain.Achievement{
					{Nominal: 50000000},
					{Nominal: 30000000},
				},
			},
			{
				Nominal: 200000000,
				Achievements: []domain.Achievement{
					{Nominal: 100000000},
				},
			},
		}

		mockTargetRepo.On("GetByEmployeeAndPeriod", mock.Anything, employeeID, month, year).Return(targets, nil).Once()

		filter := domain.TargetFilter{
			Month: month,
			Year:  year,
		}
		result, err := usecase.CalculateEmployeePerformance(context.Background(), employeeID, filter)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(300000000), result.TotalTarget)
		assert.Equal(t, int64(180000000), result.TotalAchievement) // 50+30+100 = 180jt
		assert.Equal(t, float64(60), result.Percentage)            // 180/300 * 100 = 60%

		mockEmployeeRepo.AssertExpectations(t)
		mockTargetRepo.AssertExpectations(t)
	})

	t.Run("Success - Zero Total Target Guard", func(t *testing.T) {
		employeeID := "emp-1"
		mockEmployeeRepo.On("GetByID", mock.Anything, employeeID).Return(&domain.Employee{ID: employeeID}, nil).Once()

		// Return empty targets
		mockTargetRepo.On("GetByEmployeeAndPeriod", mock.Anything, employeeID, 5, 2026).Return([]*domain.Target{}, nil).Once()

		filter := domain.TargetFilter{
			Month: 5,
			Year:  2026,
		}

		result, err := usecase.CalculateEmployeePerformance(context.Background(), employeeID, filter)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, float64(0), result.Percentage) // Pastikan tidak terjadi panic / division by zero
	})
}

func TestTargetUsecase_UpdateTargetNominal(t *testing.T) {
	mockTargetRepo, _, _, usecase := setupTargetUsecase()

	t.Run("Success", func(t *testing.T) {
		input := domain.UpdateTargetNominalInput{ID: "target-1", Nominal: 200000}
		existingTarget := &domain.Target{ID: "target-1", Nominal: 100000}

		mockTargetRepo.On("GetByID", mock.Anything, input.ID).Return(existingTarget, nil).Once()
		mockTargetRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Target")).Return(nil).Once()

		result, err := usecase.UpdateTargetNominal(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, input.Nominal, result.Nominal) // Pastikan terupdate
		mockTargetRepo.AssertExpectations(t)
	})

	t.Run("Error - Invalid Nominal", func(t *testing.T) {
		input := domain.UpdateTargetNominalInput{ID: "target-1", Nominal: -500} // Nominal tidak valid
		existingTarget := &domain.Target{ID: "target-1", Nominal: 100000}

		mockTargetRepo.On("GetByID", mock.Anything, input.ID).Return(existingTarget, nil).Once()
		// Update tidak boleh dipanggil!

		result, err := usecase.UpdateTargetNominal(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Nominal target harus lebih besar dari 0")
	})
}

func TestTargetUsecase_DeleteTarget(t *testing.T) {
	mockTargetRepo, _, _, usecase := setupTargetUsecase()

	t.Run("Success", func(t *testing.T) {
		mockTargetRepo.On("GetByID", mock.Anything, "target-1").Return(&domain.Target{ID: "target-1"}, nil).Once()
		mockTargetRepo.On("Delete", mock.Anything, "target-1").Return(nil).Once()

		err := usecase.DeleteTarget(context.Background(), "target-1")

		assert.NoError(t, err)
		mockTargetRepo.AssertExpectations(t)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		mockTargetRepo.On("GetByID", mock.Anything, "target-x").Return(nil, errors.New("db error")).Once()

		err := usecase.DeleteTarget(context.Background(), "target-x")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Target tidak ditemukan")
	})
}
