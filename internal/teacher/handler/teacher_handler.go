package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	teacherModel "github.com/shivkumar7dandin-gif/students-api/internal/teacher/model"
	teacherService "github.com/shivkumar7dandin-gif/students-api/internal/teacher/service"
)

type TeacherHandler struct {
	service *teacherService.TeacherService
}

func NewTeacherHandler(
	service *teacherService.TeacherService,
) *TeacherHandler {

	return &TeacherHandler{
		service: service,
	}
}

// =====================================================
// CREATE TEACHER
// POST /api/v1/teachers
// =====================================================

func (h *TeacherHandler) Create(c *gin.Context) {

	var req teacherModel.CreateTeacherRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request body",
			},
		)
		return
	}

	collegeIDValue, exists := c.Get("college_id")
	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "college_id not found",
			},
		)
		return
	}

	collegeID, ok := collegeIDValue.(string)
	if !ok || collegeID == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id",
			},
		)
		return
	}

	req.CollegeID = collegeID

	teacher, err := h.service.Create(
		c.Request.Context(),
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
		teacher,
	)
}

// =====================================================
// GET ALL TEACHERS
// GET /api/v1/teachers
// =====================================================

func (h *TeacherHandler) GetAll(c *gin.Context) {

	collegeIDValue, exists := c.Get(
		"college_id",
	)

	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "college_id not found",
			},
		)
		return
	}

	collegeIDString, ok :=
		collegeIDValue.(string)

	if !ok || collegeIDString == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id",
			},
		)
		return
	}

	collegeID, err :=
		bson.ObjectIDFromHex(
			collegeIDString,
		)

	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college id",
			},
		)
		return
	}

	teachers, err :=
		h.service.GetByCollegeID(
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
		teachers,
	)
}

// =====================================================
// GET TEACHER BY ID
// GET /api/v1/teachers/:id
// =====================================================

func (h *TeacherHandler) GetByID(c *gin.Context) {

	// Get college_id from JWT
	collegeIDValue, exists := c.Get("college_id")
	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "college_id not found",
			},
		)
		return
	}

	collegeIDString, ok := collegeIDValue.(string)
	if !ok || collegeIDString == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id",
			},
		)
		return
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
		return
	}

	// Get teacher only from same college
	teacher, err := h.service.GetByIDAndCollegeID(
		c.Request.Context(),
		c.Param("id"),
		collegeID,
	)

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "teacher not found",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		teacher,
	)
}

// =====================================================
// UPDATE TEACHER
// PUT /api/v1/teachers/:id
// =====================================================

func (h *TeacherHandler) Update(c *gin.Context) {

	// Get college_id from JWT
	collegeIDValue, exists := c.Get("college_id")

	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "college_id not found",
			},
		)
		return
	}

	collegeIDString, ok := collegeIDValue.(string)

	if !ok || collegeIDString == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id",
			},
		)
		return
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
		return
	}

	// Read teacher data
	var teacher teacherModel.Teacher

	if err := c.ShouldBindJSON(&teacher); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request body",
			},
		)
		return
	}

	// Tenant-safe update
	err = h.service.UpdateByCollegeID(
		c.Request.Context(),
		c.Param("id"),
		collegeID,
		teacher,
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
			"message": "teacher updated successfully",
		},
	)
}

// =====================================================
// UPDATE TEACHER STATUS
// PATCH /api/v1/teachers/:id/status
// =====================================================

func (h *TeacherHandler) UpdateStatus(c *gin.Context) {

	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid request body",
			},
		)
		return
	}

	// Get college_id from JWT
	collegeIDValue, exists := c.Get("college_id")

	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "college_id not found",
			},
		)
		return
	}

	collegeIDString, ok := collegeIDValue.(string)

	if !ok || collegeIDString == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id",
			},
		)
		return
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
		return
	}

	// Tenant-safe status update
	err = h.service.UpdateStatusByCollegeID(
		c.Request.Context(),
		c.Param("id"),
		collegeID,
		req.Status,
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
			"message": "teacher status updated successfully",
		},
	)
}

// =====================================================
// DELETE TEACHER
// DELETE /api/v1/teachers/:id
// =====================================================

func (h *TeacherHandler) Delete(c *gin.Context) {

	// Get college_id from JWT
	collegeIDValue, exists := c.Get("college_id")

	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "college_id not found",
			},
		)
		return
	}

	collegeIDString, ok := collegeIDValue.(string)

	if !ok || collegeIDString == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid college_id",
			},
		)
		return
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
		return
	}

	// Tenant-safe delete
	err = h.service.DeleteByCollegeID(
		c.Request.Context(),
		c.Param("id"),
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
		http.StatusOK,
		gin.H{
			"message": "teacher deleted successfully",
		},
	)
}

// =====================================================
// GET LOGGED-IN TEACHER DETAILS
// GET /api/v1/teachers/me
// =====================================================

func (h *TeacherHandler) GetMe(c *gin.Context) {

	referenceIDValue, exists := c.Get(
		"reference_id",
	)

	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "reference_id not found",
			},
		)
		return
	}

	referenceID, ok :=
		referenceIDValue.(string)

	if !ok || referenceID == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid reference_id",
			},
		)
		return
	}

	teacher, err := h.service.GetByID(
		c.Request.Context(),
		referenceID,
	)

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "teacher not found",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"teacher": teacher,
		},
	)
}
