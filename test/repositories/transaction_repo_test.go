package repositories_test

import (
	"testing"

	"github.com/jmoiron/sqlx"
	mocks "github.com/nabilwafi/warehouse-management-system/test/mocks/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionRepositoryBegin_Success(t *testing.T) {
	mockRepo := new(mocks.MockTransactionRepository)

	mockRepo.On("Begin").Return(nil)
	err := mockRepo.Begin()

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Begin")
}

func TestTransactionRepositoryBegin_Error(t *testing.T) {
	mockRepo := new(mocks.MockTransactionRepository)

	mockRepo.On("Begin").Return(assert.AnError)
	err := mockRepo.Begin()

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestTransactionRepositoryCommit_Success(t *testing.T) {
	mockRepo := new(mocks.MockTransactionRepository)

	mockRepo.On("Commit").Return(nil)
	err := mockRepo.Commit()

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Commit")
}

func TestTransactionRepositoryCommit_Error(t *testing.T) {
	mockRepo := new(mocks.MockTransactionRepository)

	mockRepo.On("Commit").Return(assert.AnError)
	err := mockRepo.Commit()

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestTransactionRepositoryRollback_Success(t *testing.T) {
	mockRepo := new(mocks.MockTransactionRepository)

	mockRepo.On("Rollback").Return(nil)
	err := mockRepo.Rollback()

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Rollback")
}

func TestTransactionRepositoryRollback_Error(t *testing.T) {
	mockRepo := new(mocks.MockTransactionRepository)

	mockRepo.On("Rollback").Return(assert.AnError)
	err := mockRepo.Rollback()

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestTransactionRepositoryTransaction_Success(t *testing.T) {
	mockRepo := new(mocks.MockTransactionRepository)
	callback := func() error { return nil }

	mockRepo.On("Transaction", mock.Anything).Return(nil)
	err := mockRepo.Transaction(callback)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Transaction", mock.Anything)
}

func TestTransactionRepositoryTransaction_Error(t *testing.T) {
	mockRepo := new(mocks.MockTransactionRepository)
	callback := func() error { return nil }

	mockRepo.On("Transaction", mock.Anything).Return(assert.AnError)
	err := mockRepo.Transaction(callback)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestTransactionRepositoryGetTx_Success(t *testing.T) {
	mockRepo := new(mocks.MockTransactionRepository)
	mockTx := &sqlx.Tx{}

	mockRepo.On("GetTx").Return(mockTx, nil)
	tx, err := mockRepo.GetTx()

	assert.NoError(t, err)
	assert.Equal(t, mockTx, tx)
	mockRepo.AssertCalled(t, "GetTx")
}

func TestTransactionRepositoryGetTx_Error(t *testing.T) {
	mockRepo := new(mocks.MockTransactionRepository)

	mockRepo.On("GetTx").Return((*sqlx.Tx)(nil), assert.AnError)
	tx, err := mockRepo.GetTx()

	assert.Error(t, err)
	assert.Nil(t, tx)
	assert.Equal(t, assert.AnError, err)
}
