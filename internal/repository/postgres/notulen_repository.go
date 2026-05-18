package postgres

import (
	"context"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
)

type meetingMinuteRepository struct {
	db *gorm.DB
}

// NewMeetingMinuteRepository menginisialisasi repository untuk Notulen
func NewMeetingMinuteRepository(db *gorm.DB) domain.MeetingMinuteRepository {
	return &meetingMinuteRepository{
		db: db,
	}
}

// Create menyimpan data notulen baru beserta peserta, hasil, dan gambarnya
func (r *meetingMinuteRepository) Create(ctx context.Context, meeting *domain.MeetingMinute) error {
	model := FromDomainMeetingMinute(meeting)

	// GORM otomatis melakukan insert ke tabel meeting_minutes, participants, results, dan images secara bersamaan (cascade)
	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return TranslateError(err) // Menggunakan TranslateError dari gorm_helper.go milikmu
	}

	// Kembalikan ID yang di-generate database ke domain object
	meeting.ID = model.ID
	meeting.CreatedAt = model.CreatedAt
	meeting.UpdatedAt = model.UpdatedAt

	return nil
}

// Update menyimpan perubahan pada data induk Notulen (dipaksa mengupdate zero-value / nil)
func (r *meetingMinuteRepository) Update(ctx context.Context, meeting *domain.MeetingMinute) error {
	model := FromDomainMeetingMinute(meeting)

	// PENTING: Tambahkan Select untuk kolom induk agar jika user mengubah data menjadi
	// kosong atau nil (misal menghapus Speaker atau ExternalParticipants), nilainya tetap ter-update di DB (menjadi NULL).
	err := r.db.WithContext(ctx).
		Model(&model).
		Select("Division", "Title", "MeetingDate", "MeetingType", "Summary", "Notes", "Speaker", "NumberOfParticipants", "ExternalParticipants", "UpdatedAt").
		Updates(model).Error
	if err != nil {
		return TranslateError(err)
	}

	meeting.UpdatedAt = model.UpdatedAt
	return nil
}

// GetByID mengambil detail Notulen beserta seluruh relasinya (Preload)
func (r *meetingMinuteRepository) GetByID(ctx context.Context, id string) (*domain.MeetingMinute, error) {
	var model MeetingMinuteModel

	// Preload digunakan agar GORM otomatis melakukan JOIN/mengambil data child
	err := r.db.WithContext(ctx).
		Preload("Participants.Employee"). // Ambil data peserta beserta detail Employee-nya
		Preload("Results.Employee").      // Ambil hasil rapat beserta PIC Employee-nya
		Preload("Images").
		Where("id = ?", id).
		First(&model).Error

	if err != nil {
		return nil, TranslateError(err)
	}

	return model.ToDomain(), nil
}

// GetAll mengambil daftar Notulen (biasanya untuk list view, tidak perlu preload terlalu dalam untuk efisiensi)
func (r *meetingMinuteRepository) GetAll(ctx context.Context) ([]*domain.MeetingMinute, error) {
	var models []MeetingMinuteModel

	// Urutkan dari rapat terbaru
	err := r.db.WithContext(ctx).Order("meeting_date desc").Find(&models).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	var meetings []*domain.MeetingMinute
	for _, model := range models {
		meetings = append(meetings, model.ToDomain())
	}

	return meetings, nil
}

// Delete melakukan soft delete pada notulen
func (r *meetingMinuteRepository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&MeetingMinuteModel{}).Error
	if err != nil {
		return TranslateError(err)
	}
	return nil
}

// ==========================================
// KHUSUS MEETING RESULT (Tugas / Action Items)
// ==========================================

// GetResultByID mengambil detail satu Result/Action Item (untuk validasi sebelum update status)
func (r *meetingMinuteRepository) GetResultByID(ctx context.Context, resultID string) (*domain.MeetingResult, error) {
	var model MeetingResultModel
	err := r.db.WithContext(ctx).Preload("Employee").Where("id = ?", resultID).First(&model).Error
	if err != nil {
		return nil, TranslateError(err)
	}

	result := model.ToDomain()
	return &result, nil
}

// UpdateResult mengupdate record hasil rapat (misal: saat PIC mengubah status "To Do" menjadi "Done")
func (r *meetingMinuteRepository) UpdateResult(ctx context.Context, result *domain.MeetingResult) error {
	model := FromDomainMeetingResult(result)

	err := r.db.WithContext(ctx).Save(&model).Error
	if err != nil {
		return TranslateError(err)
	}

	result.UpdatedAt = model.UpdatedAt
	return nil
}
