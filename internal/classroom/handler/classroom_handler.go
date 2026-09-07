package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shivkumar7dandin-gif/students-api/internal/classroom/model"
	"github.com/shivkumar7dandin-gif/students-api/internal/classroom/service"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ClassroomHandler struct {
	service *service.ClassroomService
}

func NewClassroomHandler(
	service *service.ClassroomService,
) *ClassroomHandler {

	return &ClassroomHandler{
		service: service,
	}
}

// =====================================================
// CREATE CLASSROOM
// =====================================================

func (h *ClassroomHandler) Create(c *gin.Context) {

	var req model.CreateClassroomRequest

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

	classroom, err := h.service.Create(
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
		classroom,
	)
}

// =====================================================
// GET ALL CLASSROOMS
// =====================================================

func (h *ClassroomHandler) GetAll(c *gin.Context) {

	collegeID, ok := getCollegeID(c)
	if !ok {
		return
	}

	classrooms, err := h.service.GetByCollegeID(
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
		classrooms,
	)
}

// =====================================================
// GET CLASSROOM BY ID
// =====================================================

func (h *ClassroomHandler) GetByID(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid classroom id",
			},
		)
		return
	}

	collegeID, ok := getCollegeID(c)
	if !ok {
		return
	}

	classroom, err := h.service.GetByIDAndCollegeID(
		c.Request.Context(),
		id,
		collegeID,
	)

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "classroom not found",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		classroom,
	)
}

// =====================================================
// UPDATE CLASSROOM
// =====================================================

func (h *ClassroomHandler) Update(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid classroom id",
			},
		)
		return
	}

	var classroom model.Classroom

	if err := c.ShouldBindJSON(&classroom); err != nil {
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

	err = h.service.UpdateByCollegeID(
		c.Request.Context(),
		id,
		collegeID,
		classroom,
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
			"message": "classroom updated successfully",
		},
	)
}

// =====================================================
// DELETE CLASSROOM
// =====================================================

func (h *ClassroomHandler) Delete(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid classroom id",
			},
		)
		return
	}

	collegeID, ok := getCollegeID(c)
	if !ok {
		return
	}

	err = h.service.DeleteByCollegeID(
		c.Request.Context(),
		id,
		collegeID,
	)

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "classroom not found",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "classroom deleted successfully",
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
