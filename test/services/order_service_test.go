package services_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/nabilwafi/warehouse-management-system/src/services"
	mocks "github.com/nabilwafi/warehouse-management-system/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReceiveOrder_Success(t *testing.T) {
	orderRepo := new(mocks.MockOrderRepository)
	productRepo := new(mocks.MockProductRepository)
	transactionRepo := new(mocks.MockTransactionRepository)

	orderService := services.NewOrderService(orderRepo, productRepo, transactionRepo)

	orderRequest := &dto.OrderCreateRequest{
		ProductID: uuid.New(),
		Quantity:  10,
	}

	orderRepo.On("SaveWithTransaction", mock.Anything, mock.Anything).Return(nil)
	productRepo.On("IncreaseStockWithTransaction", mock.Anything, mock.Anything, mock.Anything).Return(int64(10), nil)
	transactionRepo.On("Begin").Return(nil)
	transactionRepo.On("Commit").Return(nil)
	transactionRepo.On("Transaction", mock.Anything).Return(nil)
	transactionRepo.On("GetTx").Return(&sqlx.Tx{}, nil)

	status, err := orderService.ReceiveOrder(orderRequest)

	assert.NoError(t, err)
	assert.Equal(t, 200, status)
}

func TestOrderService(t *testing.T) {
	orderRepo := new(mocks.MockOrderRepository)
	productRepo := new(mocks.MockProductRepository)
	transactionRepo := new(mocks.MockTransactionRepository)

	orderService := services.NewOrderService(orderRepo, productRepo, transactionRepo)

	orderRequest := &dto.OrderCreateRequest{
		ProductID: uuid.New(),
		Quantity:  10,
	}

	orderID := uuid.New()
	mockOrder := &dto.Order{
		ID:        orderID,
		ProductID: orderRequest.ProductID,
		Quantity:  orderRequest.Quantity,
		Type:      dto.OrderTypeReceiving,
	}

	pagination := &web.PaginationRequest{
		Page: 1,
		Size: 10,
	}

	t.Run("ReceiveOrder - Success", func(t *testing.T) {
		orderRepo.On("SaveWithTransaction", mock.Anything, mock.Anything).Return(nil).Once()
		productRepo.On("IncreaseStockWithTransaction", mock.Anything, mock.Anything, mock.Anything).Return(int64(10), nil).Once()
		transactionRepo.On("Begin").Return(nil).Once()
		transactionRepo.On("Commit").Return(nil).Once()
		transactionRepo.On("Transaction", mock.Anything).Return(nil).Once()
		transactionRepo.On("GetTx").Return(&sqlx.Tx{}, nil).Once()

		status, err := orderService.ReceiveOrder(orderRequest)

		assert.NoError(t, err)
		assert.Equal(t, 200, status)
	})

	t.Run("ShipOrder - Success", func(t *testing.T) {
		orderRepo.On("SaveWithTransaction", mock.Anything, mock.Anything).Return(nil).Once()
		productRepo.On("DecreaseStockWithTransaction", mock.Anything, mock.Anything, mock.Anything).Return(int64(5), nil).Once()
		transactionRepo.On("Begin").Return(nil).Once()
		transactionRepo.On("Commit").Return(nil).Once()
		transactionRepo.On("Transaction", mock.Anything).Return(nil).Once()
		transactionRepo.On("GetTx").Return(&sqlx.Tx{}, nil).Once()

		status, err := orderService.ShipOrder(orderRequest)

		assert.NoError(t, err)
		assert.Equal(t, 200, status)
	})

	t.Run("GetAllOrders - Success", func(t *testing.T) {
		mockOrders := []*dto.Order{mockOrder}
		orderRepo.On("FindAll", pagination).Return(mockOrders, nil).Once()

		orders, status, err := orderService.GetAllOrders(pagination)

		assert.NoError(t, err)
		assert.Equal(t, 200, status)
		assert.Equal(t, mockOrders, orders)
	})

	t.Run("GetAllOrders - Failure", func(t *testing.T) {
		orderRepo.On("FindAll", pagination).Return(([]*dto.Order)(nil), errors.New("failed to get orders")).Once()

		orders, status, err := orderService.GetAllOrders(pagination)

		assert.Error(t, err)
		assert.Equal(t, 500, status)
		assert.Nil(t, orders)
	})

	t.Run("GetOrderByID - Success", func(t *testing.T) {
		orderRepo.On("FindByID", orderID).Return(mockOrder, 200, nil).Once()

		order, status, err := orderService.GetOrderByID(orderID)

		assert.NoError(t, err)
		assert.Equal(t, 200, status)
		assert.Equal(t, mockOrder, order)
	})

	t.Run("GetOrderByID - Failure", func(t *testing.T) {
		orderRepo.On("FindByID", orderID).Return((*dto.Order)(nil), 404, errors.New("order not found")).Once()

		order, status, err := orderService.GetOrderByID(orderID)

		assert.Error(t, err)
		assert.Equal(t, 404, status)
		assert.Nil(t, order)
		assert.Equal(t, "order not found", err.Error())
	})
}
