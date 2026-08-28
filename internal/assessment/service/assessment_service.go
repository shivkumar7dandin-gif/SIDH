package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/shivkumar7dandin-gif/students-api/internal/assessment/model"
	"github.com/shivkumar7dandin-gif/students-api/internal/assessment/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AssessmentService struct {
	repo *repository.AssessmentRepository
}

func NewAssessmentService(
	repo *repository.AssessmentRepository,
) *AssessmentService {

	return &AssessmentService{
		repo: repo,
	}
}

func (s *AssessmentService) Create(
	ctx context.Context,
	assessment model.Assessment,
) (*model.Assessment, error) {

	if assessment.StudentID.IsZero() {
		return nil, fmt.Errorf("student_id is required")
	}

	if strings.TrimSpace(assessment.ClassroomID) == "" {
		return nil, fmt.Errorf("classroom_id is required")
	}

	if assessment.AssessmentNo <= 0 {
		return nil, fmt.Errorf("assessment_no must be greater than 0")
	}

	if strings.TrimSpace(assessment.Subject) == "" {
		return nil, fmt.Errorf("subject is required")
	}

	if assessment.TotalMarks <= 0 {
		return nil, fmt.Errorf("total_marks must be greater than 0")
	}

	if assessment.ObtainedMarks < 0 {
		return nil, fmt.Errorf("obtained_marks cannot be negative")
	}

	if assessment.ObtainedMarks > assessment.TotalMarks {
		return nil, fmt.Errorf(
			"obtained_marks cannot be greater than total_marks",
		)
	}

	assessment.Percentage =
		(assessment.ObtainedMarks / assessment.TotalMarks) * 100

	if assessment.Percentage >= 35 {
		assessment.Result = "Pass"
	} else {
		assessment.Result = "Fail"
	}

	return s.repo.Create(ctx, assessment)
}

func (s *AssessmentService) GetAll(
	ctx context.Context,
) ([]model.Assessment, error) {

	return s.repo.GetAll(ctx)
}

func (s *AssessmentService) GetByStudent(
	ctx context.Context,
	studentID bson.ObjectID,
) ([]model.Assessment, error) {

	return s.repo.GetByStudent(ctx, studentID)
}
