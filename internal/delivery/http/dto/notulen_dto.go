package dto

import "time"

// ==========================================
// REQUEST DTO
// ==========================================

type CreateMeetingResultRequest struct {
	EmployeeID           *string    `json:"employee_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	TargetDescription    string     `json:"target_description" example:"Meningkatkan penjualan produk KUR" validate:"required"`
	TargetNominal        *int64     `json:"target_nominal" example:"500000000"`
	TargetCompletionDate *time.Time `json:"target_completion_date" example:"2024-12-31T00:00:00Z"`
}

type CreateMeetingRequest struct {
	Division       string                       `json:"division" example:"Marketing" validate:"required"`
	Title          string                       `json:"title" example:"Rapat Evaluasi Q1" validate:"required"`
	MeetingDate    time.Time                    `json:"meeting_date" example:"2024-04-01T09:00:00Z" validate:"required"`
	MeetingType    string                       `json:"meeting_type" example:"Offline" validate:"required"`
	Summary        string                       `json:"summary" example:"Evaluasi target kuartal 1 berjalan baik." validate:"required"`
	Notes          string                       `json:"notes" example:"Perlu ada peningkatan promosi di Q2."`
	Speaker        *string                      `json:"speaker" example:"Bapak Direktur"`
	ParticipantIDs []string                     `json:"participant_ids" example:"550e8400-e29b-41d4-a716-446655440000" validate:"required"`
	Results        []CreateMeetingResultRequest `json:"results" validate:"dive"`
	ImageURLs      []string                     `json:"image_urls" example:"https://storage.com/img1.jpg"`
}

type UpdateMeetingRequest struct {
	Division    string    `json:"division" example:"Marketing" validate:"required"`
	Title       string    `json:"title" example:"Rapat Evaluasi Q1" validate:"required"`
	MeetingDate time.Time `json:"meeting_date" example:"2024-04-01T09:00:00Z" validate:"required"`
	MeetingType string    `json:"meeting_type" example:"Offline" validate:"required"`
	Summary     string    `json:"summary" example:"Evaluasi target kuartal 1 berjalan baik." validate:"required"`
	Notes       string    `json:"notes" example:"Perlu ada peningkatan promosi di Q2."`
	Speaker     *string   `json:"speaker" example:"Bapak Direktur"`
}

type UpdateResultStatusRequest struct {
	AchievementStatus string `json:"achievement_status" example:"Done" validate:"required,oneof='To Do' 'In Progress' 'Done'"`
}

// ==========================================
// RESPONSE DTO
// ==========================================

type MeetingResultResponse struct {
	ID                   string            `json:"id"`
	MinuteID             string            `json:"minute_id"`
	EmployeeID           *string           `json:"employee_id"`
	Employee             *EmployeeResponse `json:"employee,omitempty"` // Re-use EmployeeResponse milikmu
	TargetDescription    string            `json:"target_description"`
	TargetNominal        *int64            `json:"target_nominal"`
	AchievementStatus    string            `json:"achievement_status"`
	TargetCompletionDate *time.Time        `json:"target_completion_date"`
	CreatedAt            time.Time         `json:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at"`
}

type MeetingParticipantResponse struct {
	ID         string            `json:"id"`
	MinuteID   string            `json:"minute_id"`
	EmployeeID string            `json:"employee_id"`
	Employee   *EmployeeResponse `json:"employee,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
}

type MeetingImageResponse struct {
	ID        string    `json:"id"`
	MinuteID  string    `json:"minute_id"`
	FileURL   string    `json:"file_url"`
	CreatedAt time.Time `json:"created_at"`
}

type MeetingMinuteResponse struct {
	ID                   string                       `json:"id"`
	Division             string                       `json:"division"`
	Title                string                       `json:"title"`
	MeetingDate          time.Time                    `json:"meeting_date"`
	MeetingType          string                       `json:"meeting_type"`
	Summary              string                       `json:"summary"`
	Notes                string                       `json:"notes"`
	Speaker              *string                      `json:"speaker"`
	NumberOfParticipants int                          `json:"number_of_participants"`
	CreatedAt            time.Time                    `json:"created_at"`
	UpdatedAt            time.Time                    `json:"updated_at"`
	Participants         []MeetingParticipantResponse `json:"participants,omitempty"`
	Results              []MeetingResultResponse      `json:"results,omitempty"`
	Images               []MeetingImageResponse       `json:"images,omitempty"`
}
