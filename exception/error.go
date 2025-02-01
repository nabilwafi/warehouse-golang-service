package exception

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err != nil {
			switch e := err.Err.(type) {
			case *CustomError:
				c.JSON(e.Code, CustomError{
					Code:    e.Code,
					Message: e.Message,
				})
			default:
				c.JSON(500, CustomError{
					Code:    500,
					Message: e.Error(),
				})
			}

			c.Abort()
		}
	}
}

func (ce *CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", ce.Code, ce.Message)
}

func NewCustomError(code int, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}
