package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/shivkumar7dandin-gif/students-api/internal/attendance/model"
	"github.com/shivkumar7dandin-gif/students-api/internal/attendance/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AttendanceService struct {
	repo *repository.AttendanceRepository
}

func NewAttendanceService(
	repo *repository.AttendanceRepository,
) *AttendanceService {

	return &AttendanceService{
		repo: repo,
	}
}

// ========================================
// CREATE ATTENDANCE
// ========================================

func (s *AttendanceService) Create(
	ctx context.Context,
	attendance model.Attendance,
) (*model.Attendance, error) {

	// ----------------------------------------
	// 1. Validate attendance status
	// ----------------------------------------

	status := strings.ToLower(attendance.Attendance)

	if status != "present" && status != "absent" {
		return nil, fmt.Errorf(
			"attendance must be Present or Absent",
		)
	}

	// Convert status into standard format
	if status == "present" {
		attendance.Attendance = "Present"
	} else {
		attendance.Attendance = "Absent"
	}

	// ----------------------------------------
	// 2. Validate date
	// ----------------------------------------

	if attendance.Date == "" {
		return nil, fmt.Errorf("date is required")
	}

	// ----------------------------------------
	// 3. Validate classroom
	// ----------------------------------------

	if attendance.ClassroomID == "" {
		return nil, fmt.Errorf("classroom_id is required")
	}

	// ----------------------------------------
	// 4. Check duplicate attendance
	// ----------------------------------------

	existingAttendance, err := s.repo.GetByStudentAndDate(
		ctx,
		attendance.StudentID,
		attendance.Date,
	)

	if err == nil && existingAttendance != nil {
		return nil, fmt.Errorf(
			"attendance already exists for this student on %s",
			attendance.Date,
		)
	}

	// If MongoDB returns another error,
	// return that error.
	//
	// mongo.ErrNoDocuments is OK because it means
	// attendance does not exist yet.
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf(
			"failed to check existing attendance: %w",
			err,
		)
	}

	// ----------------------------------------
	// 5. Create attendance
	// ----------------------------------------

	return s.repo.Create(
		ctx,
		attendance,
	)
}

// ========================================
// GET ATTENDANCE BY STUDENT
// ========================================

func (s *AttendanceService) GetByStudent(
	ctx context.Context,
	studentID bson.ObjectID,
) ([]model.Attendance, error) {

	return s.repo.GetByStudent(
		ctx,
		studentID,
	)
}

// ========================================
// GET ALL ATTENDANCE
// ========================================

func (s *AttendanceService) GetAll(
	ctx context.Context,
) ([]model.Attendance, error) {

	return s.repo.GetAll(ctx)
}
