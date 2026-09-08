package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shivkumar7dandin-gif/students-api/internal/assessment/model"
	"github.com/shivkumar7dandin-gif/students-api/internal/assessment/repository"
	calendarService "github.com/shivkumar7dandin-gif/students-api/internal/calendar/service"
	classroomService "github.com/shivkumar7dandin-gif/students-api/internal/classroom/service"
	studentService "github.com/shivkumar7dandin-gif/students-api/internal/student/service"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AssessmentService struct {
	repo             *repository.AssessmentRepository
	calendarService  *calendarService.CalendarService
	studentService   *studentService.StudentService
	classroomService *classroomService.ClassroomService
}

func NewAssessmentService(
	repo *repository.AssessmentRepository,
	calendarService *calendarService.CalendarService,
	studentService *studentService.StudentService,
	classroomService *classroomService.ClassroomService,
) *AssessmentService {

	return &AssessmentService{
		repo:             repo,
		calendarService:  calendarService,
		studentService:   studentService,
		classroomService: classroomService,
	}
}

// ========================================
// CREATE
// ========================================

func (s *AssessmentService) Create(
	ctx context.Context,
	collegeID bson.ObjectID,
	assessment model.Assessment,
) (*model.Assessment, error) {

	// College must come from JWT.
	assessment.CollegeID = collegeID

	// ========================================
	// 1. Academic year
	// ========================================

	assessment.AcademicYear =
		strings.TrimSpace(assessment.AcademicYear)

	if assessment.AcademicYear == "" {
		return nil, fmt.Errorf(
			"academic_year is required",
		)
	}

	// ========================================
	// 2. Student
	// ========================================

	if assessment.StudentID.IsZero() {
		return nil, fmt.Errorf(
			"student_id is required",
		)
	}

	if s.studentService == nil {
		return nil, fmt.Errorf(
			"student service is not configured",
		)
	}

	student, err :=
		s.studentService.GetByIDAndCollegeID(
			ctx,
			assessment.StudentID,
			collegeID,
		)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf(
				"student not found in this college",
			)
		}

		return nil, fmt.Errorf(
			"failed to validate student: %w",
			err,
		)
	}

	// ========================================
	// 3. Classroom
	// ========================================

	assessment.ClassroomID =
		strings.TrimSpace(assessment.ClassroomID)

	if assessment.ClassroomID == "" {
		return nil, fmt.Errorf(
			"classroom_id is required",
		)
	}

	classroomID, err :=
		bson.ObjectIDFromHex(
			assessment.ClassroomID,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid classroom_id",
		)
	}

	if s.classroomService == nil {
		return nil, fmt.Errorf(
			"classroom service is not configured",
		)
	}

	_, err =
		s.classroomService.GetByIDAndCollegeID(
			ctx,
			classroomID,
			collegeID,
		)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf(
				"classroom not found in this college",
			)
		}

		return nil, fmt.Errorf(
			"failed to validate classroom: %w",
			err,
		)
	}

	// Student must actually belong to classroom.
	if student.ClassroomID != assessment.ClassroomID {
		return nil, fmt.Errorf(
			"student does not belong to this classroom",
		)
	}

	// ========================================
	// 4. Assessment date
	// ========================================

	assessment.AssessmentDate =
		strings.TrimSpace(
			assessment.AssessmentDate,
		)

	if assessment.AssessmentDate == "" {
		return nil, fmt.Errorf(
			"assessment_date is required",
		)
	}

	assessmentDate, err := time.Parse(
		"2006-01-02",
		assessment.AssessmentDate,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"assessment_date must be in YYYY-MM-DD format",
		)
	}

	// ========================================
	// 5. Validate academic calendar
	// ========================================

	if s.calendarService == nil {
		return nil, fmt.Errorf(
			"calendar service is not configured",
		)
	}

	dayStatus, err :=
		s.calendarService.GetDayStatus(
			ctx,
			collegeID,
			assessment.AcademicYear,
			assessmentDate,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to validate academic calendar: %w",
			err,
		)
	}

	if !dayStatus.IsWorkingDay {
		return nil, fmt.Errorf(
			"assessment cannot be created on %s: %s",
			assessment.AssessmentDate,
			dayStatus.Reason,
		)
	}

	// ========================================
	// 6. Assessment number
	// ========================================

	if assessment.AssessmentNo <= 0 {
		return nil, fmt.Errorf(
			"assessment_no must be greater than 0",
		)
	}

	// ========================================
	// 7. Subject
	// ========================================

	assessment.Subject =
		strings.TrimSpace(assessment.Subject)

	if assessment.Subject == "" {
		return nil, fmt.Errorf(
			"subject is required",
		)
	}

	// ========================================
	// 8. Marks
	// ========================================

	if assessment.TotalMarks <= 0 {
		return nil, fmt.Errorf(
			"total_marks must be greater than 0",
		)
	}

	if assessment.ObtainedMarks < 0 {
		return nil, fmt.Errorf(
			"obtained_marks cannot be negative",
		)
	}

	if assessment.ObtainedMarks >
		assessment.TotalMarks {

		return nil, fmt.Errorf(
			"obtained_marks cannot be greater than total_marks",
		)
	}

	// ========================================
	// 9. Percentage
	// ========================================

	assessment.Percentage =
		(assessment.ObtainedMarks /
			assessment.TotalMarks) * 100

	// ========================================
	// 10. Result
	// ========================================

	if assessment.Percentage >= 35 {
		assessment.Result = "Pass"
	} else {
		assessment.Result = "Fail"
	}

	// ========================================
	// 11. Create
	// ========================================

	return s.repo.Create(
		ctx,
		assessment,
	)
}

