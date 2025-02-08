package mocks

import (
	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserRepository) Save(register *dto.RegisterRequest) error {
	args := m.Called(register)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*dto.User, int, error) {
	args := m.Called(email)
	return args.Get(0).(*dto.User), args.Int(1), args.Error(2)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*dto.User, int, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.User), args.Int(1), args.Error(2)
}

func (m *MockUserRepository) GetAll(pagination *web.PaginationRequest) ([]*dto.User, error) {
	args := m.Called(pagination)
	return args.Get(0).([]*dto.User), args.Error(1)
}

func (m *MockUserService) Register(req *dto.RegisterRequest) (int, error) {
	args := m.Called(req)
	return args.Int(0), args.Error(1)
}

func (m *MockUserService) Login(req *dto.LoginRequest) (string, int, error) {
	args := m.Called(req)
	return args.String(0), args.Int(1), args.Error(2)
}

func (m *MockUserService) GetUserByID(id uuid.UUID) (*dto.User, int, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.User), args.Int(1), args.Error(2)
}

func (m *MockUserService) GetAllUser(pagination *web.PaginationRequest) ([]*dto.User, int, error) {
	args := m.Called(pagination)
	return args.Get(0).([]*dto.User), args.Int(1), args.Error(2)
}
