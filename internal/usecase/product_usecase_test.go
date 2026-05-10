package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/faridlan/employee-tracker-api/internal/domain"
	"github.com/faridlan/employee-tracker-api/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductUsecase_CreateProduct(t *testing.T) {
	mockProductRepo := new(mocks.ProductRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	usecase := NewProductUsecase(mockProductRepo, mockCategoryRepo)

	t.Run("Success", func(t *testing.T) {
		input := domain.CreateProductInput{
			Name:       "KUR Bank",
			CategoryID: "cat-1",
		}

		expectedCategory := &domain.Category{ID: "cat-1", Name: "Kredit"}
		expectedProduct := &domain.Product{ID: "prod-1", Name: "KUR Bank", NameNorm: "kur-bank", CategoryID: "cat-1", Category: expectedCategory}

		// 1. Mock validasi kategori (harus ketemu)
		mockCategoryRepo.On("GetByID", mock.Anything, input.CategoryID).Return(expectedCategory, nil).Once()
		// 2. Mock Create Product
		mockProductRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil).Once()
		// 3. Mock Reload Product (GetByID di akhir fungsi)
		mockProductRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(expectedProduct, nil).Once()

		result, err := usecase.CreateProduct(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "kur-bank", result.NameNorm)
		mockCategoryRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("Error - Category Not Found", func(t *testing.T) {
		input := domain.CreateProductInput{Name: "KUR Bank", CategoryID: "cat-unknown"}

		mockCategoryRepo.On("GetByID", mock.Anything, input.CategoryID).Return(nil, errors.New("not found")).Once()

		result, err := usecase.CreateProduct(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Kategori tidak ditemukan")
		mockCategoryRepo.AssertExpectations(t)
	})
}

func TestProductUsecase_UpdateProduct(t *testing.T) {
	mockProductRepo := new(mocks.ProductRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	usecase := NewProductUsecase(mockProductRepo, mockCategoryRepo)

	t.Run("Success - Same Category", func(t *testing.T) {
		input := domain.UpdateProductInput{ID: "prod-1", Name: "KUR Mikro", CategoryID: "cat-1"}
		existingProduct := &domain.Product{ID: "prod-1", Name: "KUR Bank", CategoryID: "cat-1"} // Kategori sama
		reloadedProduct := &domain.Product{ID: "prod-1", Name: "KUR Mikro", NameNorm: "kur-mikro", CategoryID: "cat-1"}

		// 1. Get existing product
		mockProductRepo.On("GetByID", mock.Anything, input.ID).Return(existingProduct, nil).Once()
		// Karena kategori sama, CategoryRepo.GetByID TIDAK dipanggil.
		// 2. Update product
		mockProductRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil).Once()
		// 3. Reload product
		mockProductRepo.On("GetByID", mock.Anything, input.ID).Return(reloadedProduct, nil).Once()

		result, err := usecase.UpdateProduct(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "kur-mikro", result.NameNorm)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("Success - Change Category", func(t *testing.T) {
		input := domain.UpdateProductInput{ID: "prod-1", Name: "KUR Mikro", CategoryID: "cat-2"}
		existingProduct := &domain.Product{ID: "prod-1", Name: "KUR Bank", CategoryID: "cat-1"} // Kategori beda
		reloadedProduct := &domain.Product{ID: "prod-1", Name: "KUR Mikro", CategoryID: "cat-2"}

		mockProductRepo.On("GetByID", mock.Anything, input.ID).Return(existingProduct, nil).Once()
		// Kategori beda, wajib cek eksistensi kategori baru
		mockCategoryRepo.On("GetByID", mock.Anything, input.CategoryID).Return(&domain.Category{ID: "cat-2"}, nil).Once()
		mockProductRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil).Once()
		mockProductRepo.On("GetByID", mock.Anything, input.ID).Return(reloadedProduct, nil).Once()

		result, err := usecase.UpdateProduct(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockProductRepo.AssertExpectations(t)
		mockCategoryRepo.AssertExpectations(t)
	})

	t.Run("Error - Product Not Found", func(t *testing.T) {
		input := domain.UpdateProductInput{ID: "prod-unknown", Name: "KUR", CategoryID: "cat-1"}
		mockProductRepo.On("GetByID", mock.Anything, input.ID).Return(nil, errors.New("not found")).Once()

		result, err := usecase.UpdateProduct(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Produk tidak ditemukan")
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductUsecase_GetProductByID(t *testing.T) {
	mockProductRepo := new(mocks.ProductRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	usecase := NewProductUsecase(mockProductRepo, mockCategoryRepo)

	t.Run("Success", func(t *testing.T) {
		expectedProduct := &domain.Product{ID: "prod-1", Name: "KUR Bank"}
		mockProductRepo.On("GetByID", mock.Anything, "prod-1").Return(expectedProduct, nil).Once()

		result, err := usecase.GetProductByID(context.Background(), "prod-1")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedProduct.Name, result.Name)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		mockProductRepo.On("GetByID", mock.Anything, "prod-x").Return(nil, errors.New("db error")).Once()

		result, err := usecase.GetProductByID(context.Background(), "prod-x")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Produk tidak ditemukan")
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductUsecase_GetAllProducts(t *testing.T) {
	mockProductRepo := new(mocks.ProductRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	usecase := NewProductUsecase(mockProductRepo, mockCategoryRepo)

	t.Run("Success", func(t *testing.T) {
		expectedProducts := []*domain.Product{
			{ID: "prod-1", Name: "KUR Bank"},
			{ID: "prod-2", Name: "Deposito"},
		}

		mockProductRepo.On("GetAll", mock.Anything).Return(expectedProducts, nil).Once()

		result, err := usecase.GetAllProducts(context.Background())

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductUsecase_GetProductsByCategoryID(t *testing.T) {
	mockProductRepo := new(mocks.ProductRepository)
	mockCategoryRepo := new(mocks.CategoryRepository)
	usecase := NewProductUsecase(mockProductRepo, mockCategoryRepo)

	t.Run("Success", func(t *testing.T) {
		expectedProducts := []*domain.Product{
			{ID: "prod-1", Name: "KUR Bank", CategoryID: "cat-1"},
		}

		mockProductRepo.On("GetByCategoryID", mock.Anything, "cat-1").Return(expectedProducts, nil).Once()

		result, err := usecase.GetProductsByCategoryID(context.Background(), "cat-1")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockProductRepo.AssertExpectations(t)
	})
}
