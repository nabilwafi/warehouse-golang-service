package mocks

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderRepository) SaveWithTransaction(tx *sqlx.Tx, order *dto.Order) error {
	args := m.Called(tx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) FindAll(pagination *web.PaginationRequest) ([]*dto.Order, error) {
	args := m.Called(pagination)
	return args.Get(0).([]*dto.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByID(id uuid.UUID) (*dto.Order, int, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.Order), args.Int(1), args.Error(2)
}

func (m *MockOrderService) ReceiveOrder(order *dto.OrderCreateRequest) (int, error) {
	args := m.Called(order)
	return args.Int(0), args.Error(1)
}

func (m *MockOrderService) ShipOrder(order *dto.OrderCreateRequest) (int, error) {
	args := m.Called(order)
	return args.Int(0), args.Error(1)
}

func (m *MockOrderService) GetAllOrders(pagination *web.PaginationRequest) ([]*dto.Order, int, error) {
	args := m.Called(pagination)
	return args.Get(0).([]*dto.Order), args.Int(1), args.Error(2)
}

func (m *MockOrderService) GetOrderByID(id uuid.UUID) (*dto.Order, int, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.Order), args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}
