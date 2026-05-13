package domain

import (
	"context"
	"time"
)

// ==========================================
// 1. ENTITIES (Struktur Data Utama)
// ==========================================

type MeetingMinute struct {
	ID                   string
	Division             string
	Title                string
	MeetingDate          time.Time
	MeetingType          string
	Summary              string
	Notes                string
	Speaker              *string
	NumberOfParticipants int
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time

	// Relasi
	Participants []MeetingParticipant
	Results      []MeetingResult
	Images       []MeetingImage
}

type MeetingParticipant struct {
	ID         string
	MinuteID   string
	EmployeeID string
	CreatedAt  time.Time

	Employee *Employee // Mengambil relasi dari domain/employee.go
}

type MeetingResult struct {
	ID                   string
	MinuteID             string
	EmployeeID           *string
	TargetDescription    string
	TargetNominal        *int64
	AchievementStatus    string
	TargetCompletionDate *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time

	Employee *Employee
}

type MeetingImage struct {
	ID        string
	MinuteID  string
	FileURL   string
	CreatedAt time.Time
}

// ==========================================
// 2. INPUT STRUCTS (DTO untuk Usecase)
// ==========================================

// CreateMeetingInput dirancang untuk bisa menerima data peserta dan hasil rapat sekaligus
type CreateMeetingInput struct {
	Division    string
	Title       string
	MeetingDate time.Time
	MeetingType string
	Summary     string
	Notes       string
	Speaker     *string

	ParticipantIDs []string                   // Hanya menerima array of Employee ID
	Results        []CreateMeetingResultInput // Menerima daftar hasil rapat / action items
	ImageURLs      []string                   // Menerima array link gambar
}

type CreateMeetingResultInput struct {
	EmployeeID           *string // Bisa null jika tidak ada PIC
	TargetDescription    string
	TargetNominal        *int64
	TargetCompletionDate *time.Time
}

type UpdateMeetingInput struct {
	ID          string
	Division    string
	Title       string
	MeetingDate time.Time
	MeetingType string
	Summary     string
	Notes       string
	Speaker     *string
}

// UpdateResultStatusInput khusus digunakan saat Employee ingin mengubah status tugasnya (misal: "To Do" jadi "Done")
type UpdateResultStatusInput struct {
	ResultID          string
	AchievementStatus string
}

// ==========================================
// 3. REPOSITORY INTERFACE
// ==========================================

type MeetingMinuteRepository interface {
	// CRUD Standard untuk Notulen
	Create(ctx context.Context, meeting *MeetingMinute) error
	Update(ctx context.Context, meeting *MeetingMinute) error
	GetByID(ctx context.Context, id string) (*MeetingMinute, error)
	GetAll(ctx context.Context) ([]*MeetingMinute, error)
	Delete(ctx context.Context, id string) error

	// Khusus untuk mengupdate status Action Item / Result
	UpdateResult(ctx context.Context, result *MeetingResult) error
	GetResultByID(ctx context.Context, resultID string) (*MeetingResult, error)
}

// ==========================================
// 4. USECASE INTERFACE
// ==========================================

type MeetingMinuteUsecase interface {
	// Fitur Utama Notulen
	CreateMeeting(ctx context.Context, input CreateMeetingInput) (*MeetingMinute, error)
	UpdateMeeting(ctx context.Context, input UpdateMeetingInput) (*MeetingMinute, error)
	GetMeetingDetails(ctx context.Context, id string) (*MeetingMinute, error)
	GetAllMeetings(ctx context.Context) ([]*MeetingMinute, error)
	DeleteMeeting(ctx context.Context, id string) error

	// Fitur Khusus untuk update progress tugas / action items
	UpdateTaskStatus(ctx context.Context, input UpdateResultStatusInput) (*MeetingResult, error)
}
