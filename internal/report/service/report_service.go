package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	assessmentService "github.com/shivkumar7dandin-gif/students-api/internal/assessment/service"
	attendanceService "github.com/shivkumar7dandin-gif/students-api/internal/attendance/service"
	reportModel "github.com/shivkumar7dandin-gif/students-api/internal/report/model"
	studentService "github.com/shivkumar7dandin-gif/students-api/internal/student/service"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ReportService struct {
	studentService    *studentService.StudentService
	attendanceService *attendanceService.AttendanceService
	assessmentService *assessmentService.AssessmentService
}

func NewReportService(
	studentService *studentService.StudentService,
	attendanceService *attendanceService.AttendanceService,
	assessmentService *assessmentService.AssessmentService,
) *ReportService {

	return &ReportService{
		studentService:    studentService,
		attendanceService: attendanceService,
		assessmentService: assessmentService,
	}
}

func (s *ReportService) GetMonthlyParentReport(
	ctx context.Context,
	collegeID bson.ObjectID,
	studentID bson.ObjectID,
	academicYear string,
	year int,
	month int,
) (*reportModel.MonthlyParentReport, error) {

	academicYear = strings.TrimSpace(academicYear)

	if collegeID.IsZero() {
		return nil, errors.New("college_id is required")
	}

	if studentID.IsZero() {
		return nil, errors.New("student_id is required")
	}

	if academicYear == "" {
		return nil, errors.New("academic_year is required")
	}

	if year <= 0 {
		return nil, errors.New("valid year is required")
	}

	if month < 1 || month > 12 {
		return nil, errors.New(
			"month must be between 1 and 12",
		)
	}

	// ---------------------------------------
	// 1. Get student using tenant-safe query
	// ---------------------------------------

	student, err :=
		s.studentService.GetByIDAndCollegeID(
			ctx,
			studentID,
			collegeID,
		)

	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(
				"student not found in this college",
			)
		}

		return nil, fmt.Errorf(
			"failed to get student: %w",
			err,
		)
	}

	// ---------------------------------------
	// 2. Find guardians who receive reports
	// ---------------------------------------

	recipients :=
		make([]reportModel.ReportRecipient, 0)

	for _, guardian := range student.Guardians {

		if !guardian.ReceiveReport {
			continue
		}

		recipients = append(
			recipients,
			reportModel.ReportRecipient{
				Name:     guardian.Name,
				Relation: guardian.Relation,
				Phone:    guardian.Phone,
				Email:    guardian.Email,
				Primary:  guardian.Primary,
			},
		)
	}

	if len(recipients) == 0 {
		return nil, errors.New(
			"no guardian configured to receive reports",
		)
	}

	// ---------------------------------------
	// 3. Monthly attendance summary
	// ---------------------------------------

	attendanceSummary, err :=
		s.attendanceService.GetMonthlySummary(
			ctx,
			collegeID,
			studentID,
			academicYear,
			year,
			month,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get attendance summary: %w",
			err,
		)
	}

	// ---------------------------------------
	// 4. Monthly assessment summary
	// ---------------------------------------

	assessmentSummary, err :=
		s.assessmentService.GetMonthlySummary(
			ctx,
			collegeID,
			studentID,
			academicYear,
			year,
			month,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get assessment summary: %w",
			err,
		)
	}

	// ---------------------------------------
	// 5. Build report
	// ---------------------------------------

	report := &reportModel.MonthlyParentReport{
		StudentID:        student.ID,
		StudentName:      student.Name,
		RollNumber:       student.RollNumber,
		ClassroomID:      student.ClassroomID,
		AcademicYear:     academicYear,
		Year:             year,
		Month:            month,
		ReportRecipients: recipients,
		Attendance:       attendanceSummary,
		Assessment:       assessmentSummary,
		Address:          student.Address,
	}

	return report, nil
}
