package usecase

import (
	"context"
	"time"

	"github.com/faridlan/employee-tracker-api/internal/domain"
)

type employeeUsecase struct {
	employeeRepo domain.EmployeeRepository
}

// NewEmployeeUsecase adalah constructor untuk menginisialisasi usecase
func NewEmployeeUsecase(repo domain.EmployeeRepository) domain.EmployeeUsecase {
	return &employeeUsecase{
		employeeRepo: repo,
	}
}

// RegisterEmployee menangani logika pendaftaran karyawan baru
func (u *employeeUsecase) RegisterEmployee(ctx context.Context, input domain.CreateEmployeeInput) (*domain.Employee, error) {

	// Set waktu default jika belum diisi dari request
	if input.EntryDate.IsZero() {
		input.EntryDate = time.Now()
	}

	employeeInput := &domain.Employee{
		Name:           input.Name,
		Position:       input.Position,
		OfficeLocation: input.OfficeLocation,
		EntryDate:      input.EntryDate,
	}

	// Teruskan ke repository untuk disimpan ke database
	err := u.employeeRepo.Create(ctx, employeeInput)
	if err != nil {
		return nil, err
	}

	return employeeInput, nil
}

// UpdateEmployee menangani logika memperbaharui data karyawan
func (u *employeeUsecase) UpdateEmployee(ctx context.Context, input domain.UpdateEmployeeInput) (*domain.Employee, error) {

	existing, err := u.employeeRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Karyawan dengan ID tersebut tidak ditemukan")
	}

	existing.Name = input.Name
	existing.Position = input.Position
	existing.OfficeLocation = input.OfficeLocation
	existing.EntryDate = input.EntryDate

	// Teruskan ke repository untuk disimpan ke database
	err = u.employeeRepo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// GetEmployeeDetails mengambil data detail karyawan berdasarkan ID
func (u *employeeUsecase) GetEmployeeDetails(ctx context.Context, id string) (*domain.Employee, error) {

	employee, err := u.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Karyawan dengan ID tersebut tidak ditemukan")
	}

	return employee, nil
}

// GetAllEmployee mengambil semua data karyawan
func (u *employeeUsecase) GetAllEmployees(ctx context.Context) ([]*domain.Employee, error) {

	employees, err := u.employeeRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return employees, nil

}
