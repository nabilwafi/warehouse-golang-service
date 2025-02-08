package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nabilwafi/warehouse-management-system/src/helpers"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/nabilwafi/warehouse-management-system/src/services"
	"github.com/nabilwafi/warehouse-management-system/src/utils"
)

type LocationHandler interface {
	AddLocation(c *gin.Context)
	GetAllLocations(c *gin.Context)
}

type locationHandlerImpl struct {
	location   services.LocationService
	validation *validator.Validate
}

func NewLocationHandler(location services.LocationService, validation *validator.Validate) LocationHandler {
	return &locationHandlerImpl{
		location:   location,
		validation: validation,
	}
}

func (h *locationHandlerImpl) AddLocation(c *gin.Context) {
	var location dto.Location

	if err := c.ShouldBindJSON(&location); err != nil {
		helpers.BadRequestError(c, err.Error())
		return
	}

	if status, msg := utils.Validate(location); msg != "" {
		helpers.ErrorByCode(c, status, msg)
		return
	}

	code, err := h.location.Save(&location)
	if err != nil {
		helpers.SuccessByCode(c, code, nil)
		return
	}

	helpers.Created(c, "Successfully Created Data")
}

func (h *locationHandlerImpl) GetAllLocations(c *gin.Context) {
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

	users, code, err := h.location.GetAll(pagination)
	if err != nil {
		helpers.ErrorByCode(c, code, err.Error())
		return
	}

	helpers.OKWithMetadata(c, users, helpers.Metadata{
		Page: page,
		Size: size,
	})
}
