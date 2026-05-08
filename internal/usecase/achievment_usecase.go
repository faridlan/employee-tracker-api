package usecase

import (
	"time"

	"github.com/faridlan/employee-tracker-api/internal/domain"
)

type achievementUsecase struct {
	achievementRepo domain.AchievementRepository
	targetRepo      domain.TargetRepository // Di-inject untuk validasi
}

// NewAchievementUsecase adalah constructor untuk menginisialisasi usecase achievement
func NewAchievementUsecase(achRepo domain.AchievementRepository, tgtRepo domain.TargetRepository) domain.AchievementUsecase {
	return &achievementUsecase{
		achievementRepo: achRepo,
		targetRepo:      tgtRepo,
	}
}

// RecordAchievement mencatat riwayat pencapaian (ledger) baru untuk sebuah target
func (u *achievementUsecase) RecordAchievement(targetID string, nominal int64, description string, closingDate time.Time) error {

	// Jika tanggal closing tidak dikirim dari request, gunakan waktu saat ini
	if closingDate.IsZero() {
		closingDate = time.Now()
	}

	// 2. Business Logic: Validasi Eksistensi Target
	// Kita pastikan bahwa TargetID yang dikirim benar-benar ada di database
	_, err := u.targetRepo.GetByID(targetID)
	if err != nil {
		return domain.NewError(domain.ErrNotFound, "gagal mencatat pencapaian: target tidak ditemukan")
	}

	// 3. Susun Entity Domain
	achievement := &domain.Achievement{
		TargetID:    targetID,
		Nominal:     nominal,
		Description: description,
		ClosingDate: closingDate,
	}

	// 4. Simpan ke Database via Repository
	return u.achievementRepo.Create(achievement)
}
