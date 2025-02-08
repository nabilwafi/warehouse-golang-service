package mocks

import (
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/stretchr/testify/mock"
)

type MockLocationRepository struct {
	mock.Mock
}

type MockLocationService struct {
	mock.Mock
}

func (m *MockLocationRepository) Save(location *dto.Location) error {
	args := m.Called(location)
	return args.Error(0)
}

func (m *MockLocationRepository) FindByName(name string) (*dto.Location, int, error) {
	args := m.Called(name)
	return args.Get(0).(*dto.Location), args.Int(1), args.Error(2)
}

func (m *MockLocationRepository) GetAllLocation(pagination *web.PaginationRequest) ([]*dto.Location, error) {
	args := m.Called(pagination)
	return args.Get(0).([]*dto.Location), args.Error(1)
}

func (m *MockLocationService) Save(location *dto.Location) (int, error) {
	args := m.Called(location)
	return args.Int(0), args.Error(1)
}

func (m *MockLocationService) GetAll(pagination *web.PaginationRequest) ([]*dto.Location, int, error) {
	args := m.Called(pagination)
	return args.Get(0).([]*dto.Location), args.Int(1), args.Error(2)
}
