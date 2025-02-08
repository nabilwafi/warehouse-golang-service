package services_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/nabilwafi/warehouse-management-system/src/services"
	mocks "github.com/nabilwafi/warehouse-management-system/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := services.NewProductService(mockRepo)

	product := &dto.Product{
		ID:       uuid.New(),
		Name:     "Product 1",
		SKU:      "SKU001",
		Quantity: 5,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByName", product.Name).Return((*dto.Product)(nil), 0, sql.ErrNoRows).Once()
		mockRepo.On("FindBySKU", product.SKU).Return((*dto.Product)(nil), 0, sql.ErrNoRows).Once()
		mockRepo.On("Save", product).Return(nil).Once()

		statusCode, err := service.Create(product)

		assert.NoError(t, err)
		assert.Equal(t, 201, statusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Product Name Exists", func(t *testing.T) {
		mockRepo.On("FindByName", product.Name).Return(product, 200, nil).Once()

		statusCode, err := service.Create(product)

		assert.Error(t, err)
		assert.Equal(t, 401, statusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Product SKU Exists", func(t *testing.T) {
		mockRepo.On("FindByName", product.Name).Return((*dto.Product)(nil), 0, sql.ErrNoRows).Once()
		mockRepo.On("FindBySKU", product.SKU).Return(product, 200, nil).Once()

		statusCode, err := service.Create(product)

		assert.Error(t, err)
		assert.Equal(t, 401, statusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockRepo.On("FindByName", product.Name).Return((*dto.Product)(nil), 500, assert.AnError).Once()

		statusCode, err := service.Create(product)

		assert.Error(t, err)
		assert.Equal(t, 500, statusCode)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetProductByID(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := services.NewProductService(mockRepo)

	productID := uuid.New()
	product := &dto.Product{
		ID:       productID,
		Name:     "Product 1",
		SKU:      "SKU001",
		Quantity: 5,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByID", productID).Return(product, 200, nil).Once()

		result, statusCode, err := service.GetByID(productID)
		assert.NoError(t, err)
		assert.Equal(t, 200, statusCode)
		assert.Equal(t, product, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Product Not Found", func(t *testing.T) {
		mockRepo.On("FindByID", productID).Return((*dto.Product)(nil), 404, errors.New("not found")).Once()

		result, statusCode, err := service.GetByID(productID)
		assert.Error(t, err)
		assert.Equal(t, 404, statusCode)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllProducts(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := services.NewProductService(mockRepo)

	pagination := &web.PaginationRequest{
		Page: 1,
		Size: 10,
	}

	products := []*dto.Product{
		{ID: uuid.New(), Name: "Product 1", SKU: "SKU001", Quantity: 5},
		{ID: uuid.New(), Name: "Product 2", SKU: "SKU002", Quantity: 15},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetAllProduct", pagination).Return(products, nil).Once()

		result, statusCode, err := service.GetAll(pagination)
		assert.NoError(t, err)
		assert.Equal(t, 200, statusCode)
		assert.Equal(t, products, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("No Products Found", func(t *testing.T) {
		mockRepo.On("GetAllProduct", pagination).Return(([]*dto.Product)(nil), nil).Once()

		result, statusCode, err := service.GetAll(pagination)

		assert.NoError(t, err)
		assert.Equal(t, 200, statusCode)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockRepo.On("GetAllProduct", pagination).Return(([]*dto.Product)(nil), assert.AnError).Once()

		result, statusCode, err := service.GetAll(pagination)
		assert.Error(t, err)
		assert.Equal(t, 500, statusCode)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateProduct(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := services.NewProductService(mockRepo)

	product := &dto.Product{
		ID:       uuid.New(),
		Name:     "Updated Product",
		SKU:      "SKU001",
		Quantity: 5,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByID", product.ID).Return(product, 200, nil).Once()
		mockRepo.On("Update", product).Return(nil).Once()

		statusCode, err := service.Update(product)
		assert.NoError(t, err)
		assert.Equal(t, 200, statusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Product Not Found", func(t *testing.T) {
		mockRepo.On("FindByID", product.ID).Return((*dto.Product)(nil), 404, errors.New("not found")).Once()

		statusCode, err := service.Update(product)
		assert.Error(t, err)
		assert.Equal(t, 404, statusCode)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteProduct(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := services.NewProductService(mockRepo)

	productID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByID", productID).Return(&dto.Product{}, 200, nil).Once()
		mockRepo.On("Delete", productID).Return(200, nil).Once()

		statusCode, err := service.Delete(productID)
		assert.NoError(t, err)
		assert.Equal(t, 200, statusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Product Not Found", func(t *testing.T) {
		mockRepo.On("FindByID", productID).Return((*dto.Product)(nil), 404, errors.New("not found")).Once()

		statusCode, err := service.Delete(productID)
		assert.Error(t, err)
		assert.Equal(t, 404, statusCode)
		mockRepo.AssertExpectations(t)
	})
}
