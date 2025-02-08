package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Status       string `json:"status"`
	StatusCode   int    `json:"status_code"`
	Data         any    `json:"data,omitempty"`
	ErrorMessage any    `json:"error_message,omitempty"`
	Metadata     `json:"metadata,omitempty"`
}

type Metadata struct {
	Page int `json:"page,omitempty"`
	Size int `json:"size,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, &BaseResponse{
		Status:     "OK",
		StatusCode: http.StatusOK,
		Data:       data,
	})
}

func OKWithMetadata(c *gin.Context, data any, pagination Metadata) {
	c.JSON(http.StatusOK, &BaseResponse{
		Status:     "OK",
		StatusCode: http.StatusOK,
		Data:       data,
		Metadata:   pagination,
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, &BaseResponse{
		Status:     "CREATED",
		StatusCode: http.StatusCreated,
	})
}

func BadRequestError(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, &BaseResponse{
		Status:       "BAD REQUEST",
		StatusCode:   http.StatusBadRequest,
		ErrorMessage: msg,
	})
}

func NotFoundError(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, &BaseResponse{
		Status:       "NOT FOUND",
		StatusCode:   http.StatusNotFound,
		ErrorMessage: msg,
	})
}

func UnauthorizedError(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, &BaseResponse{
		Status:       "UNAUTHORIZED",
		StatusCode:   http.StatusUnauthorized,
		ErrorMessage: msg,
	})
}

func ForbiddenError(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, &BaseResponse{
		Status:       "FORBIDDEN",
		StatusCode:   http.StatusForbidden,
		ErrorMessage: msg,
	})
}

func InternalServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, &BaseResponse{
		Status:       "INTERNAL SERVER ERROR",
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: msg,
	})
}

func ErrorByCode(c *gin.Context, errorCode int, msg string) {
	switch errorCode {
	case 400:
		BadRequestError(c, msg)
	case 401:
		UnauthorizedError(c, msg)
	case 403:
		ForbiddenError(c, msg)
	case 404:
		NotFoundError(c, msg)
	default:
		InternalServerError(c, msg)
	}
}

func SuccessByCode(c *gin.Context, successCode int, data any) {
	switch successCode {
	case 200:
		OK(c, data)
	default:
		Created(c, data)
	}
}
