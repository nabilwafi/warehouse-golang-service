package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/src/helpers"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/nabilwafi/warehouse-management-system/src/services"
	"github.com/nabilwafi/warehouse-management-system/src/utils"
)

type OrderHandler interface {
	ReceiveOrder(c *gin.Context)
	ShipOrder(c *gin.Context)
	GetAllOrders(c *gin.Context)
	GetOrderByID(c *gin.Context)
}

type OrderHandlerImpl struct {
	order      services.OrderService
	validation *validator.Validate
}

func NewOrderHandler(order services.OrderService, validation *validator.Validate) OrderHandler {
	return &OrderHandlerImpl{order: order, validation: validation}
}

func (h *OrderHandlerImpl) ReceiveOrder(c *gin.Context) {
	var order dto.OrderCreateRequest
	if err := c.ShouldBindJSON(&order); err != nil {
		helpers.BadRequestError(c, err.Error())
		return
	}

	code, err := h.order.ReceiveOrder(&order)
	if err != nil {
		helpers.ErrorByCode(c, code, err.Error())
		return
	}

	helpers.SuccessByCode(c, code, "Successfully Receive Data")
}

func (h *OrderHandlerImpl) ShipOrder(c *gin.Context) {
	var order dto.OrderCreateRequest
	if err := c.ShouldBindJSON(&order); err != nil {
		helpers.BadRequestError(c, err.Error())
		return
	}

	code, err := h.order.ShipOrder(&order)
	if err != nil {
		helpers.ErrorByCode(c, code, err.Error())
		return
	}

	helpers.SuccessByCode(c, code, "Successfully Shipping Data")
}

func (h *OrderHandlerImpl) GetAllOrders(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		size = 10
	}

	pagination := &web.PaginationRequest{
		Page: page,
		Size: size,
	}

	if status, msg := utils.Validate(pagination); msg != "" {
		helpers.ErrorByCode(c, status, msg)
		return
	}

	users, code, err := h.order.GetAllOrders(pagination)
	if err != nil {
		helpers.ErrorByCode(c, code, err.Error())
		return
	}

	helpers.OK(c, users)
}

func (h *OrderHandlerImpl) GetOrderByID(c *gin.Context) {
	orderID := c.Param("order_id")

	orderIDConv, err := uuid.Parse(orderID)
	if err != nil {
		helpers.BadRequestError(c, "not uuid")
		return
	}

	users, code, err := h.order.GetOrderByID(orderIDConv)
	if err != nil {
		helpers.ErrorByCode(c, code, err.Error())
		return
	}

	helpers.SuccessByCode(c, code, users)
}
