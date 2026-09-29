package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	calendarService "github.com/shivkumar7dandin-gif/students-api/internal/calendar/service"
	classroomRepository "github.com/shivkumar7dandin-gif/students-api/internal/classroom/repository"
	teacherRepository "github.com/shivkumar7dandin-gif/students-api/internal/teacher/repository"
	timetableModel "github.com/shivkumar7dandin-gif/students-api/internal/timetable/model"
	timetableRepository "github.com/shivkumar7dandin-gif/students-api/internal/timetable/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TimetableService struct {
	timetableRepo   *timetableRepository.TimetableRepository
	teacherRepo     *teacherRepository.TeacherRepository
	classroomRepo   *classroomRepository.ClassroomRepository
	calendarService *calendarService.CalendarService
}

func NewTimetableService(
	timetableRepo *timetableRepository.TimetableRepository,
	teacherRepo *teacherRepository.TeacherRepository,
	classroomRepo *classroomRepository.ClassroomRepository,
	calendarService *calendarService.CalendarService,
) *TimetableService {

	return &TimetableService{
		timetableRepo:   timetableRepo,
		teacherRepo:     teacherRepo,
		classroomRepo:   classroomRepo,
		calendarService: calendarService,
	}
}

// ========================================
// CREATE TIMETABLE
// ========================================

func (s *TimetableService) Create(
	ctx context.Context,
	collegeID bson.ObjectID,
	req timetableModel.CreateTimetableRequest,
) (*timetableModel.Timetable, error) {

	// ----------------------------------------
	// 1. Validate teacher ID
	// ----------------------------------------

	if strings.TrimSpace(req.TeacherID) == "" {
		return nil, fmt.Errorf(
			"teacher id is required",
		)
	}

	teacherID, err := bson.ObjectIDFromHex(
		strings.TrimSpace(req.TeacherID),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid teacher id",
		)
	}

	// ----------------------------------------
	// 2. Validate classroom ID
	// ----------------------------------------

	if strings.TrimSpace(req.ClassroomID) == "" {
		return nil, fmt.Errorf(
			"classroom id is required",
		)
	}

	classroomID, err := bson.ObjectIDFromHex(
		strings.TrimSpace(req.ClassroomID),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid classroom id",
		)
	}

	// ----------------------------------------
	// 3. Validate subject
	// ----------------------------------------

	subject := strings.TrimSpace(req.Subject)

	if subject == "" {
		return nil, fmt.Errorf(
			"subject is required",
		)
	}

	// ----------------------------------------
	// 4. Validate academic year
	// ----------------------------------------

	academicYear := strings.TrimSpace(
		req.AcademicYear,
	)

	if academicYear == "" {
		return nil, fmt.Errorf(
			"academic year is required",
		)
	}

	// ----------------------------------------
	// Validate active academic year
	// ----------------------------------------

	if s.calendarService == nil {
		return nil, fmt.Errorf(
			"calendar service is not configured",
		)
	}

	if err := s.calendarService.ValidateActiveAcademicYear(
		ctx,
		collegeID,
		academicYear,
	); err != nil {
		return nil, err
	}

	// ----------------------------------------
	// 5. Validate day
	// ----------------------------------------

	dayOfWeek := strings.ToLower(
		strings.TrimSpace(req.DayOfWeek),
	)

	if !isValidDay(dayOfWeek) {
		return nil, fmt.Errorf(
			"invalid day of week",
		)
	}

	// ----------------------------------------
	// 6. Validate start/end time
	// ----------------------------------------

	startTime := strings.TrimSpace(
		req.StartTime,
	)

	endTime := strings.TrimSpace(
		req.EndTime,
	)

	startParsed, err := time.Parse(
		"15:04",
		startTime,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid start time, use HH:MM format",
		)
	}

	endParsed, err := time.Parse(
		"15:04",
		endTime,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid end time, use HH:MM format",
		)
	}

	if !startParsed.Before(endParsed) {
		return nil, fmt.Errorf(
			"start time must be before end time",
		)
	}

	// ----------------------------------------
	// 7. Get teacher
	// ----------------------------------------

	teacher, err := s.teacherRepo.GetByID(
		ctx,
		teacherID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"teacher not found",
		)
	}

	// ----------------------------------------
	// 8. Verify teacher belongs to college
	// ----------------------------------------

	if teacher.CollegeID != collegeID {
		return nil, fmt.Errorf(
			"teacher does not belong to this college",
		)
	}

	// ----------------------------------------
	// 9. Teacher must be active
	// ----------------------------------------

	teacherStatus := strings.ToLower(
		strings.TrimSpace(teacher.Status),
	)

	if teacherStatus == "left" {
		return nil, fmt.Errorf(
			"cannot assign timetable to a teacher who has left",
		)
	}

	// ----------------------------------------
	// 10. Get classroom
	// ----------------------------------------

	classroom, err := s.classroomRepo.GetByID(
		ctx,
		classroomID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"classroom not found",
		)
	}

	// ----------------------------------------
	// 11. Verify classroom belongs to college
	// ----------------------------------------

	if classroom.CollegeID != collegeID {
		return nil, fmt.Errorf(
			"classroom does not belong to this college",
		)
	}

	// ----------------------------------------
	// 12. Check teacher timetable conflicts
	// ----------------------------------------

	teacherTimetables, err :=
		s.timetableRepo.GetByTeacherAndDay(
			ctx,
			collegeID,
			teacherID,
			dayOfWeek,
			academicYear,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to check teacher timetable: %w",
			err,
		)
	}

	for _, existing := range teacherTimetables {

		if timesOverlap(
			startTime,
			endTime,
			existing.StartTime,
			existing.EndTime,
		) {
			return nil, fmt.Errorf(
				"teacher already has a class between %s and %s",
				existing.StartTime,
				existing.EndTime,
			)
		}
	}

	// ----------------------------------------
	// 13. Check classroom timetable conflicts
	// ----------------------------------------

	classroomTimetables, err :=
		s.timetableRepo.GetByClassroomAndDay(
			ctx,
			collegeID,
			classroomID,
			dayOfWeek,
			academicYear,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to check classroom timetable: %w",
			err,
		)
	}

	for _, existing := range classroomTimetables {

		if timesOverlap(
			startTime,
			endTime,
			existing.StartTime,
			existing.EndTime,
		) {
			return nil, fmt.Errorf(
				"classroom already has a class between %s and %s",
				existing.StartTime,
				existing.EndTime,
			)
		}
	}

	// ----------------------------------------
	// 14. Build timetable
	// ----------------------------------------

	timetable := timetableModel.Timetable{
		CollegeID:    collegeID,
		TeacherID:    teacherID,
		ClassroomID:  classroomID,
		Subject:      subject,
		DayOfWeek:    dayOfWeek,
		StartTime:    startTime,
		EndTime:      endTime,
		AcademicYear: academicYear,
	}

	// ----------------------------------------
	// 15. Save timetable
	// ----------------------------------------

	createdTimetable, err :=
		s.timetableRepo.Create(
			ctx,
			timetable,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to create timetable: %w",
			err,
		)
	}

	return createdTimetable, nil
}

