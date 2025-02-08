package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/src/handlers"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/nabilwafi/warehouse-management-system/src/utils"
	"github.com/nabilwafi/warehouse-management-system/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := new(mocks.MockUserService)
	handler := handlers.NewUserHandlerImpl(mockUserService)

	t.Run("Register_Success", func(t *testing.T) {
		reqBody := dto.RegisterRequest{
			Name:     "testuser",
			Password: "password123",
			Email:    "test@example.com",
			Role:     dto.UserRoleAdmin,
		}

		mockUserService.On("Register", &reqBody).Return(201, nil).Once()

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "api/v1/user/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req

		handler.Register(ctx)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "CREATED")
		mockUserService.AssertExpectations(t)
	})

	t.Run("Login_Success", func(t *testing.T) {
		reqBody := dto.LoginRequest{
			Email:    "testuser@gmail.com",
			Password: "password123",
		}
		mockToken := "mock_token"

		mockUserService.On("Login", &reqBody).Return(mockToken, 200, nil).Once()

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "api/v1/user/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req

		handler.Login(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), mockToken)
		mockUserService.AssertExpectations(t)
	})

	t.Run("GetMe_Success", func(t *testing.T) {
		userID := uuid.New()
		userResponse := &dto.User{ID: userID, Name: "testuser", Email: "test@test.com"}

		mockUserService.On("GetUserByID", userID).Return(userResponse, 200, nil).Once()

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		ctx.Set("user", &utils.CustomClaims{ID: userID})

		req, _ := http.NewRequest(http.MethodGet, "api/v1/user/me", nil)
		ctx.Request = req

		handler.GetMe(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "testuser")
		mockUserService.AssertExpectations(t)
	})

	t.Run("ListUsers_Success", func(t *testing.T) {
		pagination := &web.PaginationRequest{Page: 1, Size: 10}
		mockUsers := []*dto.User{
			{ID: uuid.New(), Name: "user1"},
			{ID: uuid.New(), Name: "user2"},
		}

		mockUserService.On("GetAllUser", pagination).Return(mockUsers, 200, nil).Once()

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		req, _ := http.NewRequest(http.MethodGet, "api/v1/user/users?page=1&size=10", nil)
		ctx.Request = req

		handler.ListUsers(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "user1")
		assert.Contains(t, recorder.Body.String(), "user2")
		mockUserService.AssertExpectations(t)
	})
}
