package dto

import "time"

type AssignTargetRequest struct {
	EmployeeID string `json:"employee_id" example:"uuid" validate:"required,uuid"`
	ProductID  string `json:"product_id" example:"uuid" validate:"required,uuid"`
	Nominal    int64  `json:"nominal" example:"500000000" validate:"required,gt=0"`
	Month      int    `json:"month" example:"5" validate:"required,min=1,max=12"`
	Year       int    `json:"year" example:"2026" validate:"required,min=2000"`
}

type UpdateTargetNominalRequest struct {
	Nominal int64 `json:"nominal" example:"600000000" validate:"required,gt=0"`
}

type TargetResponse struct {
	ID         string    `json:"id"`
	EmployeeID string    `json:"employee_id,omitempty"` // Tambahkan omitempty agar hilang jika tidak diperlukan
	ProductID  string    `json:"product_id,omitempty"`  // Tambahkan omitempty
	Nominal    int64     `json:"nominal"`
	Month      int       `json:"month"`
	Year       int       `json:"year"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type TargetDetailResponse struct {
	TargetResponse
	Product          *ProductResponse `json:"product,omitempty"`
	TotalAchievement int64            `json:"total_achievement"`
}

// ==========================================
// DTO KHUSUS UNTUK PERFORMANCE RESPONSE
// ==========================================

// TargetPerformanceDetail menyembunyikan field yang tidak perlu untuk tampilan nested
type TargetPerformanceDetail struct {
	ID               string           `json:"id"`
	Nominal          int64            `json:"nominal"`
	TotalAchievement int64            `json:"total_achievement"`
	Product          *ProductResponse `json:"product"` // Relasi produk langsung di sini
}

type PerformanceResponse struct {
	Employee         *EmployeeResponse         `json:"employee"` // Object Employee utuh
	Month            int                       `json:"month"`
	Year             int                       `json:"year"`
	TotalTarget      int64                     `json:"total_target"`
	TotalAchievement int64                     `json:"total_achievement"`
	Percentage       float64                   `json:"percentage"`
	Targets          []TargetPerformanceDetail `json:"targets"`
}
