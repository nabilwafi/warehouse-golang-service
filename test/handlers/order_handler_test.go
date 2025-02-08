package handlers_test

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/src/handlers"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	orderService := new(mocks.MockOrderService)
	validate := validator.New()
	orderHandler := handlers.NewOrderHandler(orderService, validate)

	t.Run("ReceiveOrder - Success", func(t *testing.T) {
		requestBody := `{"product_id": "e4c2c817-0e3d-4a87-9e4f-70856d3120a8", "quantity": 10}`
		req, _ := http.NewRequest(http.MethodPost, "api/v1/orders/receive", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		orderService.On("ReceiveOrder", mock.Anything).Return(200, nil).Once()

		orderHandler.ReceiveOrder(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Successfully Receive Data")
	})

	t.Run("ReceiveOrder - Failure", func(t *testing.T) {
		requestBody := `{"product_id": "e4c2c817-0e3d-4a87-9e4f-70856d3120a8", "quantity": 10}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders/receive", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		orderService.On("ReceiveOrder", mock.Anything).Return(500, errors.New("internal server error")).Once()

		orderHandler.ReceiveOrder(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "internal server error")
	})

	t.Run("ShipOrder - Success", func(t *testing.T) {
		requestBody := `{"product_id": "e4c2c817-0e3d-4a87-9e4f-70856d3120a8", "quantity": 5}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders/ship", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		orderService.On("ShipOrder", mock.Anything).Return(200, nil).Once()

		orderHandler.ShipOrder(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Successfully Shipping Data")
	})

	t.Run("GetAllOrders - Success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/orders/orders?page=1&size=10", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		mockOrders := []*dto.Order{
			{ID: uuid.New(), ProductID: uuid.New(), Quantity: 5},
		}
		orderService.On("GetAllOrders", mock.Anything).Return(mockOrders, 200, nil).Once()

		orderHandler.GetAllOrders(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"quantity":5`)
	})

	t.Run("GetOrderByID - Success", func(t *testing.T) {
		orderID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/orders/"+orderID.String(), nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "order_id", Value: orderID.String()}}

		mockOrder := &dto.Order{
			ID:        orderID,
			ProductID: uuid.New(),
			Quantity:  5,
		}
		orderService.On("GetOrderByID", orderID).Return(mockOrder, 200, nil).Once()

		orderHandler.GetOrderByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), orderID.String())
	})

	t.Run("GetOrderByID - Failure", func(t *testing.T) {
		orderID := uuid.New()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/orders/"+orderID.String(), nil)
		w := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Params = gin.Params{{Key: "order_id", Value: orderID.String()}}

		orderService.On("GetOrderByID", orderID).Return(nil, 404, sql.ErrNoRows).Once()

		orderHandler.GetOrderByID(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "sql: no rows in result set")
	})
}
