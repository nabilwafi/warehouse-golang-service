package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nabilwafi/warehouse-management-system/src/helpers"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/utils"
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

func RoleMiddleware(allowedRoles ...dto.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")

		if !exists {
			helpers.ForbiddenError(c, "forbidden")
			c.Abort()
		}

		decodedUser := user.(*utils.CustomClaims)

		role := decodedUser.Role

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Set("user", decodedUser)
				c.Next()
				return
			}
		}

		helpers.ForbiddenError(c, "forbidden")
		c.Abort()
	}
}