// ========================================
// GET ALL
// ========================================

func (s *AssessmentService) GetAll(
	ctx context.Context,
	collegeID bson.ObjectID,
	academicYear string,
) ([]model.Assessment, error) {

	academicYear =
		strings.TrimSpace(academicYear)

	if academicYear == "" {
		return nil, fmt.Errorf(
			"academic_year is required",
		)
	}

	return s.repo.GetAllByCollegeAndYear(
		ctx,
		collegeID,
		academicYear,
	)
}

// ========================================
// GET BY STUDENT
// ========================================

func (s *AssessmentService) GetByStudent(
	ctx context.Context,
	collegeID bson.ObjectID,
	studentID bson.ObjectID,
	academicYear string,
) ([]model.Assessment, error) {

	academicYear =
		strings.TrimSpace(academicYear)

	if academicYear == "" {
		return nil, fmt.Errorf(
			"academic_year is required",
		)
	}

	if studentID.IsZero() {
		return nil, fmt.Errorf(
			"student_id is required",
		)
	}

	return s.repo.GetByStudentAndYear(
		ctx,
		collegeID,
		studentID,
		academicYear,
	)
}

func (s *AssessmentService) GetMonthlySummary(
	ctx context.Context,
	collegeID bson.ObjectID,
	studentID bson.ObjectID,
	academicYear string,
	year int,
	month int,
) (*model.MonthlyAssessmentSummary, error) {

	academicYear = strings.TrimSpace(academicYear)

	if academicYear == "" {
		return nil, fmt.Errorf("academic_year is required")
	}

	if studentID.IsZero() {
		return nil, fmt.Errorf("student_id is required")
	}

	if year <= 0 {
		return nil, fmt.Errorf("year must be greater than 0")
	}

	if month < 1 || month > 12 {
		return nil, fmt.Errorf("month must be between 1 and 12")
	}

	assessments, err := s.repo.GetByStudentAndMonth(
		ctx,
		collegeID,
		studentID,
		academicYear,
		year,
		month,
	)

	if err != nil {
		return nil, err
	}

	summary := &model.MonthlyAssessmentSummary{
		StudentID:    studentID,
		AcademicYear: academicYear,
		Year:         year,
		Month:        month,
	}

	for _, assessment := range assessments {

		summary.TotalAssessments++

		summary.TotalMarks += assessment.TotalMarks
		summary.ObtainedMarks += assessment.ObtainedMarks

		if strings.EqualFold(assessment.Result, "Pass") {
			summary.Passed++
		} else if strings.EqualFold(assessment.Result, "Fail") {
			summary.Failed++
		}
	}

	if summary.TotalMarks > 0 {

		summary.Percentage =
			(summary.ObtainedMarks /
				summary.TotalMarks) * 100

		if summary.Percentage >= 35 {
			summary.Result = "Pass"
		} else {
			summary.Result = "Fail"
		}

	} else {
		summary.Percentage = 0
		summary.Result = "No Assessment"
	}

	return summary, nil
}
