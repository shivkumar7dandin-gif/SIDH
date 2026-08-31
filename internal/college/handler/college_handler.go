package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shivkumar7dandin-gif/students-api/internal/college/model"
	"github.com/shivkumar7dandin-gif/students-api/internal/college/service"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CollegeHandler struct {
	service *service.CollegeService
}

func NewCollegeHandler(
	service *service.CollegeService,
) *CollegeHandler {
	return &CollegeHandler{
		service: service,
	}
}

func (h *CollegeHandler) Create(c *gin.Context) {

	var req model.CreateCollegeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	createdCollege, err := h.service.Create(
		c.Request.Context(),
		req,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"message": "college created successfully",
			"college": createdCollege,
		},
	)
}

func (h *CollegeHandler) GetAll(c *gin.Context) {

	colleges, err := h.service.GetAll(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		colleges,
	)
}

func (h *CollegeHandler) GetByID(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid college id"},
		)
		return
	}

	college, err := h.service.GetByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{"error": "college not found"},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		college,
	)
}
