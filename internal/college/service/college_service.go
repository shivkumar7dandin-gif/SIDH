package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	attendanceRepository "github.com/shivkumar7dandin-gif/students-api/internal/attendance/repository"
	calendarService "github.com/shivkumar7dandin-gif/students-api/internal/calendar/service"
	classroomRepository "github.com/shivkumar7dandin-gif/students-api/internal/classroom/repository"
	collegeModel "github.com/shivkumar7dandin-gif/students-api/internal/college/model"
	collegeRepository "github.com/shivkumar7dandin-gif/students-api/internal/college/repository"
	studentModel "github.com/shivkumar7dandin-gif/students-api/internal/student/model"
	studentRepository "github.com/shivkumar7dandin-gif/students-api/internal/student/repository"
	teacherModel "github.com/shivkumar7dandin-gif/students-api/internal/teacher/model"
	teacherRepository "github.com/shivkumar7dandin-gif/students-api/internal/teacher/repository"
	timetableRepository "github.com/shivkumar7dandin-gif/students-api/internal/timetable/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type CollegeService struct {
	collegeRepo     *collegeRepository.CollegeRepository
	teacherRepo     *teacherRepository.TeacherRepository
	classroomRepo   *classroomRepository.ClassroomRepository
	studentRepo     *studentRepository.StudentRepository
	attendanceRepo  *attendanceRepository.AttendanceRepository
	timetableRepo   *timetableRepository.TimetableRepository
	calendarService *calendarService.CalendarService
}

func NewCollegeService(
	collegeRepo *collegeRepository.CollegeRepository,
	teacherRepo *teacherRepository.TeacherRepository,
	classroomRepo *classroomRepository.ClassroomRepository,
	studentRepo *studentRepository.StudentRepository,
	attendanceRepo *attendanceRepository.AttendanceRepository,
	timetableRepo *timetableRepository.TimetableRepository,
	calendarService *calendarService.CalendarService,
) *CollegeService {

	return &CollegeService{
		collegeRepo:     collegeRepo,
		teacherRepo:     teacherRepo,
		classroomRepo:   classroomRepo,
		studentRepo:     studentRepo,
		attendanceRepo:  attendanceRepo,
		timetableRepo:   timetableRepo,
		calendarService: calendarService,
	}
}

// ========================================
// CREATE COLLEGE
// ========================================

func (s *CollegeService) Create(
	ctx context.Context,
	req collegeModel.CreateCollegeRequest,
) (*collegeModel.College, error) {

	if strings.TrimSpace(req.Name) == "" {
		return nil, fmt.Errorf("college name is required")
	}

	if strings.TrimSpace(req.Code) == "" {
		return nil, fmt.Errorf("college code is required")
	}

	if strings.TrimSpace(req.Username) == "" {
		return nil, fmt.Errorf("username is required")
	}

	if strings.TrimSpace(req.Password) == "" {
		return nil, fmt.Errorf("password is required")
	}

	if len(req.Password) < 8 {
		return nil, fmt.Errorf(
			"password must be at least 8 characters",
		)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to hash password: %w",
			err,
		)
	}

	college := collegeModel.College{
		Name:         req.Name,
		Code:         req.Code,
		Address:      req.Address,
		City:         req.City,
		State:        req.State,
		Pincode:      req.Pincode,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         "college_admin",
	}

	return s.collegeRepo.Create(ctx, college)
}

// ========================================
// GET ALL COLLEGES
// ========================================

func (s *CollegeService) GetAll(
	ctx context.Context,
) ([]collegeModel.College, error) {

	return s.collegeRepo.GetAll(ctx)
}

// ========================================
// GET COLLEGE BY ID
// ========================================

func (s *CollegeService) GetByID(
	ctx context.Context,
	id bson.ObjectID,
) (*collegeModel.College, error) {

	return s.collegeRepo.GetByID(ctx, id)
}

// ========================================
// GET COLLEGE DASHBOARD
// ========================================