// ========================================
// GET ALL TIMETABLES
// ========================================

func (s *TimetableService) GetByAcademicYear(
	ctx context.Context,
	collegeID bson.ObjectID,
	academicYear string,
) ([]timetableModel.Timetable, error) {

	academicYear = strings.TrimSpace(
		academicYear,
	)

	if academicYear == "" {
		return nil, fmt.Errorf(
			"academic year is required",
		)
	}

	return s.timetableRepo.GetByCollegeAndAcademicYear(
		ctx,
		collegeID,
		academicYear,
	)
}

// ========================================
// GET TIMETABLE BY DAY
// ========================================

func (s *TimetableService) GetByDay(
	ctx context.Context,
	collegeID bson.ObjectID,
	dayOfWeek string,
	academicYear string,
) ([]timetableModel.Timetable, error) {

	dayOfWeek = strings.ToLower(
		strings.TrimSpace(dayOfWeek),
	)

	academicYear = strings.TrimSpace(
		academicYear,
	)

	if !isValidDay(dayOfWeek) {
		return nil, fmt.Errorf(
			"invalid day of week",
		)
	}

	if academicYear == "" {
		return nil, fmt.Errorf(
			"academic year is required",
		)
	}

	return s.timetableRepo.GetByCollegeAndDay(
		ctx,
		collegeID,
		dayOfWeek,
		academicYear,
	)
}

