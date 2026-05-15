package http

import (
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
func toTargetResponse(t *domain.Target) dto.TargetResponse {
	if t == nil {
		return dto.TargetResponse{}
	}
	return dto.TargetResponse{
		ID:         t.ID,
		EmployeeID: t.EmployeeID,
		ProductID:  t.ProductID,
		Nominal:    t.Nominal,
		Month:      t.Month,
		Year:       t.Year,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

func toTargetDetailResponse(t *domain.Target) dto.TargetDetailResponse {
	if t == nil {
		return dto.TargetDetailResponse{}
	}

	res := dto.TargetDetailResponse{
		TargetResponse: toTargetResponse(t),
	}

	// TAMBAHAN: Mapping Employee agar Frontend tahu milik siapa target ini
	if t.Employee != nil {
		emp := dto.EmployeeResponse{
			ID:             t.Employee.ID,
			Name:           t.Employee.Name,
			Position:       t.Employee.Position,
			OfficeLocation: t.Employee.OfficeLocation,
			EntryDate:      t.Employee.EntryDate,
		}
		res.Employee = &emp
	}

	if t.Product != nil {
		prod := dto.ProductResponse{
			ID:         t.Product.ID,
			Name:       t.Product.Name,
			NameNorm:   t.Product.NameNorm,
			CategoryID: t.Product.CategoryID,
			CreatedAt:  t.Product.CreatedAt,
			UpdatedAt:  t.Product.UpdatedAt,
		}
		res.Product = &prod
	}

	var totalAch int64 = 0
	for _, ach := range t.Achievements {
		totalAch += ach.Nominal
	}
	res.TotalAchievement = totalAch

	return res
}

func toPerformanceResponse(p *domain.EmployeePerformance) dto.PerformanceResponse {
	if p == nil {
		return dto.PerformanceResponse{}
	}

	res := dto.PerformanceResponse{
		Month:            p.Month,
		Year:             p.Year,
		TotalTarget:      p.TotalTarget,
		TotalAchievement: p.TotalAchievement,
		Percentage:       p.Percentage,
		Targets:          make([]dto.TargetPerformanceDetail, 0),
	}

	// 1. Ekstrak data Employee dari Preload (mengambil dari index 0 karena semua target milik orang yang sama)
	if len(p.Targets) > 0 && p.Targets[0].Employee != nil {
		emp := p.Targets[0].Employee
		res.Employee = &dto.EmployeeResponse{
			ID:             emp.ID,
			Name:           emp.Name,
			Position:       emp.Position,
			OfficeLocation: emp.OfficeLocation,
			EntryDate:      emp.EntryDate,
			CreatedAt:      emp.CreatedAt,
			UpdatedAt:      emp.UpdatedAt,
		}
	} else {
		// Fallback jika tidak ada target di bulan tersebut
		res.Employee = &dto.EmployeeResponse{
			ID: p.EmployeeID,
		}
	}

	// 2. Mapping list target dengan format clean
	for _, t := range p.Targets {
		detail := dto.TargetPerformanceDetail{
			ID:      t.ID,
			Nominal: t.Nominal,
		}

		// Kalkulasi pencapaian per target
		var achPerTarget int64 = 0
		for _, a := range t.Achievements {
			achPerTarget += a.Nominal
		}
		detail.TotalAchievement = achPerTarget

		// Sematkan data produk (tanpa redundansi ProductID)
		if t.Product != nil {
			detail.Product = &dto.ProductResponse{
				ID:         t.Product.ID,
				Name:       t.Product.Name,
				NameNorm:   t.Product.NameNorm,
				CategoryID: t.Product.CategoryID,
				CreatedAt:  t.Product.CreatedAt,
				UpdatedAt:  t.Product.UpdatedAt,
			}
		}

		res.Targets = append(res.Targets, detail)
	}

	return res
}

// AssignTarget godoc
// @Summary Tetapkan Target Karyawan
// @Description Menetapkan target nominal produk untuk karyawan pada bulan dan tahun tertentu
// @Tags Target
// @Accept json
// @Produce json
// @Param request body dto.AssignTargetRequest true "Payload penetapan target"
// @Success 201 {object} utils.SuccessResponse[dto.TargetDetailResponse] "Target berhasil ditetapkan"
// @Failure 400 {object} utils.ErrorResponse "Format JSON salah atau validasi gagal"
// @Failure 404 {object} utils.ErrorResponse "Karyawan atau Produk tidak ditemukan"
// @Failure 409 {object} utils.ErrorResponse "Target sudah ada (Conflict)"
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

	input := domain.AssignTargetInput{
		EmployeeID: req.EmployeeID,
		ProductID:  req.ProductID,
		Nominal:    req.Nominal,
		Month:      req.Month,
		Year:       req.Year,
	}

	result, err := h.usecase.AssignTargetToEmployee(c.Context(), input)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "Target berhasil ditetapkan", toTargetDetailResponse(result))
}

// GetEmployeePerformance godoc
// @Summary Kalkulasi Performa Karyawan
// @Description Melihat total target, total pencapaian, dan persentase performa karyawan pada bulan dan tahun tertentu
// @Tags Target
// @Produce json
// @Param employee_id path string true "ID Karyawan"
// @Param page query int false "Nomor Halaman (Default: 1)"
// @Param limit query int false "Jumlah Data per Halaman (Default: 10, Gunakan 0 untuk tanpa limit/hati-hati)"
// @Param month query int false "Bulan (1-12) - Opsional"
// @Param year query int false "Tahun (Misal: 2026) - Opsional"
// @Param product_id query string false "ID Produk - Opsional"
// @Success 200 {object} utils.SuccessResponse[dto.PerformanceResponse] "Berhasil menghitung performa"
// @Failure 400 {object} utils.ErrorResponse "Parameter tidak valid"
// @Failure 404 {object} utils.ErrorResponse "Karyawan tidak ditemukan"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/employees/{employee_id}/performance [get]
func (h *TargetHandler) GetEmployeePerformance(c *fiber.Ctx) error {
	employeeID := c.Params("employee_id")
	if err := utils.ValidateUUID(employeeID, "employee_id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10) // Secara default kita batasi 10 data per request
	month := c.QueryInt("month", 0)
	year := c.QueryInt("year", 0)
	productID := c.Query("product_id")

	// 2. Validasi & Kalkulasi
	if month < 0 || month > 12 || year < 0 {
		return utils.SendError(c, fiber.StatusBadRequest, "Parameter month atau year tidak valid")
	}
	if productID != "" {
		if err := utils.ValidateUUID(productID, "product_id"); err != nil {
			return utils.SendError(c, fiber.StatusBadRequest, err.Error())
		}
	}
	if page < 1 {
		page = 1
	}

	// Rumus mencari Offset: (Page - 1) * Limit
	offset := (page - 1) * limit

	// 3. Bungkus ke dalam Struct Filter
	filter := domain.TargetFilter{
		Month:     month,
		Year:      year,
		ProductID: productID,
		Limit:     limit,
		Offset:    offset,
	}

	result, err := h.usecase.CalculateEmployeePerformance(c.Context(), employeeID, filter)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil menghitung performa karyawan", toPerformanceResponse(result))
}

// GetAllTargets godoc
// @Summary List Semua Target (Dengan Pagination & Filter)
// @Description Mengambil daftar semua target dengan filter opsional dan pagination
// @Tags Target
// @Produce json
// @Param page query int false "Nomor Halaman (Default: 1)"
// @Param limit query int false "Jumlah Data per Halaman (Default: 10, Gunakan 0 untuk tanpa limit/hati-hati)"
// @Param month query int false "Bulan (1-12) - Opsional"
// @Param year query int false "Tahun (Misal: 2026) - Opsional"
// @Param product_id query string false "ID Produk - Opsional"
// @Success 200 {object} utils.SuccessResponse[[]dto.TargetDetailResponse] "Berhasil mengambil data target"
// @Failure 400 {object} utils.ErrorResponse "Parameter tidak valid"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/targets [get]
func (h *TargetHandler) GetAllTargets(c *fiber.Ctx) error {
	// 1. Ambil Query Params
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10) // Secara default kita batasi 10 data per request
	month := c.QueryInt("month", 0)
	year := c.QueryInt("year", 0)
	productID := c.Query("product_id")

	// 2. Validasi & Kalkulasi
	if month < 0 || month > 12 || year < 0 {
		return utils.SendError(c, fiber.StatusBadRequest, "Parameter month atau year tidak valid")
	}
	if productID != "" {
		if err := utils.ValidateUUID(productID, "product_id"); err != nil {
			return utils.SendError(c, fiber.StatusBadRequest, err.Error())
		}
	}
	if page < 1 {
		page = 1
	}

	// Rumus mencari Offset: (Page - 1) * Limit
	offset := (page - 1) * limit

	// 3. Bungkus ke dalam Struct Filter
	filter := domain.TargetFilter{
		Month:     month,
		Year:      year,
		ProductID: productID,
		Limit:     limit,
		Offset:    offset,
	}

	// 4. Panggil Usecase
	results, err := h.usecase.GetAllTargets(c.Context(), filter)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	res := make([]dto.TargetDetailResponse, 0)
	for _, t := range results {
		res = append(res, toTargetDetailResponse(t))
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil semua data target", res)
}

// UpdateTargetNominal godoc
// @Summary Update Nominal Target
// @Description Mengubah hanya nominal dari sebuah target yang sudah ditetapkan
// @Tags Target
// @Accept json
// @Produce json
// @Param id path string true "ID Target"
// @Param request body dto.UpdateTargetNominalRequest true "Payload update nominal"
// @Success 200 {object} utils.SuccessResponse[dto.TargetDetailResponse] "Nominal target berhasil diperbaharui"
// @Failure 400 {object} utils.ErrorResponse "Format JSON salah atau validasi gagal"
// @Failure 404 {object} utils.ErrorResponse "Target tidak ditemukan"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/targets/{id}/nominal [patch]
func (h *TargetHandler) UpdateTargetNominal(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := utils.ValidateUUID(id, "id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req dto.UpdateTargetNominalRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON salah")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	input := domain.UpdateTargetNominalInput{
		ID:      id,
		Nominal: req.Nominal,
	}

	result, err := h.usecase.UpdateTargetNominal(c.Context(), input)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Nominal target berhasil diperbaharui", toTargetDetailResponse(result))
}

// DeleteTarget godoc
// @Summary Hapus Target
// @Description Menghapus (Soft Delete) target karyawan berdasarkan ID
// @Tags Target
// @Produce json
// @Param id path string true "ID Target"
// @Success 200 {object} utils.SuccessResponse[utils.EmptyObj] "Target berhasil dihapus"
// @Failure 400 {object} utils.ErrorResponse "Format UUID salah"
// @Failure 404 {object} utils.ErrorResponse "Target tidak ditemukan"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/targets/{id} [delete]
func (h *TargetHandler) DeleteTarget(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := utils.ValidateUUID(id, "id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	err := h.usecase.DeleteTarget(c.Context(), id)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Target berhasil dihapus", utils.EmptyObj{})
}
