package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nabilwafi/warehouse-management-system/src/helpers"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/nabilwafi/warehouse-management-system/src/services"
	"github.com/nabilwafi/warehouse-management-system/src/utils"
)

type UserHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	GetMe(c *gin.Context)
	ListUsers(c *gin.Context)
}

type userHandlerImpl struct {
	user services.UserService
}

func NewUserHandlerImpl(user services.UserService) UserHandler {
	return &userHandlerImpl{user: user}
}

func (h *userHandlerImpl) Register(c *gin.Context) {
	var register dto.RegisterRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		helpers.BadRequestError(c, err.Error())
		return
	}

	if status, msg := utils.Validate(&register); msg != "" {
		helpers.ErrorByCode(c, status, msg)
		return
	}

	code, err := h.user.Register(&register)
	if err != nil {
		helpers.ErrorByCode(c, code, err.Error())
		return
	}

	helpers.Created(c, "Successfully Created User")
}

func (h *userHandlerImpl) Login(c *gin.Context) {
	var login dto.LoginRequest

	if err := c.ShouldBindJSON(&login); err != nil {
		helpers.BadRequestError(c, err.Error())
		return
	}

	if status, msg := utils.Validate(&login); msg != "" {
		helpers.ErrorByCode(c, status, msg)
		return
	}

	token, code, err := h.user.Login(&login)
	if err != nil {
		helpers.ErrorByCode(c, code, err.Error())
		return
	}

	helpers.OK(c, token)
}

func (h *userHandlerImpl) GetMe(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		helpers.UnauthorizedError(c, "please login first")
		return
	}

	userData, ok := user.(*utils.CustomClaims)
	if !ok {
		helpers.InternalServerError(c, "internal server error")
		return
	}

	user, code, err := h.user.GetUserByID(userData.ID)
	if err != nil {
		helpers.ErrorByCode(c, code, err.Error())
		return
	}

	helpers.OK(c, user)
}

func (h *userHandlerImpl) ListUsers(c *gin.Context) {
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

	users, code, err := h.user.GetAllUser(pagination)
	if err != nil {
		helpers.ErrorByCode(c, code, err.Error())
		return
	}

	helpers.OKWithMetadata(c, users, helpers.Metadata{
		Page: page,
		Size: size,
	})
}
