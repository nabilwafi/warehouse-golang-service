package repositories_test

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	mocks "github.com/nabilwafi/warehouse-management-system/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFindByEmail_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	email := "test@example.com"

	expectedUser := &dto.User{
		ID:    uuid.New(),
		Email: email,
		Name:  "Test User",
		Role:  "admin",
	}

	mockRepo.On("FindByEmail", email).Return(expectedUser, 200, nil)

	user, code, err := mockRepo.FindByEmail(email)

	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertCalled(t, "FindByEmail", email)
}

func TestFindByEmail_Error404(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	email := "test@example.com"

	mockRepo.On("FindByEmail", email).Return((*dto.User)(nil), 404, sql.ErrNoRows)

	_, code, err := mockRepo.FindByEmail(email)

	assert.Error(t, err)
	assert.Equal(t, 404, code)
	assert.Equal(t, sql.ErrNoRows, err)
	mockRepo.AssertCalled(t, "FindByEmail", email)
}

func TestFindByEmail_Error500(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	email := "test@example.com"

	mockRepo.On("FindByEmail", email).Return((*dto.User)(nil), 500, assert.AnError)

	_, code, err := mockRepo.FindByEmail(email)

	assert.Error(t, err)
	assert.Equal(t, 500, code)
	assert.Equal(t, assert.AnError, err)
	mockRepo.AssertCalled(t, "FindByEmail", email)
}

func TestFindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	id := uuid.New()

	expectedUser := &dto.User{
		ID:    id,
		Email: "test@gmail.com",
		Name:  "Test User",
		Role:  "admin",
	}

	mockRepo.On("FindByID", id).Return(expectedUser, 200, nil)

	user, code, err := mockRepo.FindByID(id)

	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertCalled(t, "FindByID", id)
}

func TestFindByID_Error500(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	id := uuid.New()

	mockRepo.On("FindByID", id).Return((*dto.User)(nil), 500, assert.AnError)

	_, code, err := mockRepo.FindByID(id)

	assert.Error(t, err)
	assert.Equal(t, 500, code)
	assert.Equal(t, assert.AnError, err)
	mockRepo.AssertCalled(t, "FindByID", id)
}

func TestFindByID_Error404(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	id := uuid.New()

	mockRepo.On("FindByID", id).Return((*dto.User)(nil), 404, sql.ErrNoRows)

	_, code, err := mockRepo.FindByID(id)

	assert.Error(t, err)
	assert.Equal(t, 404, code)
	assert.Equal(t, sql.ErrNoRows, err)
	mockRepo.AssertCalled(t, "FindByID", id)
}

func TestSave_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	registerRequest := &dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "securepassword",
		Name:     "Test User",
		Role:     "user",
	}

	mockRepo.On("Save", registerRequest).Return(nil)

	err := mockRepo.Save(registerRequest)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Save", registerRequest)
}

func TestSave_Error(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	registerRequest := &dto.RegisterRequest{
		Email:    "error@example.com",
		Password: "securepassword",
		Name:     "Error User",
		Role:     "user",
	}

	mockRepo.On("Save", registerRequest).Return(assert.AnError)

	err := mockRepo.Save(registerRequest)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	mockRepo.AssertCalled(t, "Save", registerRequest)
}

func TestGetAll_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	pagination := &web.PaginationRequest{
		Page: 1,
		Size: 10,
	}

	expectedUsers := []*dto.User{
		{ID: uuid.New(), Email: "user1@example.com", Name: "User One", Role: "user"},
		{ID: uuid.New(), Email: "user2@example.com", Name: "User Two", Role: "admin"},
	}

	mockRepo.On("GetAll", pagination).Return(expectedUsers, nil)

	users, err := mockRepo.GetAll(pagination)

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertCalled(t, "GetAll", pagination)
}

func TestGetAll_SuccessNil(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	pagination := &web.PaginationRequest{
		Page: 1,
		Size: 10,
	}

	mockRepo.On("GetAll", pagination).Return(([]*dto.User)(nil), nil)

	users, err := mockRepo.GetAll(pagination)

	assert.NoError(t, err)
	assert.Len(t, users, 0)
	assert.Nil(t, users)
	mockRepo.AssertCalled(t, "GetAll", pagination)
}

func TestGetAll_Nil(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	pagination := &web.PaginationRequest{
		Page: 1,
		Size: 10,
	}

	mockRepo.On("GetAll", pagination).Return(([]*dto.User)(nil), nil)

	users, err := mockRepo.GetAll(pagination)

	assert.NoError(t, err)
	assert.Len(t, users, 0)
	mockRepo.AssertCalled(t, "GetAll", pagination)
}
