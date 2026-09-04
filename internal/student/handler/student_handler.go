package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	studentModel "github.com/shivkumar7dandin-gif/students-api/internal/student/model"
	studentService "github.com/shivkumar7dandin-gif/students-api/internal/student/service"

	assessmentService "github.com/shivkumar7dandin-gif/students-api/internal/assessment/service"
	attendanceService "github.com/shivkumar7dandin-gif/students-api/internal/attendance/service"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type StudentHandler struct {
	service           *studentService.StudentService
	attendanceService *attendanceService.AttendanceService
	assessmentService *assessmentService.AssessmentService
}

func NewStudentHandler(
	service *studentService.StudentService,
	attendanceService *attendanceService.AttendanceService,
	assessmentService *assessmentService.AssessmentService,
) *StudentHandler {

	return &StudentHandler{
		service:           service,
		attendanceService: attendanceService,
		assessmentService: assessmentService,
	}
}

// =====================================================
// POST /api/v1/students
// Supports:
// 1. Single student
// 2. Bulk students
// =====================================================

func (h *StudentHandler) Create(c *gin.Context) {

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

		var req studentModel.CreateStudentRequest

		if err := json.Unmarshal(
			body,
			&req,
		); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := validateCreateStudent(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		createdStudent, err := h.service.Create(
			c.Request.Context(),
			req,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
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

		var requests []studentModel.CreateStudentRequest

		if err := json.Unmarshal(
			body,
			&requests,
		); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if len(requests) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "student list cannot be empty",
			})
			return
		}

		createdStudents := make(
			[]*studentModel.Student,
			0,
			len(requests),
		)

		for index, req := range requests {

			if err := validateCreateStudent(req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf(
						"student at index %d: %s",
						index,
						err.Error(),
					),
				})
				return
			}

			createdStudent, err := h.service.Create(
				c.Request.Context(),
				req,
			)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf(
						"student at index %d: %s",
						index,
						err.Error(),
					),
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

	c.JSON(http.StatusBadRequest, gin.H{
		"error": "request body must be a student object or an array of students",
	})
}

// =====================================================
// VALIDATE CREATE STUDENT
// =====================================================

func validateCreateStudent(
	req studentModel.CreateStudentRequest,
) error {

	if req.Name == "" {
		return fmt.Errorf("student name is required")
	}

	if req.Age <= 0 {
		return fmt.Errorf("age must be greater than 0")
	}

	if req.RollNumber <= 0 {
		return fmt.Errorf(
			"roll number must be greater than 0",
		)
	}

	if req.ClassroomID == "" {
		return fmt.Errorf("classroom_id is required")
	}

	if req.Username == "" {
		return fmt.Errorf("username is required")
	}

	if req.Password == "" {
		return fmt.Errorf("password is required")
	}

	if len(req.Password) < 8 {
		return fmt.Errorf(
			"password must contain at least 8 characters",
		)
	}

	return nil
}

// =====================================================
// VALIDATE STUDENT FOR UPDATE
// =====================================================

func validateStudent(
	student studentModel.Student,
) error {

	if student.Name == "" {
		return fmt.Errorf("student name is required")
	}

	if student.Age <= 0 {
		return fmt.Errorf("age must be greater than 0")
	}

	if student.RollNumber <= 0 {
		return fmt.Errorf(
			"roll number must be greater than 0",
		)
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

	referenceIDValue, exists := c.Get("reference_id")
	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "college reference id not found",
			},
		)
		return
	}

	referenceID, ok := referenceIDValue.(string)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college reference id",
			},
		)
		return
	}

	collegeID, err := bson.ObjectIDFromHex(referenceID)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid college id",
			},
		)
		return
	}

	students, err := h.service.GetByCollegeID(
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
		students,
	)
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

	c.JSON(
		http.StatusOK,
		student,
	)
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

	var student studentModel.Student

	if err := c.ShouldBindJSON(
		&student,
	); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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

// =====================================================
// GET LOGGED-IN STUDENT DETAILS
// GET /api/v1/students/me
// =====================================================

func (h *StudentHandler) GetMe(c *gin.Context) {

	// ------------------------------------------
	// 1. Get role from JWT context
	// ------------------------------------------

	roleValue, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "role not found",
		})
		return
	}

	role, ok := roleValue.(string)
	if !ok || role != "student" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "only students can access this API",
		})
		return
	}

	// ------------------------------------------
	// 2. Get reference_id from JWT
	// ------------------------------------------

	referenceIDValue, exists := c.Get("reference_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "reference_id not found",
		})
		return
	}

	referenceIDString, ok := referenceIDValue.(string)
	if !ok || referenceIDString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid reference_id",
		})
		return
	}

	// ------------------------------------------
	// 3. Convert reference_id to ObjectID
	// ------------------------------------------

	studentID, err := bson.ObjectIDFromHex(
		referenceIDString,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid student reference id",
		})
		return
	}

	// ------------------------------------------
	// 4. Get student profile
	// ------------------------------------------

	student, err := h.service.GetByID(
		c.Request.Context(),
		studentID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "student not found",
		})
		return
	}

	// ------------------------------------------
	// 5. Get attendance summary
	// ------------------------------------------

	attendanceSummary, err := h.attendanceService.GetSummary(
		c.Request.Context(),
		studentID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get attendance summary",
		})
		return
	}

	// ------------------------------------------
	// 6. Get assessments
	// ------------------------------------------

	assessments, err := h.assessmentService.GetByStudent(
		c.Request.Context(),
		studentID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get student assessments",
		})
		return
	}

	// ------------------------------------------
	// 7. Combined response
	// ------------------------------------------

	c.JSON(http.StatusOK, gin.H{
		"student":            student,
		"attendance_summary": attendanceSummary,
		"assessments":        assessments,
	})
}
