package dto

import "time"

// AssignTargetRequest adalah payload untuk menetapkan target ke karyawan
type AssignTargetRequest struct {
	EmployeeID string `json:"employee_id" example:"uuid-employee" validate:"required,uuid"`
	ProductID  string `json:"product_id" example:"uuid-product" validate:"required,uuid"`
	Nominal    int64  `json:"nominal" example:"500000000" validate:"required,gt=0"`
	Month      int    `json:"month" example:"5" validate:"required,min=1,max=12"`
	Year       int    `json:"year" example:"2026" validate:"required,min=2000"`
}

// TargetResponse adalah format balikan data target
type TargetResponse struct {
	ID         string    `json:"id"`
	EmployeeID string    `json:"employee_id"`
	ProductID  string    `json:"product_id"`
	Nominal    int64     `json:"nominal"`
	Month      int       `json:"month"`
	Year       int       `json:"year"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
