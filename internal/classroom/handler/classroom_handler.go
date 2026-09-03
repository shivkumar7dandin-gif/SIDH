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

// func (h *ClassroomHandler) Create(c *gin.Context) {

// 	var classroom model.Classroom

// 	// Bind JSON request
// 	if err := c.ShouldBindJSON(&classroom); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	// Validate classroom name
// 	if classroom.Name == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "classroom name is required",
// 		})
// 		return
// 	}

// 	// Validate capacity
// 	if classroom.Capacity <= 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "capacity must be greater than 0",
// 		})
// 		return
// 	}

// 	// Create classroom
// 	createdClassroom, err := h.service.Create(
// 		c.Request.Context(),
// 		classroom,
// 	)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	// Return created classroom
// 	c.JSON(http.StatusCreated, gin.H{
// 		"message":   "classroom created successfully",
// 		"classroom": createdClassroom,
// 	})
// }

func (h *ClassroomHandler) Create(c *gin.Context) {

	var req model.CreateClassroomRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	referenceIDValue, exists := c.Get("reference_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "college reference id not found",
		})
		return
	}

	collegeIDString, ok := referenceIDValue.(string)
	if !ok || collegeIDString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid college reference id",
		})
		return
	}

	collegeID, err := bson.ObjectIDFromHex(collegeIDString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid college id",
		})
		return
	}

	classroom, err := h.service.Create(
		c.Request.Context(),
		req,
		collegeID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, classroom)
}

func (h *ClassroomHandler) GetAll(c *gin.Context) {

	classrooms, err := h.service.GetAll(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, classrooms)
}

func (h *ClassroomHandler) GetByID(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid classroom id",
		})
		return
	}

	classroom, err := h.service.GetByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "classroom not found",
		})
		return
	}

	c.JSON(http.StatusOK, classroom)
}

func (h *ClassroomHandler) Update(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid classroom id",
		})
		return
	}

	var classroom model.Classroom

	if err := c.ShouldBindJSON(&classroom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.Update(
		c.Request.Context(),
		id,
		classroom,
	); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "classroom updated successfully",
	})
}

func (h *ClassroomHandler) Delete(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid classroom id",
		})
		return
	}

	if err := h.service.Delete(
		c.Request.Context(),
		id,
	); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "classroom deleted successfully",
	})
}
