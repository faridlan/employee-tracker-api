package usecase

import (
	"github.com/faridlan/employee-tracker-api/internal/domain"
)

type targetUsecase struct {
	targetRepo domain.TargetRepository
}

// NewTargetUsecase adalah constructor untuk menginisialisasi usecase target
func NewTargetUsecase(repo domain.TargetRepository) domain.TargetUsecase {
	return &targetUsecase{
		targetRepo: repo,
	}
}

// AssignTargetToEmployee menetapkan target baru untuk karyawan
func (u *targetUsecase) AssignTargetToEmployee(target *domain.Target) error {

	// Kita ambil target karyawan di bulan dan tahun tersebut
	existingTargets, err := u.targetRepo.GetByEmployeeAndPeriod(target.EmployeeID, target.Month, target.Year)
	if err != nil {
		return err
	}

	// Cek apakah di dalam target yang sudah ada, produknya sama
	for _, existingTarget := range existingTargets {
		if existingTarget.ProductID == target.ProductID {
			return domain.NewError(domain.ErrConflict, "karyawan sudah memiliki target untuk produk ini pada periode tersebut")
		}
	}

	// 3. Simpan ke Database via Repository
	return u.targetRepo.Create(target)
}

// CalculateEmployeePerformance mengambil performa karyawan berdasarkan bulan & tahun
func (u *targetUsecase) CalculateEmployeePerformance(employeeID string, month int, year int) (map[string]interface{}, error) {

	// Ambil semua target karyawan pada periode tersebut
	targets, err := u.targetRepo.GetByEmployeeAndPeriod(employeeID, month, year)
	if err != nil {
		return nil, err
	}

	var totalTargetNominal int64 = 0
	var totalAchievementNominal int64 = 0

	// Loop semua target dan hitung pencapaiannya
	// Catatan: Ini akan lebih optimal jika digabung dengan repository logic nanti,
	// tapi ini contoh business logic dasar di memori.
	for _, t := range targets {
		totalTargetNominal += t.Nominal

		// Menjumlahkan semua riwayat achievement untuk target ini
		for _, ach := range t.Achievements {
			totalAchievementNominal += ach.Nominal
		}
	}

	// Hitung persentase
	var percentage float64 = 0
	if totalTargetNominal > 0 {
		percentage = (float64(totalAchievementNominal) / float64(totalTargetNominal)) * 100
	}

	// Susun response
	result := map[string]interface{}{
		"employee_id":            employeeID,
		"month":                  month,
		"year":                   year,
		"total_target":           totalTargetNominal,
		"total_achievement":      totalAchievementNominal,
		"achievement_percentage": percentage,
		"target_details":         targets,
	}

	return result, nil
}
