package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shivkumar7dandin-gif/students-api/internal/attendance/model"
	"github.com/shivkumar7dandin-gif/students-api/internal/attendance/repository"
	calendarService "github.com/shivkumar7dandin-gif/students-api/internal/calendar/service"
	classroomService "github.com/shivkumar7dandin-gif/students-api/internal/classroom/service"
	studentService "github.com/shivkumar7dandin-gif/students-api/internal/student/service"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AttendanceService struct {
	repo             *repository.AttendanceRepository
	calendarService  *calendarService.CalendarService
	studentService   *studentService.StudentService
	classroomService *classroomService.ClassroomService
}

func NewAttendanceService(
	repo *repository.AttendanceRepository,
	calendarService *calendarService.CalendarService,
	studentService *studentService.StudentService,
	classroomService *classroomService.ClassroomService,
) *AttendanceService {

	return &AttendanceService{
		repo:             repo,
		calendarService:  calendarService,
		studentService:   studentService,
		classroomService: classroomService,
	}
}

// ========================================
// CREATE ATTENDANCE
// ========================================

func (s *AttendanceService) Create(
	ctx context.Context,
	collegeID bson.ObjectID,
	attendance model.Attendance,
) (*model.Attendance, error) {

	// College ID must always come from JWT.
	attendance.CollegeID = collegeID

	// ========================================
	// 1. Validate academic year
	// ========================================

	attendance.AcademicYear =
		strings.TrimSpace(attendance.AcademicYear)

	if attendance.AcademicYear == "" {
		return nil, fmt.Errorf(
			"academic_year is required",
		)
	}

	// ========================================
	// 2. Validate student ID
	// ========================================

	if attendance.StudentID.IsZero() {
		return nil, fmt.Errorf(
			"student_id is required",
		)
	}

	// ========================================
	// 3. Verify student belongs to college
	// ========================================

	if s.studentService == nil {
		return nil, fmt.Errorf(
			"student service is not configured",
		)
	}

	student, err :=
		s.studentService.GetByIDAndCollegeID(
			ctx,
			attendance.StudentID,
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
	// 4. Validate classroom ID
	// ========================================

	attendance.ClassroomID =
		strings.TrimSpace(attendance.ClassroomID)

	if attendance.ClassroomID == "" {
		return nil, fmt.Errorf(
			"classroom_id is required",
		)
	}

	classroomID, err :=
		bson.ObjectIDFromHex(
			attendance.ClassroomID,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid classroom_id",
		)
	}

	// ========================================
	// 5. Verify classroom belongs to college
	// ========================================

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

	// ========================================
	// 6. Verify student belongs to classroom
	// ========================================

	if student.ClassroomID != attendance.ClassroomID {
		return nil, fmt.Errorf(
			"student does not belong to this classroom",
		)
	}

	// ========================================
	// 7. Validate attendance status
	// ========================================

	status := strings.ToLower(
		strings.TrimSpace(attendance.Attendance),
	)

	switch status {

	case "present":
		attendance.Attendance = "Present"

	case "absent":
		attendance.Attendance = "Absent"

	default:
		return nil, fmt.Errorf(
			"attendance must be Present or Absent",
		)
	}

	// ========================================
	// 8. Validate date
	// ========================================

	attendance.Date =
		strings.TrimSpace(attendance.Date)

	if attendance.Date == "" {
		return nil, fmt.Errorf(
			"date is required",
		)
	}

	parsedDate, err := time.Parse(
		"2006-01-02",
		attendance.Date,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"date must be in YYYY-MM-DD format",
		)
	}

	// ========================================
	// 9. Validate academic calendar
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
			attendance.AcademicYear,
			parsedDate,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to validate academic calendar: %w",
			err,
		)
	}

	if !dayStatus.IsWorkingDay {
		return nil, fmt.Errorf(
			"attendance cannot be marked on %s: %s",
			attendance.Date,
			dayStatus.Reason,
		)
	}

	// ========================================
	// 10. Check duplicate attendance
	// ========================================

	existingAttendance, err :=
		s.repo.GetByStudentAndDate(
			ctx,
			collegeID,
			attendance.StudentID,
			attendance.AcademicYear,
			attendance.Date,
		)

	if err == nil && existingAttendance != nil {

		return nil, fmt.Errorf(
			"attendance already exists for this student on %s",
			attendance.Date,
		)
	}

	if err != nil &&
		err != mongo.ErrNoDocuments {

		return nil, fmt.Errorf(
			"failed to check existing attendance: %w",
			err,
		)
	}

	// ========================================
	// 11. Create attendance
	// ========================================

	return s.repo.Create(
		ctx,
		attendance,
	)
}