// ========================================
// GET TEACHER TIMETABLE BY DAY
// ========================================

func (s *TimetableService) GetByTeacherAndDay(
	ctx context.Context,
	collegeID bson.ObjectID,
	teacherID bson.ObjectID,
	dayOfWeek string,
	academicYear string,
) ([]timetableModel.Timetable, error) {

	dayOfWeek = strings.ToLower(
		strings.TrimSpace(dayOfWeek),
	)

	academicYear = strings.TrimSpace(
		academicYear,
	)

	if !isValidDay(dayOfWeek) {
		return nil, fmt.Errorf(
			"invalid day of week",
		)
	}

	if academicYear == "" {
		return nil, fmt.Errorf(
			"academic year is required",
		)
	}

	return s.timetableRepo.GetByTeacherAndDay(
		ctx,
		collegeID,
		teacherID,
		dayOfWeek,
		academicYear,
	)
}

// ========================================
// VALID DAY
// ========================================

func isValidDay(
	day string,
) bool {

	switch day {

	case "monday",
		"tuesday",
		"wednesday",
		"thursday",
		"friday",
		"saturday",
		"sunday":

		return true

	default:
		return false
	}
}

// ========================================
// TIME OVERLAP CHECK
// ========================================

// Example:
//
// Existing:
// 09:00 - 10:00
//
// New:
// 09:30 - 10:30
//
// Result:
// conflict
//
// Formula:
//
// newStart < existingEnd
// AND
// newEnd > existingStart
func timesOverlap(
	newStart string,
	newEnd string,
	existingStart string,
	existingEnd string,
) bool {

	newStartTime, err := time.Parse(
		"15:04",
		newStart,
	)

	if err != nil {
		return false
	}

	newEndTime, err := time.Parse(
		"15:04",
		newEnd,
	)

	if err != nil {
		return false
	}

	existingStartTime, err := time.Parse(
		"15:04",
		existingStart,
	)

	if err != nil {
		return false
	}

	existingEndTime, err := time.Parse(
		"15:04",
		existingEnd,
	)

	if err != nil {
		return false
	}

	return newStartTime.Before(existingEndTime) &&
		newEndTime.After(existingStartTime)
}

// ========================================
// UPDATE TIMETABLE
// ========================================

