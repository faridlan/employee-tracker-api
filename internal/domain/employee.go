package domain

import (
	"context"
	"time"
)

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

type CreateEmployeeInput struct {
	Name           string
	Position       string
	OfficeLocation string
	EntryDate      time.Time
}

type UpdateEmployeeInput struct {
	ID             string
	Name           string
	Position       string
	OfficeLocation string
	EntryDate      time.Time
}

type EmployeeRepository interface {
	Create(ctx context.Context, employee *Employee) error
	Update(ctx context.Context, employee *Employee) error
	GetByID(ctx context.Context, id string) (*Employee, error)
	GetAll(ctx context.Context) ([]*Employee, error)
}

type EmployeeUsecase interface {
	RegisterEmployee(ctx context.Context, input CreateEmployeeInput) (*Employee, error)
	UpdateEmployee(ctx context.Context, input UpdateEmployeeInput) (*Employee, error)
	GetEmployeeDetails(ctx context.Context, id string) (*Employee, error)
	GetAllEmployees(ctx context.Context) ([]*Employee, error)
}