func (s *CollegeService) GetDashboard(
	ctx context.Context,
	collegeID bson.ObjectID,
) (*collegeModel.CollegeDashboard, error) {

	// ----------------------------------------
	// 1. Get college
	// ----------------------------------------

	college, err := s.collegeRepo.GetByID(
		ctx,
		collegeID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get college: %w",
			err,
		)
	}

	// ----------------------------------------
	// 2. Get teachers
	// ----------------------------------------

	teachers, err := s.teacherRepo.GetByCollegeID(
		ctx,
		collegeID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get teachers: %w",
			err,
		)
	}

	// ----------------------------------------
	// 3. Calculate teacher status
	// ----------------------------------------

	activeTeachers := make(
		[]teacherModel.Teacher,
		0,
	)

	leftTeachers := make(
		[]teacherModel.Teacher,
		0,
	)

	for _, teacher := range teachers {

		status := strings.ToLower(
			strings.TrimSpace(teacher.Status),
		)

		switch status {

		case "left":
			leftTeachers = append(
				leftTeachers,
				teacher,
			)

		default:
			activeTeachers = append(
				activeTeachers,
				teacher,
			)
		}
	}

	activeTeacherCount := len(activeTeachers)
	leftTeacherCount := len(leftTeachers)

	// ----------------------------------------
	// 4. Get classrooms
	// ----------------------------------------

	classrooms, err := s.classroomRepo.GetByCollegeID(
		ctx,
		collegeID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get classrooms: %w",
			err,
		)
	}

	// ----------------------------------------
	// 5. Get students
	// ----------------------------------------

	students, err := s.studentRepo.GetByCollegeID(
		ctx,
		collegeID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get students: %w",
			err,
		)
	}

	// ----------------------------------------
	// 6. Get today's date - India timezone
	// ----------------------------------------

	indiaLocation, err := time.LoadLocation(
		"Asia/Kolkata",
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to load server timezone: %w",
			err,
		)
	}

	nowIndia := time.Now().In(indiaLocation)

	today := nowIndia.Format("2006-01-02")

	// ----------------------------------------
	// 7. Get today's attendance
	// ----------------------------------------

	todayRecords, err :=
		s.attendanceRepo.GetByCollegeAndDate(
			ctx,
			collegeID,
			today,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get today's attendance: %w",
			err,
		)
	}

	// ----------------------------------------
	// 8. Create today's attendance map
	// ----------------------------------------

	attendanceStatus := make(
		map[bson.ObjectID]string,
	)

	for _, attendance := range todayRecords {

		switch attendance.Attendance {

		case "Present", "Absent":
			attendanceStatus[attendance.StudentID] =
				attendance.Attendance
		}
	}

	// ----------------------------------------
	// 9. Separate students by attendance
	// ----------------------------------------

	presentStudents := make(
		[]studentModel.Student,
		0,
	)

	absentStudents := make(
		[]studentModel.Student,
		0,
	)

	notMarkedStudents := make(
		[]studentModel.Student,
		0,
	)

	for _, student := range students {

		status, exists :=
			attendanceStatus[student.ID]

		if !exists {

			notMarkedStudents = append(
				notMarkedStudents,
				student,
			)

			continue
		}

		switch status {

		case "Present":
			presentStudents = append(
				presentStudents,
				student,
			)

		case "Absent":
			absentStudents = append(
				absentStudents,
				student,
			)
		}
	}

	// ----------------------------------------
	// 10. Calculate attendance totals
	// ----------------------------------------

	present := len(presentStudents)
	absent := len(absentStudents)
	notMarked := len(notMarkedStudents)

	marked := present + absent

	// ----------------------------------------
	// 11. Calculate attendance percentage
	// ----------------------------------------

	percentage := 0.0

	if marked > 0 {
		percentage =
			(float64(present) /
				float64(marked)) * 100
	}

	// ----------------------------------------
	// 12. Build classroom-wise attendance
	// ----------------------------------------

	classroomAttendance := make(
		[]collegeModel.ClassroomAttendance,
		0,
		len(classrooms),
	)

	for _, classroom := range classrooms {

		classroomID := classroom.ID.Hex()

		totalStudents := 0
		classPresent := 0
		classAbsent := 0
		classNotMarked := 0

		for _, student := range students {

			if student.ClassroomID != classroomID {
				continue
			}

			totalStudents++

			status, exists :=
				attendanceStatus[student.ID]

			if !exists {
				classNotMarked++
				continue
			}

			switch status {

			case "Present":
				classPresent++

			case "Absent":
				classAbsent++

			default:
				classNotMarked++
			}
		}

		classMarked :=
			classPresent + classAbsent

		classPercentage := 0.0

		if classMarked > 0 {
			classPercentage =
				(float64(classPresent) /
					float64(classMarked)) * 100
		}

		classStatus := "Not Marked"

		if totalStudents == 0 {
			classStatus = "No Students"

		} else if classMarked == 0 {
			classStatus = "Not Marked"

		} else if classMarked < totalStudents {
			classStatus = "Partial"

		} else {
			classStatus = "Completed"
		}

		classroomAttendance = append(
			classroomAttendance,
			collegeModel.ClassroomAttendance{
				ClassroomID:   classroomID,
				ClassroomName: classroom.Name,
				Section:       classroom.Section,

				TotalStudents: totalStudents,

				Present:   classPresent,
				Absent:    classAbsent,
				Marked:    classMarked,
				NotMarked: classNotMarked,

				Percentage: classPercentage,
				Status:     classStatus,
			},
		)
	}

	// ----------------------------------------
	// 13. Calculate today's teacher classes
	// ----------------------------------------

	teacherTodayClasses := make(
		[]collegeModel.TeacherTodayClass,
		0,
	)

	// Find the active academic year that contains today.
	academicCalendars, err :=
		s.calendarService.GetByCollegeID(
			ctx,
			collegeID,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get academic calendars: %w",
			err,
		)
	}

	activeAcademicYear := ""

	todayDate := time.Date(
		nowIndia.Year(),
		nowIndia.Month(),
		nowIndia.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)

	for _, calendar := range academicCalendars {

		status := strings.ToLower(
			strings.TrimSpace(calendar.Status),
		)

		if status != "active" {
			continue
		}

		startDate := time.Date(
			calendar.StartDate.Year(),
			calendar.StartDate.Month(),
			calendar.StartDate.Day(),
			0,
			0,
			0,
			0,
			time.UTC,
		)

		endDate := time.Date(
			calendar.EndDate.Year(),
			calendar.EndDate.Month(),
			calendar.EndDate.Day(),
			0,
			0,
			0,
			0,
			time.UTC,
		)

		if !todayDate.Before(startDate) &&
			!todayDate.After(endDate) {

			activeAcademicYear =
				calendar.AcademicYear

			break
		}
	}

	// ----------------------------------------
	// 13. Count today's classes per teacher
	// ----------------------------------------

	teacherClassCount := make(
		map[bson.ObjectID]int,
	)

	if activeAcademicYear != "" {

		dayStatus, err :=
			s.calendarService.GetDayStatus(
				ctx,
				collegeID,
				activeAcademicYear,
				nowIndia,
			)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to validate today's academic calendar: %w",
				err,
			)
		}

		// Only count classes when today
		// is a working day.
		if dayStatus.IsWorkingDay {

			todayDay := strings.ToLower(
				nowIndia.Weekday().String(),
			)

			todayTimetables, err :=
				s.timetableRepo.GetByCollegeAndDay(
					ctx,
					collegeID,
					todayDay,
					activeAcademicYear,
				)

			if err != nil {
				return nil, fmt.Errorf(
					"failed to get today's timetable: %w",
					err,
				)
			}

			for _, timetable := range todayTimetables {

				teacherClassCount[timetable.TeacherID]++
			}
		}
	}

	// ----------------------------------------
	// 14. Build teacher class summaries
	// ----------------------------------------

	for _, teacher := range activeTeachers {

		teacherTodayClasses = append(
			teacherTodayClasses,
			collegeModel.TeacherTodayClass{
				TeacherID:    teacher.ID,
				TeacherName:  teacher.Name,
				Subject:      teacher.Subject,
				TotalClasses: teacherClassCount[teacher.ID],
			},
		)
	}

	// ----------------------------------------
	// 15. Build dashboard response
	// ----------------------------------------

	return &collegeModel.CollegeDashboard{

		College:    college,
		Teachers:   teachers,
		Classrooms: classrooms,
		Students:   students,

		TotalTeachers:   len(teachers),
		ActiveTeachers:  activeTeacherCount,
		LeftTeachers:    leftTeacherCount,
		TotalClassrooms: len(classrooms),
		TotalStudents:   len(students),

		// ------------------------------------
		// Teacher status
		// ------------------------------------

		TeacherStatus: collegeModel.TeacherStatusSummary{

			Total:  len(teachers),
			Active: activeTeacherCount,
			Left:   leftTeacherCount,

			ActiveTeachers: activeTeachers,
			LeftTeachers:   leftTeachers,
		},

		// ------------------------------------
		// Student attendance
		// ------------------------------------

		TodayAttendance: collegeModel.TodayAttendance{

			TotalStudents: len(students),

			Present:   present,
			Absent:    absent,
			Marked:    marked,
			NotMarked: notMarked,

			Percentage: percentage,

			PresentStudents: presentStudents,

			AbsentStudents: absentStudents,

			NotMarkedStudents: notMarkedStudents,
		},
		ClassroomAttendance: classroomAttendance,
		// ------------------------------------
		// Today's teacher classes
		// ------------------------------------

		TeacherTodayClasses: teacherTodayClasses,
	}, nil
}
