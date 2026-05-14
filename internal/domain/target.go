package domain

import (
	"context"
	"time"
)

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
	Achievements []Achievement // Relasi riwayat pencapaian
}

type AssignTargetInput struct {
	EmployeeID string
	ProductID  string
	Nominal    int64
	Month      int
	Year       int
}

// EmployeePerformance adalah DTO internal domain untuk hasil kalkulasi
type EmployeePerformance struct {
	EmployeeID       string
	Month            int
	Year             int
	TotalTarget      int64
	TotalAchievement int64
	Percentage       float64
	Targets          []*Target
}

type UpdateTargetNominalInput struct {
	ID      string
	Nominal int64
}

type TargetRepository interface {
	Create(ctx context.Context, target *Target) error
	GetByID(ctx context.Context, id string) (*Target, error)
	GetAll(ctx context.Context) ([]*Target, error)
	GetByEmployeeAndPeriod(ctx context.Context, employeeID string, month int, year int) ([]*Target, error)
	Update(ctx context.Context, target *Target) error
	Delete(ctx context.Context, id string) error
}

type TargetUsecase interface {
	AssignTargetToEmployee(ctx context.Context, input AssignTargetInput) (*Target, error)
	CalculateEmployeePerformance(ctx context.Context, employeeID string, month int, year int) (*EmployeePerformance, error)
	UpdateTargetNominal(ctx context.Context, input UpdateTargetNominalInput) (*Target, error)
	DeleteTarget(ctx context.Context, id string) error
	GetAllTargets(ctx context.Context) ([]*Target, error)
}
