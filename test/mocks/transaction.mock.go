package mocks

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Begin() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransactionRepository) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransactionRepository) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransactionRepository) Transaction(callback func() error) error {
	args := m.Called(callback)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetTx() (*sqlx.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sqlx.Tx), args.Error(1)
}
