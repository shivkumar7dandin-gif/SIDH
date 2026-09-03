package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	teacherModel "github.com/shivkumar7dandin-gif/students-api/internal/teacher/model"
	teacherService "github.com/shivkumar7dandin-gif/students-api/internal/teacher/service"
)

type TeacherHandler struct {
	service *teacherService.TeacherService
}

func NewTeacherHandler(
	service *teacherService.TeacherService,
) *TeacherHandler {

	return &TeacherHandler{
		service: service,
	}
}

func (h *TeacherHandler) Create(c *gin.Context) {

	var req teacherModel.CreateTeacherRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request body",
			},
		)
		return
	}

	teacher, err := h.service.Create(
		c.Request.Context(),
		req,
	)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		teacher,
	)
}

func (h *TeacherHandler) GetAll(c *gin.Context) {

	teachers, err := h.service.GetAll(
		c.Request.Context(),
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		teachers,
	)
}

func (h *TeacherHandler) GetByID(c *gin.Context) {

	teacher, err := h.service.GetByID(
		c.Request.Context(),
		c.Param("id"),
	)
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		teacher,
	)
}

func (h *TeacherHandler) Update(c *gin.Context) {

	var teacher teacherModel.Teacher

	if err := c.ShouldBindJSON(
		&teacher,
	); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request body",
			},
		)
		return
	}

	err := h.service.Update(
		c.Request.Context(),
		c.Param("id"),
		teacher,
	)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "teacher updated successfully",
		},
	)
}

func (h *TeacherHandler) Delete(c *gin.Context) {

	err := h.service.Delete(
		c.Request.Context(),
		c.Param("id"),
	)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "teacher deleted successfully",
		},
	)
}

// =====================================================
// GET LOGGED-IN TEACHER DETAILS
// GET /api/v1/teachers/me
// =====================================================

func (h *TeacherHandler) GetMe(c *gin.Context) {

	// Get reference_id added by AuthMiddleware
	referenceIDValue, exists := c.Get("reference_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "reference_id not found",
		})
		return
	}

	referenceID, ok := referenceIDValue.(string)

	if !ok || referenceID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid reference_id",
		})
		return
	}

	// Convert string ID to MongoDB ObjectID
	teacher, err := h.service.GetByID(
		c.Request.Context(),
		referenceID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "teacher not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "teacher not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"teacher": teacher,
	})
}
