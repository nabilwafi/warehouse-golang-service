package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/exception"
	"github.com/nabilwafi/warehouse-management-system/models/web"
	"github.com/nabilwafi/warehouse-management-system/services"
)

type userHandlerImpl struct {
	userService services.UserService
}

func NewUserHandlerImpl(userService services.UserService) UserHandler {
	return &userHandlerImpl{userService: userService}
}

func (h *userHandlerImpl) Register(c *gin.Context) {
	var user web.Register
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(exception.NewCustomError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.userService.Register(c.Request.Context(), &user); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (h *userHandlerImpl) Login(c *gin.Context) {
	var user web.Login
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(exception.NewCustomError(http.StatusBadRequest, err.Error()))
		return
	}

	token, err := h.userService.Login(c.Request.Context(), &user)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *userHandlerImpl) GetMe(c *gin.Context) {
	userID, _ := c.Get("user")

	user, err := h.userService.GetMe(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandlerImpl) ListUsers(c *gin.Context) {
	var userQuery web.UserQuery
	if err := c.ShouldBindQuery(&userQuery); err != nil {
		c.Error(exception.NewCustomError(http.StatusBadRequest, err.Error()))
		return
	}

	users, err := h.userService.ListUser(c.Request.Context(), &userQuery)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, users)
}
