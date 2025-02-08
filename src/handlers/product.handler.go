package handlers

import (
	"fmt"
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

type ProductHandler interface {
	AddProduct(c *gin.Context)
	GetAllProducts(c *gin.Context)
	GetProductByID(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
}

type ProductHandlerImpl struct {
	product    services.ProductService
	validation *validator.Validate
}

func NewProductHandler(product services.ProductService, validation *validator.Validate) ProductHandler {
	return &ProductHandlerImpl{
		product:    product,
		validation: validation,
	}
}

func (h *ProductHandlerImpl) AddProduct(c *gin.Context) {
	var product dto.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		helpers.BadRequestError(c, err.Error())
		return
	}

	if status, msg := utils.Validate(&product); msg != "" {
		helpers.ErrorByCode(c, status, msg)
		return
	}

	code, err := h.product.Create(&product)
	if err != nil {
		helpers.ErrorByCode(c, code, err.Error())
		return
	}

	helpers.Created(c, "Successfuly Created Data")
}

func (h *ProductHandlerImpl) GetAllProducts(c *gin.Context) {
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

	users, code, err := h.product.GetAll(pagination)
	if err != nil {
		helpers.ErrorByCode(c, code, err.Error())
		return
	}

	helpers.OKWithMetadata(c, users, helpers.Metadata{
		Page: page,
		Size: size,
	})
}

func (h *ProductHandlerImpl) GetProductByID(c *gin.Context) {
	productID := c.Param("product_id")

	productIDConv, err := uuid.Parse(productID)
	if err != nil {
		helpers.BadRequestError(c, "not uuid")
		return
	}

	users, code, err := h.product.GetByID(productIDConv)
	if err != nil {
		helpers.ErrorByCode(c, code, err.Error())
		return
	}

	helpers.SuccessByCode(c, code, users)
}

func (h *ProductHandlerImpl) UpdateProduct(c *gin.Context) {
	productID := c.Param("product_id")

	fmt.Println(productID)

	productDConv, err := uuid.Parse(productID)
	if err != nil {
		helpers.BadRequestError(c, "not uuid")
		return
	}

	var product dto.Product
	product.ID = productDConv

	if err := c.ShouldBindJSON(&product); err != nil {
		helpers.BadRequestError(c, err.Error())
	}

	code, err := h.product.Update(&product)
	if err != nil {
		helpers.ErrorByCode(c, code, err.Error())
		return
	}

	helpers.OK(c, "Successfully Updated Data")
}

func (h *ProductHandlerImpl) DeleteProduct(c *gin.Context) {
	productID := c.Param("product_id")

	productIDConv, err := uuid.Parse(productID)
	if err != nil {
		helpers.BadRequestError(c, "not uuid")
		return
	}

	code, err := h.product.Delete(productIDConv)
	if err != nil {
		helpers.ErrorByCode(c, code, err.Error())
		return
	}

	helpers.OK(c, "Successfully delete product")
}
