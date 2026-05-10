package http

import (
	"github.com/faridlan/employee-tracker-api/internal/delivery/http/dto"
	"github.com/faridlan/employee-tracker-api/internal/domain"
	"github.com/faridlan/employee-tracker-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	usecase domain.CategoryUsecase
}

func NewCategoryHandler(u domain.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{usecase: u}
}

// ==========================================
// HELPER: MAPPING ENTITY KE RESPONSE DTO
// ==========================================
func toCategoryResponse(cat *domain.Category) dto.CategoryResponse {
	if cat == nil {
		return dto.CategoryResponse{}
	}
	return dto.CategoryResponse{
		ID:        cat.ID,
		Name:      cat.Name,
		NameNorm:  cat.NameNorm,
		CreatedAt: cat.CreatedAt,
		UpdatedAt: cat.UpdatedAt,
	}
}

// CreateCategory godoc
// @Summary Buat Kategori Baru
// @Description Menambahkan data kategori produk master ke sistem
// @Tags Category
// @Accept json
// @Produce json
// @Param request body dto.CategoryRequest true "Payload data kategori"
// @Success 201 {object} utils.SuccessResponse[dto.CategoryResponse] "Kategori berhasil dibuat"
// @Failure 400 {object} utils.ErrorResponse "Format JSON salah atau validasi gagal"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/categories [post]
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req dto.CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON salah")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	categoryInput := domain.CreateCategoryInput{
		Name: req.Name,
	}

	result, err := h.usecase.CreateCategory(c.Context(), categoryInput)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	res := toCategoryResponse(result)
	return utils.SendSuccess(c, fiber.StatusCreated, "Kategori berhasil dibuat", res)
}

// UpdateCategory godoc
// @Summary Update Data Kategori
// @Description Memperbaharui data kategori produk yang ada di dalam sistem
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "ID Kategori"
// @Param request body dto.CategoryRequest true "Payload data kategori"
// @Success 200 {object} utils.SuccessResponse[dto.CategoryResponse] "Kategori berhasil diperbaharui"
// @Failure 400 {object} utils.ErrorResponse "Format JSON salah atau validasi gagal"
// @Failure 404 {object} utils.ErrorResponse "Kategori tidak ditemukan"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := utils.ValidateUUID(id, "id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req dto.CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON salah")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	categoryInput := domain.UpdateCategoryInput{
		ID:   id,
		Name: req.Name,
	}

	result, err := h.usecase.UpdateCategory(c.Context(), categoryInput)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	res := toCategoryResponse(result)
	return utils.SendSuccess(c, fiber.StatusOK, "Kategori berhasil diperbaharui", res) // Menggunakan StatusOK (200)
}

// GetCategoryDetail godoc
// @Summary Detail Kategori
// @Description Menampilkan detail informasi kategori berdasarkan ID
// @Tags Category
// @Produce json
// @Param id path string true "ID Kategori"
// @Success 200 {object} utils.SuccessResponse[dto.CategoryResponse] "Berhasil mengambil detail kategori"
// @Failure 400 {object} utils.ErrorResponse "Format UUID salah"
// @Failure 404 {object} utils.ErrorResponse "Kategori tidak ditemukan"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/categories/{id} [get]
func (h *CategoryHandler) GetCategoryDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := utils.ValidateUUID(id, "id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	result, err := h.usecase.GetCategoryByID(c.Context(), id)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	res := toCategoryResponse(result)
	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil detail kategori", res)
}

// GetAllCategories godoc
// @Summary List Semua Kategori
// @Description Mengambil daftar semua kategori produk yang tersedia
// @Tags Category
// @Produce json
// @Success 200 {object} utils.SuccessResponse[[]dto.CategoryResponse] "Berhasil mengambil daftar kategori"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/categories [get]
func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.usecase.GetAllCategories(c.Context())
	if err != nil {
		return utils.HandleDomainError(c, err) // Menggunakan HandleDomainError agar konsisten
	}

	// Mencegah return null array
	res := make([]dto.CategoryResponse, 0)
	for _, cat := range categories {
		res = append(res, toCategoryResponse(cat))
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil daftar kategori", res)
}
