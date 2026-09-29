package service

import (
	"context"
	"fmt"
	"sort"
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

	if s.calendarService == nil {
		return nil, fmt.Errorf(
			"calendar service is not configured",
		)
	}

	if err := s.calendarService.ValidateActiveAcademicYear(
		ctx,
		collegeID,
		attendance.AcademicYear,
	); err != nil {
		return nil, err
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
	// Allow attendance only for today's date
	// ========================================

	indiaLocation, err :=
		time.LoadLocation("Asia/Kolkata")

	if err != nil {
		return nil, fmt.Errorf(
			"failed to load server timezone",
		)
	}

	today :=
		time.Now().
			In(indiaLocation).
			Format("2006-01-02")

	if attendance.Date != today {
		return nil, fmt.Errorf(
			"attendance can only be marked or updated for today's date",
		)
	}

	// ========================================
	// 9. Validate academic calendar
	// ========================================

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

		// Attendance already exists.
		// Update Present / Absent instead of
		// creating a duplicate record.

		return s.repo.Update(
			ctx,
			collegeID,
			existingAttendance.ID,
			attendance.Attendance,
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
// CREATE / UPDATE BULK ATTENDANCE
// ========================================

func (s *AttendanceService) CreateBulk(
	ctx context.Context,
	collegeID bson.ObjectID,
	req model.BulkAttendanceRequest,
) ([]model.Attendance, error) {

	req.ClassroomID =
		strings.TrimSpace(req.ClassroomID)

	req.AcademicYear =
		strings.TrimSpace(req.AcademicYear)

	req.Date =
		strings.TrimSpace(req.Date)

	// ========================================
	// Basic request validation
	// ========================================

	if req.ClassroomID == "" {
		return nil, fmt.Errorf(
			"classroom_id is required",
		)
	}

	if req.AcademicYear == "" {
		return nil, fmt.Errorf(
			"academic_year is required",
		)
	}

	if req.Date == "" {
		return nil, fmt.Errorf(
			"date is required",
		)
	}

	if len(req.Attendances) == 0 {
		return nil, fmt.Errorf(
			"attendance list cannot be empty",
		)
	}

	// ========================================
	// Prevent duplicate students in request
	// ========================================

	seenStudents :=
		make(map[bson.ObjectID]bool)

	for _, item := range req.Attendances {

		if item.StudentID.IsZero() {
			return nil, fmt.Errorf(
				"student_id is required",
			)
		}

		if seenStudents[item.StudentID] {
			return nil, fmt.Errorf(
				"duplicate student found in attendance request",
			)
		}

		seenStudents[item.StudentID] = true
	}

	// ========================================
	// Create / update each attendance
	//
	// Existing Create() already validates:
	// - college
	// - academic year
	// - student
	// - classroom
	// - today's date
	// - working day
	// - Present / Absent
	// - create or update
	// ========================================

	results :=
		make([]model.Attendance, 0, len(req.Attendances))

	for _, item := range req.Attendances {

		attendance := model.Attendance{
			CollegeID:    collegeID,
			StudentID:    item.StudentID,
			ClassroomID:  req.ClassroomID,
			AcademicYear: req.AcademicYear,
			Date:         req.Date,
			Attendance:   item.Attendance,
		}

		savedAttendance, err :=
			s.Create(
				ctx,
				collegeID,
				attendance,
			)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to save attendance for student %s: %w",
				item.StudentID.Hex(),
				err,
			)
		}

		results = append(
			results,
			*savedAttendance,
		)
	}

	return results, nil
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

// ========================================
// CLASSROOM MONTHLY ATTENDANCE SUMMARY
// ========================================

func (s *AttendanceService) GetClassroomMonthlySummary(
	ctx context.Context,
	collegeID bson.ObjectID,
	classroomID string,
	academicYear string,
	year int,
	month int,
) (*model.ClassroomMonthlyAttendanceSummary, error) {

	// ========================================
	// 1. BASIC VALIDATION
	// ========================================

	classroomID = strings.TrimSpace(classroomID)
	academicYear = strings.TrimSpace(academicYear)

	if classroomID == "" {
		return nil, fmt.Errorf(
			"classroom_id is required",
		)
	}

	if academicYear == "" {
		return nil, fmt.Errorf(
			"academic_year is required",
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
	// 2. VALIDATE CLASSROOM ID
	// ========================================

	classroomObjectID, err :=
		bson.ObjectIDFromHex(classroomID)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid classroom_id",
		)
	}

	// ========================================
	// 3. VERIFY CLASSROOM BELONGS TO COLLEGE
	// ========================================

	if s.classroomService == nil {
		return nil, fmt.Errorf(
			"classroom service is not configured",
		)
	}

	_, err =
		s.classroomService.GetByIDAndCollegeID(
			ctx,
			classroomObjectID,
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
	// 4. GET ALL STUDENTS FROM THIS COLLEGE
	// ========================================

	if s.studentService == nil {
		return nil, fmt.Errorf(
			"student service is not configured",
		)
	}

	collegeStudents, err :=
		s.studentService.GetByCollegeID(
			ctx,
			collegeID,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get students: %w",
			err,
		)
	}

	// Used to resolve names and roll numbers.
	studentLookup :=
		make(
			map[bson.ObjectID]int,
		)

	// Final student summary map.
	studentSummary :=
		make(
			map[bson.ObjectID]*model.ClassroomStudentAttendanceSummary,
		)

	// ========================================
	// 5. ADD CURRENT CLASSROOM STUDENTS
	// ========================================

	for index, student := range collegeStudents {

		studentLookup[student.ID] =
			index

		if student.ClassroomID !=
			classroomID {
			continue
		}

		studentSummary[student.ID] =
			&model.ClassroomStudentAttendanceSummary{
				StudentID:  student.ID,
				RollNumber: student.RollNumber,
				Name:       student.Name,
			}
	}

	// ========================================
	// 6. GET MONTH ATTENDANCE
	// ========================================

	attendanceList, err :=
		s.repo.GetByClassroomAndMonth(
			ctx,
			collegeID,
			classroomID,
			academicYear,
			year,
			month,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get classroom monthly attendance: %w",
			err,
		)
	}

	// ========================================
	// 7. CALCULATE STUDENT + CLASS TOTALS
	// ========================================

	present := 0
	absent := 0

	// We need UNIQUE dates here.
	//
	// Example:
	// 50 students marked on 01-09-2026
	// means 1 marked school day,
	// NOT 50 marked days.
	markedDates :=
		make(map[string]bool)

	for _, attendance := range attendanceList {

		markedDates[attendance.Date] = true

		summary, exists := studentSummary[attendance.StudentID]

		if !exists {

			index, found := studentLookup[attendance.StudentID]

			if !found {
				continue
			}

			student := collegeStudents[index]

			summary = &model.ClassroomStudentAttendanceSummary{
				StudentID:  student.ID,
				RollNumber: student.RollNumber,
				Name:       student.Name,
			}

			studentSummary[attendance.StudentID] = summary
		}

		switch strings.ToLower(
			strings.TrimSpace(attendance.Attendance),
		) {

		case "present":
			summary.Present++
			present++

		case "absent":
			summary.Absent++
			absent++
		}

		summary.AttendanceMarkedDays++
	}

	// ========================================
	// 8. GET REAL WORKING DAYS
	// ========================================

	if s.calendarService == nil {
		return nil, fmt.Errorf(
			"calendar service is not configured",
		)
	}

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
	// 9. BUILD STUDENT LIST
	// ========================================

	students :=
		make(
			[]model.ClassroomStudentAttendanceSummary,
			0,
			len(studentSummary),
		)

	for _, summary := range studentSummary {

		if summary.AttendanceMarkedDays > 0 {

			summary.AttendancePercentage =
				(float64(summary.Present) /
					float64(
						summary.AttendanceMarkedDays,
					)) * 100
		}

		students = append(
			students,
			*summary,
		)
	}

	// Sort using roll number.
	sort.Slice(
		students,
		func(i int, j int) bool {

			if students[i].RollNumber ==
				students[j].RollNumber {

				return students[i].Name <
					students[j].Name
			}

			return students[i].RollNumber <
				students[j].RollNumber
		},
	)

	// ========================================
	// 10. CLASSROOM PERCENTAGE
	// ========================================

	totalMarkedAttendance :=
		present + absent

	percentage := 0.0

	if totalMarkedAttendance > 0 {

		percentage =
			(float64(present) /
				float64(
					totalMarkedAttendance,
				)) * 100
	}

	// ========================================
	// 11. RESPONSE
	// ========================================

	return &model.ClassroomMonthlyAttendanceSummary{
		ClassroomID:          classroomID,
		AcademicYear:         academicYear,
		Year:                 year,
		Month:                month,
		WorkingDays:          workingDays,
		AttendanceMarkedDays: len(markedDates),
		TotalStudents:        len(students),
		Present:              present,
		Absent:               absent,
		AttendancePercentage: percentage,
		Students:             students,
	}, nil
}
