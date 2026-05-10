package http

import (
	"github.com/faridlan/employee-tracker-api/internal/delivery/http/dto"
	"github.com/faridlan/employee-tracker-api/internal/domain"
	"github.com/faridlan/employee-tracker-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type EmployeeHandler struct {
	usecase domain.EmployeeUsecase
}

func NewEmployeeHandler(u domain.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{usecase: u}
}

// ==========================================
// HELPER: MAPPING ENTITY KE RESPONSE DTO
// ==========================================
func toEmployeeResponse(emp *domain.Employee) dto.EmployeeResponse {
	if emp == nil {
		return dto.EmployeeResponse{}
	}
	return dto.EmployeeResponse{
		ID:             emp.ID,
		Name:           emp.Name,
		Position:       emp.Position,
		OfficeLocation: emp.OfficeLocation,
		EntryDate:      emp.EntryDate,
		CreatedAt:      emp.CreatedAt,
		UpdatedAt:      emp.UpdatedAt,
	}
}

// RegisterEmployee godoc
// @Summary Registrasi Karyawan Baru
// @Description Mendaftarkan data karyawan baru ke dalam sistem
// @Tags Employee
// @Accept json
// @Produce json
// @Param request body dto.RegisterEmployeeRequest true "Payload data karyawan"
// @Success 201 {object} utils.SuccessResponse[dto.EmployeeResponse] "Karyawan berhasil didaftarkan"
// @Failure 400 {object} utils.ErrorResponse "Format JSON salah atau validasi gagal"
// @Failure 500 {object} utils.ErrorResponse "Gagal menyimpan data karyawan"
// @Router /api/employees [post]
func (h *EmployeeHandler) RegisterEmployee(c *fiber.Ctx) error {
	var req dto.RegisterEmployeeRequest

	// Parsing Request Body
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON salah")
	}

	// Validasi DTO menggunakan Validator (Logic validasi berpindah ke sini)
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Mapping DTO ke Domain Entity
	employeeDomain := &domain.Employee{
		Name:           req.Name,
		Position:       req.Position,
		OfficeLocation: req.OfficeLocation,
		EntryDate:      req.EntryDate,
	}

	// Panggil Usecase
	err := h.usecase.RegisterEmployee(employeeDomain)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	// Mapping balikan Domain Entity ke Response DTO
	res := toEmployeeResponse(employeeDomain)

	return utils.SendSuccess(c, fiber.StatusCreated, "Karyawan berhasil didaftarkan", res)
}

// GetEmployeeDetail godoc
// @Summary Detail Karyawan
// @Description Menampilkan detail informasi karyawan berdasarkan ID
// @Tags Employee
// @Produce json
// @Param id path string true "ID Karyawan"
// @Success 200 {object} utils.SuccessResponse[dto.EmployeeResponse] "Berhasil mengambil detail karyawan"
// @Failure 400 {object} utils.ErrorResponse "Format UUID salah"
// @Failure 404 {object} utils.ErrorResponse "Karyawan tidak ditemukan"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/employees/{id} [get]
func (h *EmployeeHandler) GetEmployeeDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validasi UUID
	if err := utils.ValidateUUID(id, "id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Panggil Usecase
	result, err := h.usecase.GetEmployeeDetails(id)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}
	if result == nil {
		return utils.SendError(c, fiber.StatusNotFound, "Karyawan tidak ditemukan")
	}

	// Mapping Domain Entity ke Response DTO
	res := toEmployeeResponse(result)

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil detail karyawan", res)
}
