package usecase

import (
	"context"

	"github.com/faridlan/employee-tracker-api/internal/domain"
)

type meetingMinuteUsecase struct {
	notulenRepo domain.MeetingMinuteRepository
}

// NewMeetingMinuteUsecase adalah constructor untuk menginisialisasi usecase notulen
func NewMeetingMinuteUsecase(repo domain.MeetingMinuteRepository) domain.MeetingMinuteUsecase {
	return &meetingMinuteUsecase{
		notulenRepo: repo,
	}
}

// CreateMeeting menangani logika pembuatan Notulen Rapat beserta Peserta, Tugas, dan Dokumentasinya
func (u *meetingMinuteUsecase) CreateMeeting(ctx context.Context, input domain.CreateMeetingInput) (*domain.MeetingMinute, error) {

	// 1. Siapkan Entitas Induk (Meeting Minute)
	meeting := &domain.MeetingMinute{
		Division:             input.Division,
		Title:                input.Title,
		MeetingDate:          input.MeetingDate,
		MeetingType:          input.MeetingType,
		Summary:              input.Summary,
		Notes:                input.Notes,
		Speaker:              input.Speaker,
		NumberOfParticipants: len(input.ParticipantIDs), // Otomatis hitung jumlah peserta dari array ID
	}

	// 2. Mapping data array ParticipantIDs menjadi struct domain.MeetingParticipant
	for _, employeeID := range input.ParticipantIDs {
		participant := domain.MeetingParticipant{
			EmployeeID: employeeID,
		}
		meeting.Participants = append(meeting.Participants, participant)
	}

	// 3. Mapping data array Results menjadi struct domain.MeetingResult
	for _, resInput := range input.Results {
		result := domain.MeetingResult{
			EmployeeID:           resInput.EmployeeID,
			TargetDescription:    resInput.TargetDescription,
			TargetNominal:        resInput.TargetNominal,
			TargetCompletionDate: resInput.TargetCompletionDate,
			AchievementStatus:    "To Do", // Set status default ke "To Do"
		}
		meeting.Results = append(meeting.Results, result)
	}

	// 4. Mapping data array ImageURLs menjadi struct domain.MeetingImage
	for _, url := range input.ImageURLs {
		img := domain.MeetingImage{
			FileURL: url,
		}
		meeting.Images = append(meeting.Images, img)
	}

	// Teruskan ke repository untuk disimpan ke database (GORM akan handle insert ke 4 tabel otomatis)
	err := u.notulenRepo.Create(ctx, meeting)
	if err != nil {
		return nil, err
	}

	return meeting, nil
}

// UpdateMeeting menangani logika pembaruan data dasar Notulen
func (u *meetingMinuteUsecase) UpdateMeeting(ctx context.Context, input domain.UpdateMeetingInput) (*domain.MeetingMinute, error) {

	existing, err := u.notulenRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Notulen rapat dengan ID tersebut tidak ditemukan")
	}

	// Update field dasar
	existing.Division = input.Division
	existing.Title = input.Title
	existing.MeetingDate = input.MeetingDate
	existing.MeetingType = input.MeetingType
	existing.Summary = input.Summary
	existing.Notes = input.Notes
	existing.Speaker = input.Speaker

	// Teruskan ke repository untuk diupdate
	err = u.notulenRepo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// GetMeetingDetails mengambil data detail Notulen beserta relasinya (Peserta, Tugas, dll)
func (u *meetingMinuteUsecase) GetMeetingDetails(ctx context.Context, id string) (*domain.MeetingMinute, error) {
	meeting, err := u.notulenRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Notulen rapat dengan ID tersebut tidak ditemukan")
	}

	return meeting, nil
}

// GetAllMeetings mengambil semua daftar Notulen rapat
func (u *meetingMinuteUsecase) GetAllMeetings(ctx context.Context) ([]*domain.MeetingMinute, error) {
	meetings, err := u.notulenRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return meetings, nil
}

// DeleteMeeting menghapus data Notulen berdasarkan ID
func (u *meetingMinuteUsecase) DeleteMeeting(ctx context.Context, id string) error {
	// Pengecekan data eksis
	_, err := u.notulenRepo.GetByID(ctx, id)
	if err != nil {
		return domain.NewError(domain.ErrNotFound, "Notulen rapat dengan ID tersebut tidak ditemukan")
	}

	return u.notulenRepo.Delete(ctx, id)
}

// ==========================================
// KHUSUS MEETING RESULT (Tugas / Action Items)
// ==========================================

// UpdateTaskStatus memungkinkan Karyawan untuk memperbarui status pekerjaannya (misal: "To Do" -> "Done")
func (u *meetingMinuteUsecase) UpdateTaskStatus(ctx context.Context, input domain.UpdateResultStatusInput) (*domain.MeetingResult, error) {

	// Validasi apakah action item tersebut ada
	result, err := u.notulenRepo.GetResultByID(ctx, input.ResultID)
	if err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Data tugas/hasil rapat tidak ditemukan")
	}

	// Timpa status yang lama dengan yang baru
	result.AchievementStatus = input.AchievementStatus

	// Simpan perubahan ke database
	err = u.notulenRepo.UpdateResult(ctx, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
