package usecase

import (
	"context"
	"strings"

	"github.com/faridlan/employee-tracker-api/internal/domain"
)

type categoryUsecase struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryUsecase(repo domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{
		categoryRepo: repo,
	}
}

func (u *categoryUsecase) CreateCategory(ctx context.Context, input domain.CreateCategoryInput) (*domain.Category, error) {
	// Business Logic: Generate NameNorm
	norm := strings.ToLower(strings.TrimSpace(input.Name))
	nameNorm := strings.ReplaceAll(norm, " ", "-")

	categoryInput := &domain.Category{
		Name:     input.Name,
		NameNorm: nameNorm,
	}

	err := u.categoryRepo.Create(ctx, categoryInput)
	if err != nil {
		return nil, err
	}

	return categoryInput, nil
}

func (u *categoryUsecase) UpdateCategory(ctx context.Context, input domain.UpdateCategoryInput) (*domain.Category, error) {
	existing, err := u.categoryRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Kategori dengan ID tersebut tidak ditemukan")
	}

	existing.Name = input.Name

	// Update NameNorm juga jika namanya berubah
	norm := strings.ToLower(strings.TrimSpace(input.Name))
	existing.NameNorm = strings.ReplaceAll(norm, " ", "-")

	err = u.categoryRepo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (u *categoryUsecase) GetCategoryByID(ctx context.Context, id string) (*domain.Category, error) {
	category, err := u.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.NewError(domain.ErrNotFound, "Kategori dengan ID tersebut tidak ditemukan")
	}

	return category, nil
}

func (u *categoryUsecase) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	categories, err := u.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
