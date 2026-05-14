package usecase

import (
	"context"

	"github.com/faridlan/employee-tracker-api/internal/domain"
)

type targetUsecase struct {
	targetRepo   domain.TargetRepository
	employeeRepo domain.EmployeeRepository
	productRepo  domain.ProductRepository
}

func NewTargetUsecase(tRepo domain.TargetRepository, eRepo domain.EmployeeRepository, pRepo domain.ProductRepository) domain.TargetUsecase {
	return &targetUsecase{
		targetRepo:   tRepo,
		employeeRepo: eRepo,
		productRepo:  pRepo,
	}
}

func (u *targetUsecase) AssignTargetToEmployee(ctx context.Context, input domain.AssignTargetInput) (*domain.Target, error) {
	// 1. Validasi eksistensi Employee & Product
	if _, err := u.employeeRepo.GetByID(ctx, input.EmployeeID); err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Karyawan tidak ditemukan")
	}
	if _, err := u.productRepo.GetByID(ctx, input.ProductID); err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Produk tidak ditemukan")
	}

	// 2. Cek Duplikasi Target di periode yang sama
	existingTargets, err := u.targetRepo.GetByEmployeeAndPeriod(ctx, input.EmployeeID, input.Month, input.Year)
	if err != nil {
		return nil, err
	}

	for _, t := range existingTargets {
		if t.ProductID == input.ProductID {
			return nil, domain.NewError(domain.ErrConflict, "Karyawan sudah memiliki target untuk produk ini pada periode tersebut")
		}
	}

	// 3. Simpan Target
	target := &domain.Target{
		EmployeeID: input.EmployeeID,
		ProductID:  input.ProductID,
		Nominal:    input.Nominal,
		Month:      input.Month,
		Year:       input.Year,
	}

	if err := u.targetRepo.Create(ctx, target); err != nil {
		return nil, err
	}

	return u.targetRepo.GetByID(ctx, target.ID) // Reload untuk ambil relasi
}

func (u *targetUsecase) CalculateEmployeePerformance(ctx context.Context, employeeID string, month int, year int) (*domain.EmployeePerformance, error) {
	// Validasi eksistensi employee
	if _, err := u.employeeRepo.GetByID(ctx, employeeID); err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Karyawan tidak ditemukan")
	}

	targets, err := u.targetRepo.GetByEmployeeAndPeriod(ctx, employeeID, month, year)
	if err != nil {
		return nil, err
	}

	var totalTarget int64 = 0
	var totalAchievement int64 = 0

	for _, t := range targets {
		totalTarget += t.Nominal
		for _, ach := range t.Achievements {
			totalAchievement += ach.Nominal
		}
	}

	var percentage float64 = 0
	if totalTarget > 0 {
		percentage = (float64(totalAchievement) / float64(totalTarget)) * 100
	}

	performance := &domain.EmployeePerformance{
		EmployeeID:       employeeID,
		Month:            month,
		Year:             year,
		TotalTarget:      totalTarget,
		TotalAchievement: totalAchievement,
		Percentage:       percentage,
		Targets:          targets,
	}

	return performance, nil
}

func (u *targetUsecase) GetAllTargets(ctx context.Context, filter domain.TargetFilter) ([]*domain.Target, error) {
	return u.targetRepo.GetAll(ctx, filter)
}

// Tambahkan fungsi UpdateTargetNominal
func (u *targetUsecase) UpdateTargetNominal(ctx context.Context, input domain.UpdateTargetNominalInput) (*domain.Target, error) {
	// 1. Cek eksistensi target
	existing, err := u.targetRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Target tidak ditemukan")
	}

	// 2. Business Logic: Pastikan nominal valid
	if input.Nominal <= 0 {
		return nil, domain.NewError(domain.ErrBadParamInput, "Nominal target harus lebih besar dari 0")
	}

	// 3. Hanya ubah nominalnya saja
	existing.Nominal = input.Nominal

	// 4. Simpan perubahan
	if err := u.targetRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// Tambahkan fungsi DeleteTarget
func (u *targetUsecase) DeleteTarget(ctx context.Context, id string) error {
	// Pastikan target ada sebelum dihapus
	_, err := u.targetRepo.GetByID(ctx, id)
	if err != nil {
		return domain.NewError(domain.ErrNotFound, "Target tidak ditemukan")
	}

	return u.targetRepo.Delete(ctx, id)
}
