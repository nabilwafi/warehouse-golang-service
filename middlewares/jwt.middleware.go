package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nabilwafi/warehouse-management-system/exception"
	"github.com/nabilwafi/warehouse-management-system/utils"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user", claims)

		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")

		if !exists {
			c.Error(exception.NewCustomError(http.StatusForbidden, "unauthorized"))
			c.Abort()
			return
		}

		decodedUser := user.(*utils.Claims)

		role := decodedUser.Role

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Set("user", decodedUser.ID)
				c.Next()
				return
			}
		}

		c.Error(exception.NewCustomError(http.StatusForbidden, "unauthorized"))
		c.Abort()
	}
}
