package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/src/handlers"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/nabilwafi/warehouse-management-system/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestLocationHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockLocationService := new(mocks.MockLocationService)
	validator := validator.New()
	handler := handlers.NewLocationHandler(mockLocationService, validator)

	t.Run("AddLocation_Success", func(t *testing.T) {
		location := dto.Location{
			ID:   uuid.New(),
			Name: "Main Warehouse",
		}

		mockLocationService.On("Save", &location).Return(201, nil).Once()

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		body, _ := json.Marshal(location)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/locations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		ctx.Request = req

		handler.AddLocation(ctx)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "CREATED")
		mockLocationService.AssertExpectations(t)
	})

	t.Run("AddLocation_BadRequest", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		body := []byte(`{"name":}`)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/locations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req

		handler.AddLocation(ctx)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "invalid character")
	})

	t.Run("GetAllLocations_Success", func(t *testing.T) {
		pagination := &web.PaginationRequest{
			Page: 1,
			Size: 10,
		}

		mockLocations := []*dto.Location{
			{ID: uuid.New(), Name: "Main Warehouse"},
			{ID: uuid.New(), Name: "Secondary Warehouse"},
		}

		mockLocationService.On("GetAll", pagination).Return(mockLocations, 200, nil).Once()

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/locations?page=1&size=10", nil)
		ctx.Request = req

		handler.GetAllLocations(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Main Warehouse")
		assert.Contains(t, recorder.Body.String(), "Secondary Warehouse")
		mockLocationService.AssertExpectations(t)
	})

	t.Run("GetAllLocations_Error", func(t *testing.T) {
		pagination := &web.PaginationRequest{
			Page: 1,
			Size: 10,
		}

		mockLocationService.On("GetAll", pagination).Return(([]*dto.Location)(nil), 500, errors.New("internal server error")).Once()

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/locations?page=1&size=10", nil)
		ctx.Request = req

		handler.GetAllLocations(ctx)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "internal server error")
		mockLocationService.AssertExpectations(t)
	})
}
