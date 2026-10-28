package handler

import (
	"net/http"

	subjectModel "github.com/shivkumar7dandin-gif/students-api/internal/subject/model"
	subjectService "github.com/shivkumar7dandin-gif/students-api/internal/subject/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SubjectHandler struct {
	service *subjectService.SubjectService
}

func NewSubjectHandler(
	service *subjectService.SubjectService,
) *SubjectHandler {
	return &SubjectHandler{
		service: service,
	}
}

func getCollegeID(c *gin.Context) (bson.ObjectID, error) {
	collegeIDValue, exists := c.Get("college_id")

	if !exists {
		return bson.NilObjectID, &collegeIDError{
			message: "college ID not found in token",
		}
	}

	collegeIDString, ok := collegeIDValue.(string)

	if !ok || collegeIDString == "" {
		return bson.NilObjectID, &collegeIDError{
			message: "invalid college ID in token",
		}
	}

	collegeID, err := bson.ObjectIDFromHex(collegeIDString)

	if err != nil {
		return bson.NilObjectID, &collegeIDError{
			message: "invalid college ID",
		}
	}

	return collegeID, nil
}

type collegeIDError struct {
	message string
}

func (e *collegeIDError) Error() string {
	return e.message
}

func (h *SubjectHandler) Create(c *gin.Context) {
	collegeID, err := getCollegeID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req subjectModel.CreateSubjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	subject, err := h.service.Create(
		c.Request.Context(),
		collegeID,
		req,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, subject)
}

func (h *SubjectHandler) GetAll(c *gin.Context) {
	collegeID, err := getCollegeID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	subjects, err := h.service.GetAll(
		c.Request.Context(),
		collegeID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get subjects",
		})
		return
	}

	c.JSON(http.StatusOK, subjects)
}

func (h *SubjectHandler) GetByID(c *gin.Context) {
	collegeID, err := getCollegeID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	subjectID, err := bson.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid subject ID",
		})
		return
	}

	subject, err := h.service.GetByID(
		c.Request.Context(),
		subjectID,
		collegeID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, subject)
}

func (h *SubjectHandler) Delete(c *gin.Context) {
	collegeID, err := getCollegeID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	subjectID, err := bson.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid subject ID",
		})
		return
	}

	err = h.service.Delete(
		c.Request.Context(),
		subjectID,
		collegeID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "subject deleted successfully",
	})
}
