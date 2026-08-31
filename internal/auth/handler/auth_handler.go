package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	authModel "github.com/shivkumar7dandin-gif/students-api/internal/auth/model"
	authService "github.com/shivkumar7dandin-gif/students-api/internal/auth/service"
)

type AuthHandler struct {
	service *authService.AuthService
}

func NewAuthHandler(
	service *authService.AuthService,
) *AuthHandler {

	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {

	var req authModel.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request body",
			},
		)
		return
	}

	req.Username = strings.TrimSpace(req.Username)

	if req.Username == "" || req.Password == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "username and password are required",
			},
		)
		return
	}

	token, role, err := h.service.Login(
		c.Request.Context(),
		req.Username,
		req.Password,
	)

	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		authModel.LoginResponse{
			Token: token,
			Role:  role,
		},
	)
}
