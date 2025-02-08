package repositories_test

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	mocks "github.com/nabilwafi/warehouse-management-system/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMockOrderRepositorySaveWithTransaction_Success(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	tx := &sqlx.Tx{}
	order := &dto.Order{
		ID:        uuid.New(),
		Type:      dto.OrderTypeReceiving,
		ProductID: uuid.New(),
		Quantity:  10,
	}

	mockRepo.On("SaveWithTransaction", tx, order).Return(nil)
	err := mockRepo.SaveWithTransaction(tx, order)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "SaveWithTransaction", tx, order)
}

func TestMockOrderRepositorySaveWithTransaction_Error(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	tx := &sqlx.Tx{}
	order := &dto.Order{
		ID:        uuid.New(),
		Type:      dto.OrderTypeReceiving,
		ProductID: uuid.New(),
		Quantity:  10,
	}

	mockRepo.On("SaveWithTransaction", tx, order).Return(assert.AnError)
	err := mockRepo.SaveWithTransaction(tx, order)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	mockRepo.AssertCalled(t, "SaveWithTransaction", tx, order)
}

func TestMockOrderRepositoryFindAll_Success(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	pagination := &web.PaginationRequest{
		Page: 1,
		Size: 10,
	}

	orders := []*dto.Order{
		{ID: uuid.New(), Type: "purchase", ProductID: uuid.New(), Quantity: 10},
		{ID: uuid.New(), Type: "sale", ProductID: uuid.New(), Quantity: 5},
	}

	mockRepo.On("FindAll", pagination).Return(orders, nil)
	result, err := mockRepo.FindAll(pagination)

	assert.NoError(t, err)
	assert.Equal(t, orders, result)
	mockRepo.AssertCalled(t, "FindAll", pagination)
}

func TestMockOrderRepositoryFindAll_SuccessNil(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	pagination := &web.PaginationRequest{
		Page: 1,
		Size: 10,
	}

	mockRepo.On("FindAll", pagination).Return(([]*dto.Order)(nil), nil)
	result, err := mockRepo.FindAll(pagination)

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockRepo.AssertCalled(t, "FindAll", pagination)
}

func TestMockOrderRepositoryFindAll_Error(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	pagination := &web.PaginationRequest{
		Page: 1,
		Size: 10,
	}

	mockRepo.On("FindAll", pagination).Return(([]*dto.Order)(nil), assert.AnError)
	result, err := mockRepo.FindAll(pagination)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, assert.AnError, err)
	mockRepo.AssertCalled(t, "FindAll", pagination)
}

func TestMockOrderRepositoryFindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	id := uuid.New()
	order := &dto.Order{
		ID:        id,
		Type:      "purchase",
		ProductID: uuid.New(),
		Quantity:  10,
	}

	mockRepo.On("FindByID", id).Return(order, 200, nil)
	result, code, err := mockRepo.FindByID(id)

	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, order, result)
	mockRepo.AssertCalled(t, "FindByID", id)
}

func TestMockOrderRepositoryFindByID_Error404(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	id := uuid.New()

	mockRepo.On("FindByID", id).Return((*dto.Order)(nil), 404, sql.ErrNoRows)
	result, code, err := mockRepo.FindByID(id)

	assert.Error(t, err)
	assert.Equal(t, 404, code)
	assert.Nil(t, result)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestMockOrderRepositoryFindByID_Error500(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	id := uuid.New()

	mockRepo.On("FindByID", id).Return((*dto.Order)(nil), 500, assert.AnError)
	result, code, err := mockRepo.FindByID(id)

	assert.Error(t, err)
	assert.Equal(t, 500, code)
	assert.Equal(t, assert.AnError, err)
	assert.Nil(t, result)
}
