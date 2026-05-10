package usecase

import (
	"errors"
	"strings"

	"github.com/faridlan/employee-tracker-api/internal/domain"
)

type productUsecase struct {
	productRepo domain.ProductRepository
}

// NewProductUsecase adalah constructor untuk menginisialisasi usecase produk
func NewProductUsecase(repo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepo: repo,
	}
}

// CreateProduct menangani logika pembuatan produk baru
func (u *productUsecase) CreateProduct(product *domain.Product) error {
	// Validasi dasar
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("nama produk tidak boleh kosong")
	}
	if strings.TrimSpace(product.CategoryID) == "" {
		return errors.New("kategori ID tidak boleh kosong")
	}

	// Business Logic: Otomatis membuat name_norm dari Name
	// Contoh: "KUR BANK" -> "kur-bank"
	norm := strings.ToLower(strings.TrimSpace(product.Name))
	product.NameNorm = strings.ReplaceAll(norm, " ", "-")

	// Teruskan ke repository untuk disimpan
	return u.productRepo.Create(product)
}

// GetProductsByCategory mengambil daftar produk berdasarkan ID kategorinya
func (u *productUsecase) GetProductsByCategory(categoryID string) ([]domain.Product, error) {
	if strings.TrimSpace(categoryID) == "" {
		return nil, errors.New("kategori ID tidak valid")
	}

	return u.productRepo.GetByCategoryID(categoryID)
}
