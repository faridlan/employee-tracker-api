package usecase

import (
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
func (u *employeeUsecase) RegisterEmployee(employee *domain.Employee) error {

	// Set waktu default jika belum diisi dari request
	if employee.EntryDate.IsZero() {
		employee.EntryDate = time.Now()
	}

	// Teruskan ke repository untuk disimpan ke database
	err := u.employeeRepo.Create(employee)
	if err != nil {
		return err
	}

	return nil
}

// GetEmployeeDetails mengambil data detail karyawan berdasarkan ID
func (u *employeeUsecase) GetEmployeeDetails(id string) (*domain.Employee, error) {

	employee, err := u.employeeRepo.GetByID(id)
	if err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Karyawan dengan ID tersebut tidak ditemukan")
	}

	return employee, nil
}
