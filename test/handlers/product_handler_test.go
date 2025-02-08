package handlers_test

import (
	"bytes"
	"encoding/json"
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

func TestProductHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockProductService := new(mocks.MockProductService)
	validator := validator.New()
	handler := handlers.NewProductHandler(mockProductService, validator)

	t.Run("AddProduct_Success", func(t *testing.T) {
		product := dto.Product{
			Name:       "Product A",
			SKU:        "PROD123",
			Quantity:   5,
			LocationID: uuid.New(),
		}

		mockProductService.On("Create", &product).Return(201, nil).Once()

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		body, _ := json.Marshal(product)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req

		handler.AddProduct(ctx)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "CREATED")
		mockProductService.AssertExpectations(t)
	})

	t.Run("GetAllProducts_Success", func(t *testing.T) {
		pagination := &web.PaginationRequest{Page: 1, Size: 10}
		mockProducts := []*dto.Product{
			{ID: uuid.New(), Name: "Product A"},
			{ID: uuid.New(), Name: "Product B"},
		}

		mockProductService.On("GetAll", pagination).Return(mockProducts, 200, nil).Once()

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/products?page=1&size=10", nil)
		ctx.Request = req

		handler.GetAllProducts(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Product A")
		assert.Contains(t, recorder.Body.String(), "Product B")
		mockProductService.AssertExpectations(t)
	})

	t.Run("GetProductByID_Success", func(t *testing.T) {
		productID := uuid.New()
		product := &dto.Product{ID: productID, Name: "Product A"}

		mockProductService.On("GetByID", productID).Return(product, 200, nil).Once()

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/products/"+productID.String(), nil)
		ctx.Params = gin.Params{{Key: "product_id", Value: productID.String()}}
		ctx.Request = req

		handler.GetProductByID(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Product A")
		mockProductService.AssertExpectations(t)
	})

	t.Run("UpdateProduct_Success", func(t *testing.T) {
		productID := uuid.New()
		product := dto.Product{ID: productID, Name: "Updated Product", SKU: "SK-1000", Quantity: 6, LocationID: uuid.New()}

		mockProductService.On("Update", &product).Return(200, nil).Once()

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		body, _ := json.Marshal(product)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/products/"+productID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{{Key: "product_id", Value: productID.String()}}
		ctx.Request = req

		handler.UpdateProduct(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Successfully Updated Data")
		mockProductService.AssertExpectations(t)
	})

	t.Run("DeleteProduct_Success", func(t *testing.T) {
		productID := uuid.New()

		mockProductService.On("Delete", productID).Return(200, nil).Once()

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/products/"+productID.String(), nil)
		ctx.Params = gin.Params{{Key: "product_id", Value: productID.String()}}
		ctx.Request = req

		handler.DeleteProduct(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Successfully delete product")
		mockProductService.AssertExpectations(t)
	})
}
