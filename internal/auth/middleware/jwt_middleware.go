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

		// ---------------------------------
		// GET AUTHORIZATION HEADER
		// ---------------------------------

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

		// ---------------------------------
		// CHECK "Bearer <token>"
		// ---------------------------------

		parts := strings.SplitN(
			authHeader,
			" ",
			2,
		)

		if len(parts) != 2 ||
			!strings.EqualFold(
				parts[0],
				"Bearer",
			) {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid authorization header",
				},
			)

			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(
			parts[1],
		)

		if tokenString == "" {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "token is required",
				},
			)

			c.Abort()
			return
		}

		// ---------------------------------
		// CREATE CLAIMS OBJECT
		// ---------------------------------

		claims := &authService.Claims{}

		// ---------------------------------
		// PARSE AND VALIDATE JWT
		// ---------------------------------

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (
				interface{},
				error,
			) {
				return []byte(jwtSecret), nil
			},
			jwt.WithValidMethods(
				[]string{
					jwt.SigningMethodHS256.Alg(),
				},
			),
		)

		if err != nil ||
			token == nil ||
			!token.Valid {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid or expired token",
				},
			)

			c.Abort()
			return
		}

		// ---------------------------------
		// VALIDATE REQUIRED CLAIMS
		// ---------------------------------

		if claims.UserID == "" ||
			claims.Username == "" ||
			claims.Role == "" {

			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid token claims",
				},
			)

			c.Abort()
			return
		}

		// ---------------------------------
		// STORE JWT DATA IN GIN CONTEXT
		// ---------------------------------

		c.Set(
			"user_id",
			claims.UserID,
		)

		c.Set(
			"reference_id",
			claims.ReferenceID,
		)

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

		// ---------------------------------
		// CONTINUE TO NEXT HANDLER
		// ---------------------------------

		c.Next()
	}
}
