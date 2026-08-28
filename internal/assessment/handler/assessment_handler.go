package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shivkumar7dandin-gif/students-api/internal/assessment/model"
	"github.com/shivkumar7dandin-gif/students-api/internal/assessment/service"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AssessmentHandler struct {
	service *service.AssessmentService
}

func NewAssessmentHandler(
	service *service.AssessmentService,
) *AssessmentHandler {

	return &AssessmentHandler{
		service: service,
	}
}

func (h *AssessmentHandler) Create(c *gin.Context) {

	var assessment model.Assessment

	if err := c.ShouldBindJSON(&assessment); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	if assessment.StudentID.IsZero() {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "student_id is required"},
		)
		return
	}

	createdAssessment, err := h.service.Create(
		c.Request.Context(),
		assessment,
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
			"message":    "assessment created successfully",
			"assessment": createdAssessment,
		},
	)
}

func (h *AssessmentHandler) GetAll(c *gin.Context) {

	assessments, err := h.service.GetAll(
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
		assessments,
	)
}

func (h *AssessmentHandler) GetByStudent(c *gin.Context) {

	studentID, err := bson.ObjectIDFromHex(
		c.Param("studentId"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid student id"},
		)
		return
	}

	assessments, err := h.service.GetByStudent(
		c.Request.Context(),
		studentID,
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
		assessments,
	)
}
