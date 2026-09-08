package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	calendarModel "github.com/shivkumar7dandin-gif/students-api/internal/calendar/model"
	calendarService "github.com/shivkumar7dandin-gif/students-api/internal/calendar/service"
)

type CalendarHandler struct {
	service *calendarService.CalendarService
}

func NewCalendarHandler(
	service *calendarService.CalendarService,
) *CalendarHandler {

	return &CalendarHandler{
		service: service,
	}
}

func getCalendarCollegeID(
	c *gin.Context,
) (bson.ObjectID, bool) {

	value, exists := c.Get("college_id")

	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "college_id not found",
			},
		)
		return bson.ObjectID{}, false
	}

	collegeIDString, ok := value.(string)

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

// POST /api/v1/calendars
func (h *CalendarHandler) Create(
	c *gin.Context,
) {

	collegeID, ok := getCalendarCollegeID(c)

	if !ok {
		return
	}

	var req calendarModel.CreateAcademicCalendarRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request body",
			},
		)
		return
	}

	calendar, err := h.service.Create(
		c.Request.Context(),
		req,
		collegeID,
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
		calendar,
	)
}

// GET /api/v1/calendars
func (h *CalendarHandler) GetAll(
	c *gin.Context,
) {

	collegeID, ok := getCalendarCollegeID(c)

	if !ok {
		return
	}

	calendars, err := h.service.GetByCollegeID(
		c.Request.Context(),
		collegeID,
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
		calendars,
	)
}

// GET /api/v1/calendars/:academicYear
func (h *CalendarHandler) GetByAcademicYear(
	c *gin.Context,
) {

	collegeID, ok := getCalendarCollegeID(c)

	if !ok {
		return
	}

	academicYear := c.Param(
		"academicYear",
	)

	calendar, err := h.service.GetByAcademicYear(
		c.Request.Context(),
		collegeID,
		academicYear,
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
		calendar,
	)
}
