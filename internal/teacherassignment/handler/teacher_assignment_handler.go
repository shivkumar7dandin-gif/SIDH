package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	assignmentModel "github.com/shivkumar7dandin-gif/students-api/internal/teacherassignment/model"
	assignmentService "github.com/shivkumar7dandin-gif/students-api/internal/teacherassignment/service"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TeacherAssignmentHandler struct {
	service *assignmentService.TeacherAssignmentService
}

func NewTeacherAssignmentHandler(
	service *assignmentService.TeacherAssignmentService,
) *TeacherAssignmentHandler {

	return &TeacherAssignmentHandler{
		service: service,
	}
}

// =====================================================
// CREATE TEACHER CLASS ASSIGNMENT
// =====================================================

func (h *TeacherAssignmentHandler) Create(
	c *gin.Context,
) {

	var req assignmentModel.CreateTeacherClassAssignmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	collegeID, ok := getCollegeID(c)
	if !ok {
		return
	}

	assignment, err := h.service.Create(
		c.Request.Context(),
		collegeID,
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
		assignment,
	)
}

// =====================================================
// GET ASSIGNMENTS FOR ONE TEACHER
// =====================================================

func (h *TeacherAssignmentHandler) GetByTeacher(
	c *gin.Context,
) {

	teacherID := c.Param("id")

	collegeID, ok := getCollegeID(c)
	if !ok {
		return
	}

	assignments, err := h.service.GetByTeacher(
		c.Request.Context(),
		collegeID,
		teacherID,
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
		assignments,
	)
}

// =====================================================
// DELETE TEACHER CLASS ASSIGNMENT
// =====================================================

func (h *TeacherAssignmentHandler) Delete(
	c *gin.Context,
) {

	assignmentID := c.Param("id")

	collegeID, ok := getCollegeID(c)
	if !ok {
		return
	}

	err := h.service.Delete(
		c.Request.Context(),
		collegeID,
		assignmentID,
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
			"message": "teacher assignment deleted successfully",
		},
	)
}

// =====================================================
// HELPER: GET COLLEGE ID FROM JWT CONTEXT
// =====================================================

func getCollegeID(
	c *gin.Context,
) (bson.ObjectID, bool) {

	collegeIDValue, exists := c.Get("college_id")

	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "college_id not found",
			},
		)

		return bson.ObjectID{}, false
	}

	collegeIDString, ok := collegeIDValue.(string)

	if !ok || collegeIDString == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id",
			},
		)

		return bson.ObjectID{}, false
	}

	collegeID, err := bson.ObjectIDFromHex(
		collegeIDString,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid college id",
			},
		)

		return bson.ObjectID{}, false
	}

	return collegeID, true
}
