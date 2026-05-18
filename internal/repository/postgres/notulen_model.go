package postgres

import (
	"time"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"gorm.io/gorm"
)

// ==========================================
// MODELS
// ==========================================

type MeetingMinuteModel struct {
	ID                   string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Division             string         `gorm:"not null"`
	Title                string         `gorm:"not null"`
	MeetingDate          time.Time      `gorm:"not null"`
	MeetingType          string         `gorm:"not null"`
	Summary              string         `gorm:"type:text"`
	Notes                string         `gorm:"type:text"`
	Speaker              *string        `gorm:"type:varchar(255)"`
	NumberOfParticipants int            `gorm:"default:0"`
	ExternalParticipants *string        `gorm:"type:text"`
	CreatedAt            time.Time      `gorm:"autoCreateTime"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime"`
	DeletedAt            gorm.DeletedAt `gorm:"index"`

	// Relasi
	Participants []MeetingParticipantModel `gorm:"foreignKey:MinuteID;constraint:OnDelete:CASCADE;"`
	Results      []MeetingResultModel      `gorm:"foreignKey:MinuteID;constraint:OnDelete:CASCADE;"`
	Images       []MeetingImageModel       `gorm:"foreignKey:MinuteID;constraint:OnDelete:CASCADE;"`
}

func (MeetingMinuteModel) TableName() string { return "meeting_minutes" }

type MeetingParticipantModel struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	MinuteID   string    `gorm:"type:uuid;not null"`
	EmployeeID string    `gorm:"type:uuid;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	Employee *EmployeeModel `gorm:"foreignKey:EmployeeID"`
}

func (MeetingParticipantModel) TableName() string { return "meeting_participants" }

type MeetingResultModel struct {
	ID                   string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	MinuteID             string  `gorm:"type:uuid;not null"`
	EmployeeID           *string `gorm:"type:uuid"`
	TargetDescription    string  `gorm:"type:text;not null"`
	TargetNominal        *int64  `gorm:"type:bigint"`
	AchievementStatus    string  `gorm:"default:'To Do'"`
	TargetCompletionDate *time.Time
	CreatedAt            time.Time `gorm:"autoCreateTime"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime"`

	Employee *EmployeeModel `gorm:"foreignKey:EmployeeID"`
}

func (MeetingResultModel) TableName() string { return "meeting_results" }

type MeetingImageModel struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	MinuteID  string    `gorm:"type:uuid;not null"`
	FileURL   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (MeetingImageModel) TableName() string { return "meeting_images" }

// ==========================================
// MAPPERS (Database -> Domain)
// ==========================================

func (m *MeetingMinuteModel) ToDomain() *domain.MeetingMinute {
	minute := &domain.MeetingMinute{
		ID:                   m.ID,
		Division:             m.Division,
		Title:                m.Title,
		MeetingDate:          m.MeetingDate,
		MeetingType:          m.MeetingType,
		Summary:              m.Summary,
		Notes:                m.Notes,
		Speaker:              m.Speaker,
		NumberOfParticipants: m.NumberOfParticipants,
		ExternalParticipants: m.ExternalParticipants,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}

	if m.DeletedAt.Valid {
		minute.DeletedAt = &m.DeletedAt.Time
	}

	// Map relasi jika sudah di-preload
	for _, p := range m.Participants {
		minute.Participants = append(minute.Participants, p.ToDomain())
	}
	for _, r := range m.Results {
		minute.Results = append(minute.Results, r.ToDomain())
	}
	for _, i := range m.Images {
		minute.Images = append(minute.Images, i.ToDomain())
	}

	return minute
}

func (m *MeetingParticipantModel) ToDomain() domain.MeetingParticipant {
	p := domain.MeetingParticipant{
		ID:         m.ID,
		MinuteID:   m.MinuteID,
		EmployeeID: m.EmployeeID,
		CreatedAt:  m.CreatedAt,
	}
	if m.Employee != nil {
		p.Employee = m.Employee.ToDomain()
	}
	return p
}

func (m *MeetingResultModel) ToDomain() domain.MeetingResult {
	r := domain.MeetingResult{
		ID:                   m.ID,
		MinuteID:             m.MinuteID,
		EmployeeID:           m.EmployeeID,
		TargetDescription:    m.TargetDescription,
		TargetNominal:        m.TargetNominal,
		AchievementStatus:    m.AchievementStatus,
		TargetCompletionDate: m.TargetCompletionDate,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
	if m.Employee != nil {
		r.Employee = m.Employee.ToDomain()
	}
	return r
}

func (m *MeetingImageModel) ToDomain() domain.MeetingImage {
	return domain.MeetingImage{
		ID:        m.ID,
		MinuteID:  m.MinuteID,
		FileURL:   m.FileURL,
		CreatedAt: m.CreatedAt,
	}
}

// ==========================================
// MAPPERS (Domain -> Database)
// ==========================================

func FromDomainMeetingMinute(d *domain.MeetingMinute) MeetingMinuteModel {
	m := MeetingMinuteModel{
		ID:                   d.ID,
		Division:             d.Division,
		Title:                d.Title,
		MeetingDate:          d.MeetingDate,
		MeetingType:          d.MeetingType,
		Summary:              d.Summary,
		Notes:                d.Notes,
		Speaker:              d.Speaker,
		NumberOfParticipants: d.NumberOfParticipants,
		ExternalParticipants: d.ExternalParticipants,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
	}

	for _, p := range d.Participants {
		m.Participants = append(m.Participants, FromDomainMeetingParticipant(&p))
	}
	for _, r := range d.Results {
		m.Results = append(m.Results, FromDomainMeetingResult(&r))
	}
	for _, i := range d.Images {
		m.Images = append(m.Images, FromDomainMeetingImage(&i))
	}

	return m
}

func FromDomainMeetingParticipant(d *domain.MeetingParticipant) MeetingParticipantModel {
	return MeetingParticipantModel{
		ID:         d.ID,
		MinuteID:   d.MinuteID,
		EmployeeID: d.EmployeeID,
		CreatedAt:  d.CreatedAt,
	}
}

func FromDomainMeetingResult(d *domain.MeetingResult) MeetingResultModel {
	return MeetingResultModel{
		ID:                   d.ID,
		MinuteID:             d.MinuteID,
		EmployeeID:           d.EmployeeID,
		TargetDescription:    d.TargetDescription,
		TargetNominal:        d.TargetNominal,
		AchievementStatus:    d.AchievementStatus,
		TargetCompletionDate: d.TargetCompletionDate,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
	}
}

func FromDomainMeetingImage(d *domain.MeetingImage) MeetingImageModel {
	return MeetingImageModel{
		ID:        d.ID,
		MinuteID:  d.MinuteID,
		FileURL:   d.FileURL,
		CreatedAt: d.CreatedAt,
	}
}
