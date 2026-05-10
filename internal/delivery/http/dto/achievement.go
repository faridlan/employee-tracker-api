package dto

import "time"

type RecordAchievementRequest struct {
	TargetID    string    `json:"target_id" example:"uuid-target" validate:"required,uuid"`
	Nominal     int64     `json:"nominal" example:"150000000" validate:"required,gt=0"`
	Description string    `json:"description" example:"Pencairan KUR Nasabah A"`
	ClosingDate time.Time `json:"closing_date" example:"2026-05-10T14:00:00Z"`
}

type AchievementResponse struct {
	ID          string    `json:"id"`
	TargetID    string    `json:"target_id"`
	Nominal     int64     `json:"nominal"`
	Description string    `json:"description"`
	ClosingDate time.Time `json:"closing_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
