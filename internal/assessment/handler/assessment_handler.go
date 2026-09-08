package handler

import (
	"net/http"
	"strconv"
	"strings"

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

// ========================================
// HELPER - GET COLLEGE ID FROM JWT CONTEXT
// ========================================

func getCollegeID(c *gin.Context) (bson.ObjectID, bool) {

	collegeIDValue, exists := c.Get("college_id")

	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "college_id not found in token"},
		)
		return bson.NilObjectID, false
	}

	collegeIDString, ok := collegeIDValue.(string)

	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "invalid college_id in token"},
		)
		return bson.NilObjectID, false
	}

	collegeID, err := bson.ObjectIDFromHex(collegeIDString)

	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "invalid college_id in token"},
		)
		return bson.NilObjectID, false
	}

	return collegeID, true
}

// ========================================
// CREATE
// ========================================

func (h *AssessmentHandler) Create(c *gin.Context) {

	collegeID, ok := getCollegeID(c)

	if !ok {
		return
	}

	var assessment model.Assessment

	if err := c.ShouldBindJSON(&assessment); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	createdAssessment, err := h.service.Create(
		c.Request.Context(),
		collegeID,
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

// ========================================
// GET ALL
// ========================================

func (h *AssessmentHandler) GetAll(c *gin.Context) {

	collegeID, ok := getCollegeID(c)

	if !ok {
		return
	}

	academicYear :=
		strings.TrimSpace(
			c.Query("academic_year"),
		)

	if academicYear == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "academic_year is required",
			},
		)
		return
	}

	assessments, err := h.service.GetAll(
		c.Request.Context(),
		collegeID,
		academicYear,
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

// ========================================
// GET BY STUDENT
// ========================================

func (h *AssessmentHandler) GetByStudent(c *gin.Context) {

	collegeID, ok := getCollegeID(c)

	if !ok {
		return
	}

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

	academicYear :=
		strings.TrimSpace(
			c.Query("academic_year"),
		)

	if academicYear == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "academic_year is required",
			},
		)
		return
	}

	assessments, err :=
		h.service.GetByStudent(
			c.Request.Context(),
			collegeID,
			studentID,
			academicYear,
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

func (h *AssessmentHandler) GetMonthlySummary(c *gin.Context) {

	collegeID, ok := getCollegeID(c)
	if !ok {
		return
	}

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

	academicYear := strings.TrimSpace(
		c.Query("academic_year"),
	)

	if academicYear == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "academic_year is required"},
		)
		return
	}

	year, err := strconv.Atoi(
		c.Query("year"),
	)
	if err != nil || year <= 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "valid year is required"},
		)
		return
	}

	month, err := strconv.Atoi(
		c.Query("month"),
	)
	if err != nil || month < 1 || month > 12 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "month must be between 1 and 12"},
		)
		return
	}

	summary, err := h.service.GetMonthlySummary(
		c.Request.Context(),
		collegeID,
		studentID,
		academicYear,
		year,
		month,
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
		summary,
	)
}
