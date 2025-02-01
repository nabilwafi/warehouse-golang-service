package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nabilwafi/warehouse-management-system/exception"
	"github.com/nabilwafi/warehouse-management-system/models/web"
	"github.com/nabilwafi/warehouse-management-system/services"
)

type LocationHandlerImpl struct {
	service services.LocationService
}

func NewLocationHandler(service services.LocationService) LocationHandler {
	return &LocationHandlerImpl{service: service}
}

func (h *LocationHandlerImpl) AddLocation(c *gin.Context) {
	var location web.LocationCreateRequest
	if err := c.ShouldBindJSON(&location); err != nil {
		c.Error(exception.NewCustomError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.service.AddLocation(c.Request.Context(), &location); err != nil {
		c.Error(exception.NewCustomError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, location)
}

func (h *LocationHandlerImpl) GetAllLocations(c *gin.Context) {
	locations, err := h.service.GetAllLocations(c.Request.Context())
	if err != nil {
		c.Error(exception.NewCustomError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, locations)
}
