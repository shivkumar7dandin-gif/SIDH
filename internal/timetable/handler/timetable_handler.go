package handler

import (
	"net/http"
	"strings"

	timetableModel "github.com/shivkumar7dandin-gif/students-api/internal/timetable/model"
	timetableService "github.com/shivkumar7dandin-gif/students-api/internal/timetable/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TimetableHandler struct {
	service *timetableService.TimetableService
}

func NewTimetableHandler(
	service *timetableService.TimetableService,
) *TimetableHandler {

	return &TimetableHandler{
		service: service,
	}
}

// ========================================
// CREATE TIMETABLE
// ========================================

func (h *TimetableHandler) Create(c *gin.Context) {

	collegeIDValue, exists :=
		c.Get("college_id")

	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "college_id not found in token",
			},
		)
		return
	}

	collegeIDHex, ok :=
		collegeIDValue.(string)

	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id in token",
			},
		)
		return
	}

	collegeID, err :=
		bson.ObjectIDFromHex(collegeIDHex)

	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id in token",
			},
		)
		return
	}

	var req timetableModel.CreateTimetableRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request body",
			},
		)
		return
	}

	created, err :=
		h.service.Create(
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
		gin.H{
			"message":   "timetable created successfully",
			"timetable": created,
		},
	)
}

// ========================================
// GET BY ACADEMIC YEAR
// ========================================

func (h *TimetableHandler) GetAll(c *gin.Context) {

	collegeID, ok :=
		getCollegeID(c)

	if !ok {
		return
	}

	academicYear :=
		strings.TrimSpace(
			c.Query("academic_year"),
		)

	timetables, err :=
		h.service.GetByAcademicYear(
			c.Request.Context(),
			collegeID,
			academicYear,
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
			"timetables": timetables,
			"total":      len(timetables),
		},
	)
}

// ========================================
// GET BY DAY
// ========================================

func (h *TimetableHandler) GetByDay(c *gin.Context) {

	collegeID, ok :=
		getCollegeID(c)

	if !ok {
		return
	}

	day :=
		strings.TrimSpace(
			c.Param("day"),
		)

	academicYear :=
		strings.TrimSpace(
			c.Query("academic_year"),
		)

	timetables, err :=
		h.service.GetByDay(
			c.Request.Context(),
			collegeID,
			day,
			academicYear,
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
			"day":        day,
			"total":      len(timetables),
			"timetables": timetables,
		},
	)
}

// ========================================
// GET BY TEACHER + DAY
// ========================================

func (h *TimetableHandler) GetByTeacherAndDay(
	c *gin.Context,
) {

	collegeID, ok :=
		getCollegeID(c)

	if !ok {
		return
	}

	teacherIDHex :=
		strings.TrimSpace(
			c.Param("teacherId"),
		)

	teacherID, err :=
		bson.ObjectIDFromHex(
			teacherIDHex,
		)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid teacher id",
			},
		)
		return
	}

	day :=
		strings.TrimSpace(
			c.Query("day"),
		)

	academicYear :=
		strings.TrimSpace(
			c.Query("academic_year"),
		)

	timetables, err :=
		h.service.GetByTeacherAndDay(
			c.Request.Context(),
			collegeID,
			teacherID,
			day,
			academicYear,
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
			"teacher_id": teacherIDHex,
			"day":        day,
			"total":      len(timetables),
			"timetables": timetables,
		},
	)
}

// ========================================
// COMMON COLLEGE ID HELPER
// ========================================

func getCollegeID(
	c *gin.Context,
) (bson.ObjectID, bool) {

	collegeIDValue, exists :=
		c.Get("college_id")

	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "college_id not found in token",
			},
		)

		return bson.NilObjectID, false
	}

	collegeIDHex, ok :=
		collegeIDValue.(string)

	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id in token",
			},
		)

		return bson.NilObjectID, false
	}

	collegeID, err :=
		bson.ObjectIDFromHex(
			collegeIDHex,
		)

	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id in token",
			},
		)

		return bson.NilObjectID, false
	}

	return collegeID, true
}

// ========================================
// UPDATE TIMETABLE
// ========================================

func (h *TimetableHandler) Update(c *gin.Context) {

	// Get college ID from JWT context
	collegeID, ok := getCollegeID(c)

	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college id",
			},
		)
		return
	}

	// Convert timetable ID
	timetableID, err := bson.ObjectIDFromHex(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid timetable id",
			},
		)
		return
	}

	// Read request body
	var req timetableModel.UpdateTimetableRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request body",
			},
		)
		return
	}

	// Update timetable
	updatedTimetable, err := h.service.Update(
		c.Request.Context(),
		collegeID,
		timetableID,
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
			"message":   "timetable updated successfully",
			"timetable": updatedTimetable,
		},
	)
}

// ========================================
// DELETE TIMETABLE
// ========================================

func (h *TimetableHandler) Delete(c *gin.Context) {

	// Get college ID from JWT context
	collegeID, ok := getCollegeID(c)

	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college id",
			},
		)
		return
	}

	// Convert timetable ID
	timetableID, err := bson.ObjectIDFromHex(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid timetable id",
			},
		)
		return
	}

	// Delete timetable
	err = h.service.Delete(
		c.Request.Context(),
		collegeID,
		timetableID,
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
			"message": "timetable deleted successfully",
		},
	)
}