func (s *TimetableService) Update(
	ctx context.Context,
	collegeID bson.ObjectID,
	timetableID bson.ObjectID,
	req timetableModel.UpdateTimetableRequest,
) (*timetableModel.Timetable, error) {

	// ----------------------------------------
	// 1. Check existing timetable
	// ----------------------------------------

	_, err := s.timetableRepo.GetByID(
		ctx,
		collegeID,
		timetableID,
	)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------
	// 2. Validate teacher ID
	// ----------------------------------------

	if strings.TrimSpace(req.TeacherID) == "" {
		return nil, fmt.Errorf(
			"teacher id is required",
		)
	}

	teacherID, err := bson.ObjectIDFromHex(
		strings.TrimSpace(req.TeacherID),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid teacher id",
		)
	}

	// ----------------------------------------
	// 3. Validate classroom ID
	// ----------------------------------------

	if strings.TrimSpace(req.ClassroomID) == "" {
		return nil, fmt.Errorf(
			"classroom id is required",
		)
	}

	classroomID, err := bson.ObjectIDFromHex(
		strings.TrimSpace(req.ClassroomID),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid classroom id",
		)
	}

	// ----------------------------------------
	// 4. Validate subject
	// ----------------------------------------

	subject := strings.TrimSpace(
		req.Subject,
	)

	if subject == "" {
		return nil, fmt.Errorf(
			"subject is required",
		)
	}

	// ----------------------------------------
	// 5. Validate academic year
	// ----------------------------------------

	academicYear := strings.TrimSpace(
		req.AcademicYear,
	)

	if academicYear == "" {
		return nil, fmt.Errorf(
			"academic year is required",
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
		academicYear,
	); err != nil {
		return nil, err
	}

	// ----------------------------------------
	// 6. Validate day
	// ----------------------------------------

	dayOfWeek := strings.ToLower(
		strings.TrimSpace(req.DayOfWeek),
	)

	if !isValidDay(dayOfWeek) {
		return nil, fmt.Errorf(
			"invalid day of week",
		)
	}

	// ----------------------------------------
	// 7. Validate time
	// ----------------------------------------

	startTime := strings.TrimSpace(
		req.StartTime,
	)

	endTime := strings.TrimSpace(
		req.EndTime,
	)

	startParsed, err := time.Parse(
		"15:04",
		startTime,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid start time, use HH:MM format",
		)
	}

	endParsed, err := time.Parse(
		"15:04",
		endTime,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid end time, use HH:MM format",
		)
	}

	if !startParsed.Before(endParsed) {
		return nil, fmt.Errorf(
			"start time must be before end time",
		)
	}

	// ----------------------------------------
	// 8. Validate teacher
	// ----------------------------------------

	teacher, err := s.teacherRepo.GetByID(
		ctx,
		teacherID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"teacher not found",
		)
	}

	if teacher.CollegeID != collegeID {
		return nil, fmt.Errorf(
			"teacher does not belong to this college",
		)
	}

	teacherStatus := strings.ToLower(
		strings.TrimSpace(teacher.Status),
	)

	if teacherStatus == "left" {
		return nil, fmt.Errorf(
			"cannot assign timetable to a teacher who has left",
		)
	}

	// ----------------------------------------
	// 9. Validate classroom
	// ----------------------------------------

	classroom, err := s.classroomRepo.GetByID(
		ctx,
		classroomID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"classroom not found",
		)
	}

	if classroom.CollegeID != collegeID {
		return nil, fmt.Errorf(
			"classroom does not belong to this college",
		)
	}

	// ----------------------------------------
	// 10. Teacher conflict check
	// ----------------------------------------

	teacherTimetables, err :=
		s.timetableRepo.GetByTeacherAndDay(
			ctx,
			collegeID,
			teacherID,
			dayOfWeek,
			academicYear,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to check teacher timetable: %w",
			err,
		)
	}

	for _, existing := range teacherTimetables {

		// Ignore current timetable being edited
		if existing.ID == timetableID {
			continue
		}

		if timesOverlap(
			startTime,
			endTime,
			existing.StartTime,
			existing.EndTime,
		) {
			return nil, fmt.Errorf(
				"teacher already has a class between %s and %s",
				existing.StartTime,
				existing.EndTime,
			)
		}
	}

	// ----------------------------------------
	// 11. Classroom conflict check
	// ----------------------------------------

	classroomTimetables, err :=
		s.timetableRepo.GetByClassroomAndDay(
			ctx,
			collegeID,
			classroomID,
			dayOfWeek,
			academicYear,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to check classroom timetable: %w",
			err,
		)
	}

	for _, existing := range classroomTimetables {

		// Ignore current timetable being edited
		if existing.ID == timetableID {
			continue
		}

		if timesOverlap(
			startTime,
			endTime,
			existing.StartTime,
			existing.EndTime,
		) {
			return nil, fmt.Errorf(
				"classroom already has a class between %s and %s",
				existing.StartTime,
				existing.EndTime,
			)
		}
	}

	// ----------------------------------------
	// 12. Build updated timetable
	// ----------------------------------------

	timetable := timetableModel.Timetable{
		ID:           timetableID,
		CollegeID:    collegeID,
		TeacherID:    teacherID,
		ClassroomID:  classroomID,
		Subject:      subject,
		DayOfWeek:    dayOfWeek,
		StartTime:    startTime,
		EndTime:      endTime,
		AcademicYear: academicYear,
	}

	// ----------------------------------------
	// 13. Update
	// ----------------------------------------

	updatedTimetable, err :=
		s.timetableRepo.Update(
			ctx,
			collegeID,
			timetableID,
			timetable,
		)

	if err != nil {
		return nil, err
	}

	return updatedTimetable, nil
}

// ========================================
// DELETE TIMETABLE
// ========================================

func (s *TimetableService) Delete(
	ctx context.Context,
	collegeID bson.ObjectID,
	timetableID bson.ObjectID,
) error {

	// Make sure timetable belongs to this college
	_, err := s.timetableRepo.GetByID(
		ctx,
		collegeID,
		timetableID,
	)

	if err != nil {
		return err
	}

	return s.timetableRepo.Delete(
		ctx,
		collegeID,
		timetableID,
	)
}
