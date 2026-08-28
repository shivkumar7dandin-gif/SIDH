package service

import (
	"context"
	"fmt"

	"github.com/shivkumar7dandin-gif/students-api/internal/classroom/repository"
	"github.com/shivkumar7dandin-gif/students-api/internal/student/model"
	studentRepository "github.com/shivkumar7dandin-gif/students-api/internal/student/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type StudentService struct {
	studentRepo   *studentRepository.StudentRepository
	classroomRepo *repository.ClassroomRepository
}

func NewStudentService(
	studentRepo *studentRepository.StudentRepository,
	classroomRepo *repository.ClassroomRepository,
) *StudentService {

	return &StudentService{
		studentRepo:   studentRepo,
		classroomRepo: classroomRepo,
	}
}

// Create Student
func (s *StudentService) Create(
	ctx context.Context,
	student model.Student,
) (*model.Student, error) {

	// ------------------------------------------------
	// 1. Check classroom exists
	// ------------------------------------------------

	classroom, err := s.classroomRepo.GetByName(
		ctx,
		student.ClassroomID,
	)

	if err != nil {
		return nil, fmt.Errorf("invalid classroom_id")
	}

	// ------------------------------------------------
	// 2. Count existing students in classroom
	// ------------------------------------------------

	count, err := s.studentRepo.CountByClassroom(
		ctx,
		student.ClassroomID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to count classroom students: %w",
			err,
		)
	}

	// ------------------------------------------------
	// 3. Check classroom capacity
	// ------------------------------------------------

	if count >= int64(classroom.Capacity) {
		return nil, fmt.Errorf(
			"classroom seats are full, not able to enroll in this classroom",
		)
	}

	// ------------------------------------------------
	// 4. Create student
	// ------------------------------------------------

	return s.studentRepo.Create(
		ctx,
		student,
	)
}

// Get all students
func (s *StudentService) GetAll(
	ctx context.Context,
) ([]model.Student, error) {

	return s.studentRepo.GetAll(ctx)
}

// Get student by ID
func (s *StudentService) GetByID(
	ctx context.Context,
	id bson.ObjectID,
) (*model.Student, error) {

	return s.studentRepo.GetByID(
		ctx,
		id,
	)
}

// Update student
func (s *StudentService) Update(
	ctx context.Context,
	id bson.ObjectID,
	student model.Student,
) error {

	return s.studentRepo.Update(
		ctx,
		id,
		student,
	)
}

// Delete student
func (s *StudentService) Delete(
	ctx context.Context,
	id bson.ObjectID,
) error {

	return s.studentRepo.Delete(
		ctx,
		id,
	)
}
