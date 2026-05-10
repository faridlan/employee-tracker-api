package usecase

import (
	"errors"
	"strings"

	"github.com/faridlan/employee-tracker-api/internal/domain"
)

type categoryUsecase struct {
	categoryRepo domain.CategoryRepository
}

// NewCategoryUsecase adalah constructor untuk menginisialisasi usecase kategori
func NewCategoryUsecase(repo domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{
		categoryRepo: repo,
	}
}

// CreateCategory menangani logika penambahan kategori baru
func (u *categoryUsecase) CreateCategory(category *domain.Category) error {
	// Validasi dasar
	if strings.TrimSpace(category.Name) == "" {
		return errors.New("nama kategori tidak boleh kosong")
	}

	// Business Logic: Otomatis membuat name_norm dari Name
	// Contoh: "KREDIT MIKRO" -> "kredit-mikro"
	norm := strings.ToLower(strings.TrimSpace(category.Name))
	category.NameNorm = strings.ReplaceAll(norm, " ", "-")

	// Teruskan ke repository untuk disimpan
	return u.categoryRepo.Create(category)
}

// GetAllCategories mengambil seluruh daftar kategori
func (u *categoryUsecase) GetAllCategories() ([]domain.Category, error) {
	return u.categoryRepo.GetAll()
}
