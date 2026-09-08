package handler

import (
	"net/http"
	"strconv"
	"strings"

	reportService "github.com/shivkumar7dandin-gif/students-api/internal/report/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ReportHandler struct {
	service *reportService.ReportService
}

func NewReportHandler(
	service *reportService.ReportService,
) *ReportHandler {

	return &ReportHandler{
		service: service,
	}
}

func (h *ReportHandler) GetMonthlyParentReport(
	c *gin.Context,
) {

	// ---------------------------------------
	// 1. Get college ID from JWT
	// ---------------------------------------

	collegeIDValue, exists := c.Get("college_id")

	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "college_id missing from token"},
		)
		return
	}

	collegeIDString, ok := collegeIDValue.(string)

	if !ok || strings.TrimSpace(collegeIDString) == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "invalid college_id in token"},
		)
		return
	}

	collegeID, err :=
		bson.ObjectIDFromHex(collegeIDString)

	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "invalid college_id in token"},
		)
		return
	}

	// ---------------------------------------
	// 2. Student ID
	// ---------------------------------------

	studentID, err :=
		bson.ObjectIDFromHex(c.Param("studentId"))

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid student_id"},
		)
		return
	}

	// ---------------------------------------
	// 3. Academic year
	// ---------------------------------------

	academicYear :=
		strings.TrimSpace(c.Query("academic_year"))

	if academicYear == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "academic_year is required"},
		)
		return
	}

	// ---------------------------------------
	// 4. Year
	// ---------------------------------------

	year, err := strconv.Atoi(c.Query("year"))

	if err != nil || year <= 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "valid year is required"},
		)
		return
	}

	// ---------------------------------------
	// 5. Month
	// ---------------------------------------

	month, err := strconv.Atoi(c.Query("month"))

	if err != nil || month < 1 || month > 12 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "month must be between 1 and 12"},
		)
		return
	}

	// ---------------------------------------
	// 6. Generate report
	// ---------------------------------------

	report, err :=
		h.service.GetMonthlyParentReport(
			c.Request.Context(),
			collegeID,
			studentID,
			academicYear,
			year,
			month,
		)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(http.StatusOK, report)
}
