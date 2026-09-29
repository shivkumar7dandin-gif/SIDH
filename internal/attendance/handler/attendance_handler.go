package handler

import (
	"net/http"
	"strings"
	"time"

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
// CREATE / UPDATE BULK ATTENDANCE
// ========================================

func (h *AttendanceHandler) CreateBulk(
	c *gin.Context,
) {

	// College ID always comes from JWT.
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

	var req model.BulkAttendanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	attendanceList, err :=
		h.service.CreateBulk(
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
		http.StatusOK,
		gin.H{
			"message":    "attendance saved successfully",
			"attendance": attendanceList,
			"total":      len(attendanceList),
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
// GET CLASSROOM MONTHLY ATTENDANCE SUMMARY
// ========================================

func (h *AttendanceHandler) GetClassroomMonthlySummary(
	c *gin.Context,
) {

	// ========================================
	// 1. Get college ID from JWT
	// ========================================

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

	// ========================================
	// 2. Get classroom ID
	// ========================================

	classroomID :=
		strings.TrimSpace(
			c.Param("classroomId"),
		)

	if classroomID == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "classroom_id is required",
			},
		)
		return
	}

	// Validate ObjectID format
	if _, err := bson.ObjectIDFromHex(
		classroomID,
	); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid classroom_id",
			},
		)
		return
	}

	// ========================================
	// 3. Get academic year
	// ========================================

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

	// ========================================
	// 4. Get month
	// Expected format: 2026-09
	// ========================================

	monthValue :=
		strings.TrimSpace(
			c.Query("month"),
		)

	if monthValue == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "month is required in YYYY-MM format",
			},
		)
		return
	}

	parsedMonth, err :=
		time.Parse(
			"2006-01",
			monthValue,
		)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "month must be in YYYY-MM format",
			},
		)
		return
	}

	year := parsedMonth.Year()
	month := int(parsedMonth.Month())

	// ========================================
	// 5. Get classroom monthly summary
	// ========================================

	summary, err :=
		h.service.GetClassroomMonthlySummary(
			c.Request.Context(),
			collegeID,
			classroomID,
			academicYear,
			year,
			month,
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

	// ========================================
	// 6. Response
	// ========================================

	c.JSON(
		http.StatusOK,
		summary,
	)
}