// ========================================
// GET BY STUDENT
// ========================================

func (s *AttendanceService) GetByStudent(
	ctx context.Context,
	collegeID bson.ObjectID,
	studentID bson.ObjectID,
	academicYear string,
) ([]model.Attendance, error) {

	academicYear =
		strings.TrimSpace(academicYear)

	if academicYear == "" {
		return nil, fmt.Errorf(
			"academic_year is required",
		)
	}

	return s.repo.GetByStudent(
		ctx,
		collegeID,
		studentID,
		academicYear,
	)
}

// ========================================
// GET ALL
// ========================================

func (s *AttendanceService) GetAll(
	ctx context.Context,
	collegeID bson.ObjectID,
	academicYear string,
) ([]model.Attendance, error) {

	academicYear =
		strings.TrimSpace(academicYear)

	if academicYear == "" {
		return nil, fmt.Errorf(
			"academic_year is required",
		)
	}

	return s.repo.GetAll(
		ctx,
		collegeID,
		academicYear,
	)
}

// ========================================
// SUMMARY
// ========================================

func (s *AttendanceService) GetSummary(
	ctx context.Context,
	collegeID bson.ObjectID,
	studentID bson.ObjectID,
	academicYear string,
) (*model.AttendanceSummary, error) {

	academicYear =
		strings.TrimSpace(academicYear)

	if academicYear == "" {
		return nil, fmt.Errorf(
			"academic_year is required",
		)
	}

	attendanceList, err :=
		s.repo.GetByStudent(
			ctx,
			collegeID,
			studentID,
			academicYear,
		)

	if err != nil {
		return nil, err
	}

	totalDays := len(attendanceList)

	present := 0
	absent := 0

	for _, attendance := range attendanceList {

		switch strings.ToLower(
			attendance.Attendance,
		) {

		case "present":
			present++

		case "absent":
			absent++
		}
	}

	percentage := 0.0

	if totalDays > 0 {

		percentage =
			(float64(present) /
				float64(totalDays)) * 100
	}

	return &model.AttendanceSummary{
		StudentID:            studentID,
		AcademicYear:         academicYear,
		TotalDays:            totalDays,
		Present:              present,
		Absent:               absent,
		AttendancePercentage: percentage,
	}, nil
}

// ========================================
// MONTHLY SUMMARY
// ========================================

func (s *AttendanceService) GetMonthlySummary(
	ctx context.Context,
	collegeID bson.ObjectID,
	studentID bson.ObjectID,
	academicYear string,
	year int,
	month int,
) (*model.MonthlyAttendanceSummary, error) {

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

	if year <= 0 {
		return nil, fmt.Errorf(
			"valid year is required",
		)
	}

	if month < 1 || month > 12 {
		return nil, fmt.Errorf(
			"month must be between 1 and 12",
		)
	}

	// ========================================
	// 1. Get monthly attendance records
	// ========================================

	attendanceList, err :=
		s.repo.GetByStudentAndMonth(
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

	attendanceMarkedDays :=
		len(attendanceList)

	present := 0
	absent := 0

	for _, attendance := range attendanceList {

		switch strings.ToLower(
			attendance.Attendance,
		) {

		case "present":
			present++

		case "absent":
			absent++
		}
	}

	// ========================================
	// 2. Get real working days from calendar
	// ========================================

	workingDays, err :=
		s.calendarService.GetMonthlyWorkingDays(
			ctx,
			collegeID,
			academicYear,
			year,
			month,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to calculate working days: %w",
			err,
		)
	}

	// ========================================
	// 3. Calculate percentage
	// ========================================

	percentage := 0.0

	if attendanceMarkedDays > 0 {

		percentage =
			(float64(present) /
				float64(attendanceMarkedDays)) * 100
	}

	return &model.MonthlyAttendanceSummary{
		StudentID:            studentID,
		AcademicYear:         academicYear,
		Year:                 year,
		Month:                month,
		WorkingDays:          workingDays,
		AttendanceMarkedDays: attendanceMarkedDays,
		Present:              present,
		Absent:               absent,
		AttendancePercentage: percentage,
	}, nil
}
