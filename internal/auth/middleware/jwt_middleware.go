package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	authService "github.com/shivkumar7dandin-gif/students-api/internal/auth/service"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "authorization header is required",
				},
			)
			c.Abort()
			return
		}

		parts := strings.SplitN(
			authHeader,
			" ",
			2,
		)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid authorization header",
				},
			)
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims := &authService.Claims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			},
			jwt.WithValidMethods(
				[]string{
					jwt.SigningMethodHS256.Alg(),
				},
			),
		)

		if err != nil || !token.Valid {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid or expired token",
				},
			)
			c.Abort()
			return
		}

		c.Set(
			"college_id",
			claims.CollegeID,
		)

		c.Set(
			"username",
			claims.Username,
		)

		c.Set(
			"role",
			claims.Role,
		)

		c.Next()
	}
}
