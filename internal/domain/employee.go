package domain

import "time"

type Employee struct {
	ID             string
	Name           string
	Position       string
	OfficeLocation string
	EntryDate      time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time

	Targets []Target
}

type EmployeeRepository interface {
	Create(employee *Employee) error
	GetByID(id string) (*Employee, error)
	GetAll() ([]Employee, error)
}

type EmployeeUsecase interface {
	RegisterEmployee(employee *Employee) error
	GetEmployeeDetails(id string) (*Employee, error)
}
