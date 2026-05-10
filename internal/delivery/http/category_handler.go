package http

import (
	// Sesuaikan path import dengan module-mu
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
// @Param request body dto.CreateCategoryRequest true "Payload data kategori"
// @Success 201 {object} utils.SuccessResponse[dto.CategoryResponse] "Kategori berhasil dibuat"
// @Failure 400 {object} utils.ErrorResponse "Format JSON salah atau validasi gagal"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/categories [post]
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req dto.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON salah")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	categoryDomain := &domain.Category{
		Name: req.Name,
	}

	err := h.usecase.CreateCategory(categoryDomain)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	res := toCategoryResponse(categoryDomain)
	return utils.SendSuccess(c, fiber.StatusCreated, "Kategori berhasil dibuat", res)
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
	categories, err := h.usecase.GetAllCategories()
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	var res []dto.CategoryResponse
	for _, cat := range categories {
		res = append(res, toCategoryResponse(&cat))
	}

	// Memastikan array kosong tidak bernilai null di JSON (opsional tapi best practice)
	if res == nil {
		res = make([]dto.CategoryResponse, 0)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil daftar kategori", res)
}
