package handler

import (
	"net/http"
	"strconv"
	"strings"

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

// ========================================
// GET COLLEGE ID FROM JWT
// ========================================

func getCollegeID(
	c *gin.Context,
) (bson.ObjectID, error) {

	value, exists := c.Get("college_id")
	if !exists {
		return bson.NilObjectID,
			http.ErrNoCookie
	}

	collegeIDString, ok := value.(string)

	if !ok || collegeIDString == "" {
		return bson.NilObjectID,
			http.ErrNoCookie
	}

	return bson.ObjectIDFromHex(
		collegeIDString,
	)
}

// ========================================
// CREATE
// ========================================

func (h *AttendanceHandler) Create(
	c *gin.Context,
) {

	collegeID, err := getCollegeID(c)

	if err != nil {

		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id",
			},
		)
		return
	}

	var attendance model.Attendance

	if err := c.ShouldBindJSON(
		&attendance,
	); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	createdAttendance, err :=
		h.service.Create(
			c.Request.Context(),
			collegeID,
			attendance,
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
		gin.H{
			"message":    "attendance created successfully",
			"attendance": createdAttendance,
		},
	)
}

// ========================================
// GET BY STUDENT
// ========================================

func (h *AttendanceHandler) GetByStudent(
	c *gin.Context,
) {

	collegeID, err := getCollegeID(c)

	if err != nil {

		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id",
			},
		)
		return
	}

	studentID, err :=
		bson.ObjectIDFromHex(
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

	academicYear :=
		strings.TrimSpace(
			c.Query("academic_year"),
		)

	if academicYear == "" {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "academic_year query parameter is required",
			},
		)
		return
	}

	attendance, err :=
		h.service.GetByStudent(
			c.Request.Context(),
			collegeID,
			studentID,
			academicYear,
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
		attendance,
	)
}

// ========================================
// GET ALL
// ========================================

func (h *AttendanceHandler) GetAll(
	c *gin.Context,
) {

	collegeID, err := getCollegeID(c)

	if err != nil {

		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id",
			},
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
				"error": "academic_year query parameter is required",
			},
		)
		return
	}

	attendance, err :=
		h.service.GetAll(
			c.Request.Context(),
			collegeID,
			academicYear,
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
		attendance,
	)
}

// ========================================
// SUMMARY
// ========================================

func (h *AttendanceHandler) GetSummary(
	c *gin.Context,
) {

	collegeID, err := getCollegeID(c)

	if err != nil {

		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id",
			},
		)
		return
	}

	studentID, err :=
		bson.ObjectIDFromHex(
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

	academicYear :=
		strings.TrimSpace(
			c.Query("academic_year"),
		)

	if academicYear == "" {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "academic_year query parameter is required",
			},
		)
		return
	}

	summary, err :=
		h.service.GetSummary(
			c.Request.Context(),
			collegeID,
			studentID,
			academicYear,
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

// ========================================
// GET MONTHLY ATTENDANCE SUMMARY
// ========================================

func (h *AttendanceHandler) GetMonthlySummary(
	c *gin.Context,
) {

	// Get college ID from JWT
	collegeID, err := getCollegeID(c)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "invalid college"},
		)
		return
	}

	// Get student ID
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

	// Get academic year
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

	// Get year
	year, err := strconv.Atoi(
		c.Query("year"),
	)
	if err != nil || year <= 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "valid year is required",
			},
		)
		return
	}

	// Get month
	month, err := strconv.Atoi(
		c.Query("month"),
	)
	if err != nil ||
		month < 1 ||
		month > 12 {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "month must be between 1 and 12",
			},
		)
		return
	}

	// Get monthly summary
	summary, err :=
		h.service.GetMonthlySummary(
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
