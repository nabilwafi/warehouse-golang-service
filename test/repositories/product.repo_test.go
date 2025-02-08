package repositories_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	mocks "github.com/nabilwafi/warehouse-management-system/test/mocks/repositories"
	"github.com/stretchr/testify/assert"
)

func TestMockProductRepositorySave_Error(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	product := &dto.Product{ID: uuid.New(), Name: "Test Product"}

	mockRepo.On("Save", product).Return(assert.AnError)

	err := mockRepo.Save(product)

	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
	mockRepo.AssertCalled(t, "Save", product)
}

func TestMockProductRepositoryFindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	product := &dto.Product{ID: uuid.New(), Name: "Test Product"}
	id := product.ID

	mockRepo.On("FindByID", id).Return(product, 200, nil)

	result, code, err := mockRepo.FindByID(id)

	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, product, result)
	mockRepo.AssertCalled(t, "FindByID", id)
}

func TestMockProductRepositoryFindByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	id := uuid.New()

	mockRepo.On("FindByID", id).Return((*dto.Product)(nil), 404, errors.New("product not found"))

	result, code, err := mockRepo.FindByID(id)

	assert.Nil(t, result)
	assert.Equal(t, 404, code)
	assert.EqualError(t, err, "product not found")
	mockRepo.AssertCalled(t, "FindByID", id)
}

func TestMockProductRepositoryGetAllProduct_Error(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	pagination := &web.PaginationRequest{Page: 1, Size: 10}

	mockRepo.On("GetAllProduct", pagination).Return(([]*dto.Product)(nil), assert.AnError)

	result, err := mockRepo.GetAllProduct(pagination)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, err, assert.AnError)
	mockRepo.AssertCalled(t, "GetAllProduct", pagination)
}

func TestMockProductRepositoryDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	id := uuid.New()

	mockRepo.On("Delete", id).Return(200, nil)

	code, err := mockRepo.Delete(id)

	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	mockRepo.AssertCalled(t, "Delete", id)
}

func TestMockProductRepositoryDelete_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	id := uuid.New()

	mockRepo.On("Delete", id).Return(404, assert.AnError)

	code, err := mockRepo.Delete(id)

	assert.Equal(t, 404, code)
	assert.Equal(t, assert.AnError, err)
	mockRepo.AssertCalled(t, "Delete", id)
}
