package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/exception"
	"github.com/nabilwafi/warehouse-management-system/models/domain"
	"github.com/nabilwafi/warehouse-management-system/models/web"
	"github.com/nabilwafi/warehouse-management-system/services"
)

type ProductHandlerImpl struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) ProductHandler {
	return &ProductHandlerImpl{service: service}
}

func (h *ProductHandlerImpl) AddProduct(c *gin.Context) {
	var product web.ProductCreateRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		c.Error(exception.NewCustomError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.service.Create(c.Request.Context(), &product); err != nil {
		c.Error(exception.NewCustomError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandlerImpl) GetAllProducts(c *gin.Context) {
	products, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.Error(exception.NewCustomError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandlerImpl) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		c.Error(exception.NewCustomError(http.StatusBadRequest, "invalid product ID"))
		return
	}

	fmt.Println(productID)

	product, err := h.service.GetByID(c.Request.Context(), productID)
	if err != nil {
		c.Error(exception.NewCustomError(http.StatusNotFound, err.Error()))
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandlerImpl) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		c.Error(exception.NewCustomError(http.StatusBadRequest, "invalid product ID"))
		return
	}

	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.Error(exception.NewCustomError(http.StatusBadRequest, err.Error()))
		return
	}

	product.ID = productID

	if err := h.service.Update(c.Request.Context(), &product); err != nil {
		c.Error(exception.NewCustomError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandlerImpl) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		c.Error(exception.NewCustomError(http.StatusBadRequest, "invalid product ID"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), productID); err != nil {
		c.Error(exception.NewCustomError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, nil)
}
