package http

import (
	// Sesuaikan path import dengan module-mu
	"github.com/faridlan/employee-tracker-api/internal/delivery/http/dto"
	"github.com/faridlan/employee-tracker-api/internal/domain"
	"github.com/faridlan/employee-tracker-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type TargetHandler struct {
	usecase domain.TargetUsecase
}

func NewTargetHandler(u domain.TargetUsecase) *TargetHandler {
	return &TargetHandler{usecase: u}
}

// ==========================================
// HELPER: MAPPING ENTITY KE RESPONSE DTO
// ==========================================
func toTargetResponse(tgt *domain.Target) dto.TargetResponse {
	if tgt == nil {
		return dto.TargetResponse{}
	}
	return dto.TargetResponse{
		ID:         tgt.ID,
		EmployeeID: tgt.EmployeeID,
		ProductID:  tgt.ProductID,
		Nominal:    tgt.Nominal,
		Month:      tgt.Month,
		Year:       tgt.Year,
		CreatedAt:  tgt.CreatedAt,
		UpdatedAt:  tgt.UpdatedAt,
	}
}

// AssignTarget godoc
// @Summary Tetapkan Target Karyawan
// @Description Menetapkan target nominal produk untuk karyawan pada bulan dan tahun tertentu
// @Tags Target
// @Accept json
// @Produce json
// @Param request body dto.AssignTargetRequest true "Payload penetapan target"
// @Success 201 {object} utils.SuccessResponse[dto.TargetResponse] "Target berhasil ditetapkan"
// @Failure 400 {object} utils.ErrorResponse "Format JSON salah atau validasi gagal"
// @Failure 409 {object} utils.ErrorResponse "Target untuk produk ini pada periode tersebut sudah ada"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/targets [post]
func (h *TargetHandler) AssignTarget(c *fiber.Ctx) error {
	var req dto.AssignTargetRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON salah")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	targetDomain := &domain.Target{
		EmployeeID: req.EmployeeID,
		ProductID:  req.ProductID,
		Nominal:    req.Nominal,
		Month:      req.Month,
		Year:       req.Year,
	}

	err := h.usecase.AssignTargetToEmployee(targetDomain)
	if err != nil {
		// Asumsi error duplikasi akan ditangani oleh HandleDomainError atau balikan spesifik
		return utils.HandleDomainError(c, err)
	}

	res := toTargetResponse(targetDomain)
	return utils.SendSuccess(c, fiber.StatusCreated, "Target berhasil ditetapkan", res)
}

// GetEmployeePerformance godoc
// @Summary Kalkulasi Performa Karyawan
// @Description Melihat total target, total pencapaian, dan persentase performa karyawan pada bulan dan tahun tertentu
// @Tags Target
// @Produce json
// @Param employee_id path string true "ID Karyawan"
// @Param month query int true "Bulan"
// @Param year query int true "Tahun"
// @Success 200 {object} utils.SuccessResponse[dto.EmployeeResponse] "Berhasil kalkulasi performa"
// @Failure 400 {object} utils.ErrorResponse "Parameter tidak valid"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/employees/{employee_id}/performance [get]
func (h *TargetHandler) GetEmployeePerformance(c *fiber.Ctx) error {
	employeeID := c.Params("employee_id")
	if err := utils.ValidateUUID(employeeID, "employee_id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	month := c.QueryInt("month", 0)
	year := c.QueryInt("year", 0)

	if month < 1 || month > 12 || year < 2000 {
		return utils.SendError(c, fiber.StatusBadRequest, "Parameter month atau year tidak valid")
	}

	result, err := h.usecase.CalculateEmployeePerformance(employeeID, month, year)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil menghitung performa karyawan", result)
}
