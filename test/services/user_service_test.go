package services_test

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/nabilwafi/warehouse-management-system/src/services"
	"github.com/nabilwafi/warehouse-management-system/src/utils"
	mocks "github.com/nabilwafi/warehouse-management-system/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := services.NewUserService(mockRepo)

	loginRequest := &dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	passwordHash, _ := utils.HashPassword(loginRequest.Password)

	user := &dto.User{
		ID:       uuid.New(),
		Name:     "Test User",
		Email:    "test@example.com",
		Password: passwordHash,
		Role:     "user",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByEmail", loginRequest.Email).Return(user, 200, nil).Once()
		token, statusCode, err := service.Login(loginRequest)

		assert.NoError(t, err)
		assert.Equal(t, 200, statusCode)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Wrong Password", func(t *testing.T) {
		wrongPasswordRequest := &dto.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}
		mockRepo.On("FindByEmail", wrongPasswordRequest.Email).Return(user, 200, nil).Once()

		token, statusCode, err := service.Login(wrongPasswordRequest)

		assert.Error(t, err)
		assert.Equal(t, 400, statusCode)
		assert.Equal(t, "", token)
		assert.Equal(t, "wrong password", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo.On("FindByEmail", loginRequest.Email).Return((*dto.User)(nil), 404, sql.ErrNoRows).Once()

		token, statusCode, err := service.Login(loginRequest)

		assert.Error(t, err)
		assert.Equal(t, 404, statusCode)
		assert.Equal(t, "", token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockRepo.On("FindByEmail", loginRequest.Email).Return((*dto.User)(nil), 500, assert.AnError).Once()

		token, statusCode, err := service.Login(loginRequest)

		assert.Error(t, err)
		assert.Equal(t, 500, statusCode)
		assert.Equal(t, "", token)
		mockRepo.AssertExpectations(t)
	})
}

func TestRegister(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := services.NewUserService(mockRepo)

	registerRequest := &dto.RegisterRequest{
		Name:     "New User",
		Email:    "newuser@example.com",
		Password: "password123",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByEmail", registerRequest.Email).Return((*dto.User)(nil), 404, sql.ErrNoRows).Once()
		mockRepo.On("Save", mock.Anything).Return(nil).Once()

		statusCode, err := service.Register(registerRequest)

		assert.NoError(t, err)
		assert.Equal(t, 201, statusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		existingUser := &dto.User{Email: registerRequest.Email}
		mockRepo.On("FindByEmail", registerRequest.Email).Return(existingUser, 200, nil).Once()

		statusCode, err := service.Register(registerRequest)

		assert.Error(t, err)
		assert.Equal(t, 401, statusCode)
		assert.Equal(t, "email is exists", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Internal Server Error on Save", func(t *testing.T) {
		mockRepo.On("FindByEmail", registerRequest.Email).Return((*dto.User)(nil), 404, sql.ErrNoRows).Once()
		mockRepo.On("Save", mock.Anything).Return(assert.AnError).Once()

		statusCode, err := service.Register(registerRequest)

		assert.Error(t, err)
		assert.Equal(t, 500, statusCode)
		assert.Equal(t, assert.AnError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := services.NewUserService(mockRepo)

	userID := uuid.New()
	user := &dto.User{
		ID:    userID,
		Name:  "Test User",
		Email: "test@example.com",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByID", userID).Return(user, 200, nil).Once()

		result, statusCode, err := service.GetUserByID(userID)

		assert.NoError(t, err)
		assert.Equal(t, 200, statusCode)
		assert.Equal(t, user, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo.On("FindByID", userID).Return((*dto.User)(nil), 404, sql.ErrNoRows).Once()

		result, statusCode, err := service.GetUserByID(userID)

		assert.Error(t, err)
		assert.Equal(t, 404, statusCode)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockRepo.On("FindByID", userID).Return((*dto.User)(nil), 500, assert.AnError).Once()

		result, statusCode, err := service.GetUserByID(userID)

		assert.Error(t, err)
		assert.Equal(t, 500, statusCode)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := services.NewUserService(mockRepo)

	pagination := &web.PaginationRequest{
		Page: 1,
		Size: 10,
	}

	users := []*dto.User{
		{
			ID:    uuid.New(),
			Name:  "User 1",
			Email: "user1@example.com",
		},
		{
			ID:    uuid.New(),
			Name:  "User 2",
			Email: "user2@example.com",
		},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetAll", pagination).Return(users, nil).Once()

		result, statusCode, err := service.GetAllUser(pagination)

		assert.NoError(t, err)
		assert.Equal(t, 200, statusCode)
		assert.Equal(t, users, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("No Users Found", func(t *testing.T) {
		mockRepo.On("GetAll", pagination).Return(([]*dto.User)(nil), nil).Once()

		result, statusCode, err := service.GetAllUser(pagination)

		assert.NoError(t, err)
		assert.Equal(t, 200, statusCode)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockRepo.On("GetAll", pagination).Return(([]*dto.User)(nil), assert.AnError).Once()

		result, statusCode, err := service.GetAllUser(pagination)

		assert.Error(t, err)
		assert.Equal(t, 500, statusCode)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
