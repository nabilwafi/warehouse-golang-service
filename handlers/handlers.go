package handlers

import (
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	GetMe(c *gin.Context)
	ListUsers(c *gin.Context)
}

type ProductHandler interface {
	AddProduct(c *gin.Context)
	GetAllProducts(c *gin.Context)
	GetProductByID(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
}

type LocationHandler interface {
	AddLocation(c *gin.Context)
	GetAllLocations(c *gin.Context)
}

type OrderHandler interface {
	ReceiveOrder(c *gin.Context)
	ShipOrder(c *gin.Context)
	GetAllOrders(c *gin.Context)
	GetOrderByID(c *gin.Context)
}
