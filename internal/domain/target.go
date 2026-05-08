package domain

import "time"

type Target struct {
	ID         string
	EmployeeID string
	ProductID  string
	Nominal    int64
	Month      int
	Year       int
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Employee     *Employee
	Product      *Product
	Achievements []Achievement
}

type TargetRepository interface {
	Create(target *Target) error
	GetByID(id string) (*Target, error)
	GetByEmployeeAndPeriod(employeeID string, month int, year int) ([]Target, error)
}

type TargetUsecase interface {
	AssignTargetToEmployee(target *Target) error

	// Kalkulasi performa berdasarkan target bulan & tahun tertentu
	CalculateEmployeePerformance(employeeID string, month int, year int) (map[string]interface{}, error)
}
