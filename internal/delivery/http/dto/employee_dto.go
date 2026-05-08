package dto

import "time"

// RegisterEmployeeRequest adalah payload untuk mendaftarkan karyawan baru
type RegisterEmployeeRequest struct {
	Name           string    `json:"name" example:"Budi Santoso" validate:"required"`
	Position       string    `json:"position" example:"Sales Executive" validate:"required"`
	OfficeLocation string    `json:"office_location" example:"Cabang Utama Jakarta" validate:"required"`
	EntryDate      time.Time `json:"entry_date" example:"2024-01-15T08:00:00Z" validate:"required"`
}

// EmployeeResponse adalah format balikan data karyawan
type EmployeeResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Position       string    `json:"position"`
	OfficeLocation string    `json:"office_location"`
	EntryDate      time.Time `json:"entry_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
