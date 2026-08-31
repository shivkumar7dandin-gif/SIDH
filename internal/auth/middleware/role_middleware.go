package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "role not found"},
			)
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid role"},
			)
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(
			http.StatusForbidden,
			gin.H{
				"error": "you are not authorized to perform this action",
			},
		)
		c.Abort()
	}
}
