package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shivkumar7dandin-gif/students-api/internal/attendance/model"
	"github.com/shivkumar7dandin-gif/students-api/internal/attendance/service"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AttendanceHandler struct {
	service *service.AttendanceService
}

func NewAttendanceHandler(
	service *service.AttendanceService,
) *AttendanceHandler {
	return &AttendanceHandler{
		service: service,
	}
}

func (h *AttendanceHandler) Create(c *gin.Context) {

	var attendance model.Attendance

	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	if attendance.StudentID.IsZero() {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "student_id is required"},
		)
		return
	}

	createdAttendance, err := h.service.Create(
		c.Request.Context(),
		attendance,
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
			"message":    "attendance created successfully",
			"attendance": createdAttendance,
		},
	)
}

func (h *AttendanceHandler) GetByStudent(c *gin.Context) {

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

	attendance, err := h.service.GetByStudent(
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

	c.JSON(http.StatusOK, attendance)
}

func (h *AttendanceHandler) GetAll(c *gin.Context) {

	attendance, err := h.service.GetAll(
		c.Request.Context(),
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(http.StatusOK, attendance)
}

func (h *AttendanceHandler) GetSummary(c *gin.Context) {

	studentID, err := bson.ObjectIDFromHex(
		c.Param("studentId"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid student id",
			},
		)
		return
	}

	summary, err := h.service.GetSummary(
		c.Request.Context(),
		studentID,
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
		summary,
	)
}
