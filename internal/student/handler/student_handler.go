package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

// =====================================================
// POST /api/v1/students
// Supports:
// 1. Single student
// 2. Bulk students
// =====================================================

func (h *StudentHandler) Create(c *gin.Context) {

	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read request body",
		})
		return
	}

	body = bytes.TrimSpace(body)

	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "request body cannot be empty",
		})
		return
	}

	// =====================================================
	// SINGLE STUDENT
	// =====================================================

	if body[0] == '{' {

		var student model.Student

		if err := json.Unmarshal(body, &student); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Validate student
		if err := validateStudent(student); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Create student
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

		c.JSON(http.StatusCreated, gin.H{
			"message": "student created successfully",
			"student": createdStudent,
		})

		return
	}

	// =====================================================
	// BULK STUDENTS
	// =====================================================

	if body[0] == '[' {

		var students []model.Student

		if err := json.Unmarshal(body, &students); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Check empty array
		if len(students) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "student list cannot be empty",
			})
			return
		}

		createdStudents := make([]*model.Student, 0, len(students))

		for _, student := range students {

			// Validate student
			if err := validateStudent(student); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			// Create student
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

			createdStudents = append(
				createdStudents,
				createdStudent,
			)
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":  "students created successfully",
			"students": createdStudents,
		})

		return
	}

	// =====================================================
	// INVALID REQUEST
	// =====================================================

	c.JSON(http.StatusBadRequest, gin.H{
		"error": "request body must be a student object or an array of students",
	})
}

// =====================================================
// VALIDATE STUDENT
// =====================================================

func validateStudent(student model.Student) error {

	if student.Name == "" {
		return fmt.Errorf("student name is required")
	}

	if student.Age <= 0 {
		return fmt.Errorf("age must be greater than 0")
	}

	if student.RollNumber <= 0 {
		return fmt.Errorf("roll number must be greater than 0")
	}

	if student.ClassroomID == "" {
		return fmt.Errorf("classroom_id is required")
	}

	return nil
}

// =====================================================
// GET ALL STUDENTS
// GET /api/v1/students
// =====================================================

func (h *StudentHandler) GetAll(c *gin.Context) {

	students, err := h.service.GetAll(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, students)
}

// =====================================================
// GET STUDENT BY ID
// GET /api/v1/students/:id
// =====================================================

func (h *StudentHandler) GetByID(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(
		c.Param("id"),
	)

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

// =====================================================
// UPDATE STUDENT
// PUT /api/v1/students/:id
// =====================================================

func (h *StudentHandler) Update(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(
		c.Param("id"),
	)

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

	// Validate student
	if err := validateStudent(student); err != nil {
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

// =====================================================
// DELETE STUDENT
// DELETE /api/v1/students/:id
// =====================================================

func (h *StudentHandler) Delete(c *gin.Context) {

	id, err := bson.ObjectIDFromHex(
		c.Param("id"),
	)

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
