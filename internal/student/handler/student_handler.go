package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shivkumar7dandin-gif/students-api/internal/student/model"
	"github.com/shivkumar7dandin-gif/students-api/internal/student/service"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type StudentHandler struct {
	service *service.StudentService
}

func NewStudentHandler(
	service *service.StudentService,
) *StudentHandler {

	return &StudentHandler{
		service: service,
	}
}

// POST /api/v1/students
func (h *StudentHandler) Create(c *gin.Context) {

	var students []model.Student

	// Read array of students
	if err := c.ShouldBindJSON(&students); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check empty request
	if len(students) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "student list cannot be empty",
		})
		return
	}

	// Create students one by one
	createdStudents := make([]*model.Student, 0, len(students))

	for _, student := range students {

		if student.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "student name is required",
			})
			return
		}

		if student.Age <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "age must be greater than 0",
			})
			return
		}

		if student.RollNumber <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "roll number must be greater than 0",
			})
			return
		}

		if student.ClassroomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "classroom_id is required",
			})
			return
		}

		createdStudent, err := h.service.Create(
			c.Request.Context(),
			student,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		createdStudents = append(createdStudents, createdStudent)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "students created successfully",
		"students": createdStudents,
	})
}

// GET /api/v1/students
func (h *StudentHandler) GetAll(c *gin.Context) {

	students, err := h.service.GetAll(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, students)
}

// GET /api/v1/students/:id
func (h *StudentHandler) GetByID(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid student id",
		})
		return
	}

	student, err := h.service.GetByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "student not found",
		})
		return
	}

	c.JSON(http.StatusOK, student)
}

// PUT /api/v1/students/:id
func (h *StudentHandler) Update(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid student id",
		})
		return
	}

	var student model.Student

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.service.Update(
		c.Request.Context(),
		id,
		student,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "student not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "student updated successfully",
	})
}

// DELETE /api/v1/students/:id
func (h *StudentHandler) Delete(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid student id",
		})
		return
	}

	err = h.service.Delete(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "student not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "student deleted successfully",
	})
}
