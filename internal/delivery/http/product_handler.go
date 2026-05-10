package http

import (
	// Sesuaikan path import dengan module-mu
	"github.com/faridlan/employee-tracker-api/internal/delivery/http/dto"
	"github.com/faridlan/employee-tracker-api/internal/domain"
	"github.com/faridlan/employee-tracker-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	usecase domain.ProductUsecase
}

func NewProductHandler(u domain.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: u}
}

// ==========================================
// HELPER: MAPPING ENTITY KE RESPONSE DTO
// ==========================================
func toProductWithCategoryResponse(p *domain.Product) dto.ProductWithCategoryResponse {
	if p == nil {
		return dto.ProductWithCategoryResponse{}
	}
	res := dto.ProductWithCategoryResponse{
		ProductResponse: dto.ProductResponse{
			ID:         p.ID,
			Name:       p.Name,
			NameNorm:   p.NameNorm,
			CategoryID: p.CategoryID,
			CreatedAt:  p.CreatedAt,
			UpdatedAt:  p.UpdatedAt,
		},
	}
	if p.Category != nil {
		res.Category = dto.CategoryResponse{
			ID:        p.Category.ID,
			Name:      p.Category.Name,
			NameNorm:  p.Category.NameNorm,
			CreatedAt: p.Category.CreatedAt,
			UpdatedAt: p.Category.UpdatedAt,
		}
	}
	return res
}

// CreateProduct godoc
// @Summary Buat Produk Baru
// @Description Menambahkan data produk baru di bawah kategori tertentu
// @Tags Product
// @Accept json
// @Produce json
// @Param request body dto.ProductRequest true "Payload data produk"
// @Success 201 {object} utils.SuccessResponse[dto.ProductWithCategoryResponse] "Produk berhasil dibuat"
// @Failure 400 {object} utils.ErrorResponse "Format JSON salah atau validasi gagal"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req dto.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON salah")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	input := domain.CreateProductInput{Name: req.Name, CategoryID: req.CategoryID}
	result, err := h.usecase.CreateProduct(c.Context(), input)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "Produk berhasil dibuat", toProductWithCategoryResponse(result))
}

// UpdateProduct godoc
// @Summary Update Data Produk
// @Description Memperbaharui data produk yang ada di dalam sistem
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "ID Produk"
// @Param request body dto.ProductRequest true "Payload data produk"
// @Success 200 {object} utils.SuccessResponse[dto.ProductWithCategoryResponse] "Produk berhasil diperbaharui"
// @Failure 400 {object} utils.ErrorResponse "Format JSON salah atau validasi gagal"
// @Failure 404 {object} utils.ErrorResponse "Produk tidak ditemukan"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := utils.ValidateUUID(id, "id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req dto.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON salah")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	input := domain.UpdateProductInput{ID: id, Name: req.Name, CategoryID: req.CategoryID}
	result, err := h.usecase.UpdateProduct(c.Context(), input)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Produk berhasil diperbaharui", toProductWithCategoryResponse(result))
}

// GetAllProducts godoc
// @Summary List Semua Produk
// @Description Mengambil daftar semua produk beserta detail kategorinya
// @Tags Product
// @Produce json
// @Success 200 {object} utils.SuccessResponse[[]dto.ProductWithCategoryResponse] "Berhasil mengambil data produk"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/products [get]
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	results, err := h.usecase.GetAllProducts(c.Context())
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	res := make([]dto.ProductWithCategoryResponse, 0)
	for _, p := range results {
		res = append(res, toProductWithCategoryResponse(p))
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil data produk", res)
}

// GetProductByID godoc
// @Summary Detail Produk
// @Description Menampilkan detail informasi produk berdasarkan ID
// @Tags Product
// @Produce json
// @Param id path string true "ID Produk"
// @Success 200 {object} utils.SuccessResponse[dto.ProductWithCategoryResponse] "Berhasil mengambil detail produk"
// @Failure 400 {object} utils.ErrorResponse "Format UUID salah"
// @Failure 404 {object} utils.ErrorResponse "Produk tidak ditemukan"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/products/{id} [get]
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := utils.ValidateUUID(id, "id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	result, err := h.usecase.GetProductByID(c.Context(), id)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil detail produk", toProductWithCategoryResponse(result))
}

// GetProductsByCategory godoc
// @Summary List Produk Berdasarkan Kategori
// @Description Mengambil daftar produk berdasarkan ID Kategorinya
// @Tags Product
// @Produce json
// @Param category_id path string true "ID Kategori"
// @Success 200 {object} utils.SuccessResponse[[]dto.ProductWithCategoryResponse] "Berhasil mengambil produk"
// @Failure 400 {object} utils.ErrorResponse "Format UUID salah"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/categories/{category_id}/products [get]
func (h *ProductHandler) GetProductsByCategory(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")

	// Validasi UUID untuk path param
	if err := utils.ValidateUUID(categoryID, "category_id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	products, err := h.usecase.GetProductsByCategoryID(c.Context(), categoryID)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	res := make([]dto.ProductWithCategoryResponse, 0)
	for _, prod := range products {
		res = append(res, toProductWithCategoryResponse(prod))
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil daftar produk berdasarkan kategori", res)
}
