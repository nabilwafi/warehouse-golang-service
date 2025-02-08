package mocks

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

type MockProductService struct {
	mock.Mock
}

func (m *MockProductRepository) Save(product *dto.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Update(product *dto.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) FindByID(id uuid.UUID) (*dto.Product, int, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.Product), args.Int(1), args.Error(2)
}

func (m *MockProductRepository) FindByName(name string) (*dto.Product, int, error) {
	args := m.Called(name)
	return args.Get(0).(*dto.Product), args.Int(1), args.Error(2)
}

func (m *MockProductRepository) FindBySKU(sku string) (*dto.Product, int, error) {
	args := m.Called(sku)
	return args.Get(0).(*dto.Product), args.Int(1), args.Error(2)
}

func (m *MockProductRepository) GetAllProduct(pagination *web.PaginationRequest) ([]*dto.Product, error) {
	args := m.Called(pagination)
	return args.Get(0).([]*dto.Product), args.Error(1)
}

func (m *MockProductRepository) Delete(id uuid.UUID) (int, error) {
	args := m.Called(id)
	return args.Int(0), args.Error(1)
}

func (m *MockProductRepository) IncreaseStockWithTransaction(tx *sqlx.Tx, productID uuid.UUID, quantity int64) (int, error) {
	args := m.Called(tx, productID, quantity)
	return args.Int(0), args.Error(1)
}

func (m *MockProductRepository) DecreaseStockWithTransaction(tx *sqlx.Tx, productID uuid.UUID, quantity int64) (int, error) {
	args := m.Called(tx, productID, quantity)
	return args.Int(0), args.Error(1)
}

func (m *MockProductService) Create(product *dto.Product) (int, error) {
	args := m.Called(product)
	return args.Int(0), args.Error(1)
}

func (m *MockProductService) GetAll(pagination *web.PaginationRequest) ([]*dto.Product, int, error) {
	args := m.Called(pagination)
	return args.Get(0).([]*dto.Product), args.Int(1), args.Error(2)
}

func (m *MockProductService) GetByID(id uuid.UUID) (*dto.Product, int, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.Product), args.Int(1), args.Error(2)
}

func (m *MockProductService) Update(product *dto.Product) (int, error) {
	args := m.Called(product)
	return args.Int(0), args.Error(1)
}

func (m *MockProductService) Delete(id uuid.UUID) (int, error) {
	args := m.Called(id)
	return args.Int(0), args.Error(1)
}
