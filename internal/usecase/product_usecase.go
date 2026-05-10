package usecase

import (
	"context"
	"strings"

	"github.com/faridlan/employee-tracker-api/internal/domain"
)

type productUsecase struct {
	productRepo  domain.ProductRepository
	categoryRepo domain.CategoryRepository
}

func NewProductUsecase(pRepo domain.ProductRepository, cRepo domain.CategoryRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepo:  pRepo,
		categoryRepo: cRepo,
	}
}

func (u *productUsecase) CreateProduct(ctx context.Context, input domain.CreateProductInput) (*domain.Product, error) {
	// Validasi Category Exist
	if _, err := u.categoryRepo.GetByID(ctx, input.CategoryID); err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Kategori tidak ditemukan")
	}

	norm := strings.ToLower(strings.TrimSpace(input.Name))
	nameNorm := strings.ReplaceAll(norm, " ", "-")

	product := &domain.Product{
		Name:       input.Name,
		NameNorm:   nameNorm,
		CategoryID: input.CategoryID,
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return u.GetProductByID(ctx, product.ID) // Reload untuk dpt detail category
}

func (u *productUsecase) UpdateProduct(ctx context.Context, input domain.UpdateProductInput) (*domain.Product, error) {
	existing, err := u.productRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Produk tidak ditemukan")
	}

	// Jika ganti kategori, cek eksistensi kategori baru
	if existing.CategoryID != input.CategoryID {
		if _, err := u.categoryRepo.GetByID(ctx, input.CategoryID); err != nil {
			return nil, domain.NewError(domain.ErrNotFound, "Kategori baru tidak ditemukan")
		}
	}

	norm := strings.ToLower(strings.TrimSpace(input.Name))
	existing.NameNorm = strings.ReplaceAll(norm, " ", "-")
	existing.Name = input.Name
	existing.CategoryID = input.CategoryID

	if err := u.productRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return u.GetProductByID(ctx, existing.ID)
}

func (u *productUsecase) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {

	result, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Produk tidak ditemukan")
	}

	return result, nil

}

func (u *productUsecase) GetAllProducts(ctx context.Context) ([]*domain.Product, error) {
	return u.productRepo.GetAll(ctx)
}

func (u *productUsecase) GetProductsByCategoryID(ctx context.Context, categoryID string) ([]*domain.Product, error) {
	return u.productRepo.GetByCategoryID(ctx, categoryID)
}
