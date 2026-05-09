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
func toProductResponse(prod *domain.Product) dto.ProductResponse {
	if prod == nil {
		return dto.ProductResponse{}
	}
	return dto.ProductResponse{
		ID:         prod.ID,
		Name:       prod.Name,
		NameNorm:   prod.NameNorm,
		CategoryID: prod.CategoryID,
		CreatedAt:  prod.CreatedAt,
		UpdatedAt:  prod.UpdatedAt,
	}
}

func toProductWithCategoryResponse(prod *domain.Product) dto.ProductWithCategoryResponse {
	if prod == nil {
		return dto.ProductWithCategoryResponse{}
	}

	res := dto.ProductWithCategoryResponse{
		ProductResponse: toProductResponse(prod),
	}

	// Jika ada relasi Category, map sekalian
	if prod.Category != nil {
		res.Category = toCategoryResponse(prod.Category) // toCategoryResponse bisa di-import/dibuat public jika beda package
	}

	return res
}

// CreateProduct godoc
// @Summary Buat Produk Baru
// @Description Menambahkan data produk baru di bawah kategori tertentu
// @Tags Product
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Payload data produk"
// @Success 201 {object} utils.SuccessResponse[dto.ProductResponse] "Produk berhasil dibuat"
// @Failure 400 {object} utils.ErrorResponse "Format JSON salah atau validasi gagal"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req dto.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON salah")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	productDomain := &domain.Product{
		Name:       req.Name,
		CategoryID: req.CategoryID,
	}

	err := h.usecase.CreateProduct(productDomain)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	res := toProductResponse(productDomain)
	return utils.SendSuccess(c, fiber.StatusCreated, "Produk berhasil dibuat", res)
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

	products, err := h.usecase.GetProductsByCategory(categoryID)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	var res []dto.ProductWithCategoryResponse
	for _, prod := range products {
		res = append(res, toProductWithCategoryResponse(&prod))
	}

	if res == nil {
		res = make([]dto.ProductWithCategoryResponse, 0)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil daftar produk", res)
}
