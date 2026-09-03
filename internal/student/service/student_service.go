package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	classroomRepository "github.com/shivkumar7dandin-gif/students-api/internal/classroom/repository"
	studentModel "github.com/shivkumar7dandin-gif/students-api/internal/student/model"
	studentRepository "github.com/shivkumar7dandin-gif/students-api/internal/student/repository"
	userModel "github.com/shivkumar7dandin-gif/students-api/internal/user/model"
	userRepository "github.com/shivkumar7dandin-gif/students-api/internal/user/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type StudentService struct {
	studentRepo   *studentRepository.StudentRepository
	classroomRepo *classroomRepository.ClassroomRepository
	userRepo      *userRepository.UserRepository
}

func NewStudentService(
	studentRepo *studentRepository.StudentRepository,
	classroomRepo *classroomRepository.ClassroomRepository,
	userRepo *userRepository.UserRepository,
) *StudentService {

	return &StudentService{
		studentRepo:   studentRepo,
		classroomRepo: classroomRepo,
		userRepo:      userRepo,
	}
}

// Create student
func (s *StudentService) Create(
	ctx context.Context,
	req studentModel.CreateStudentRequest,
) (*studentModel.Student, error) {

	req.Name = strings.TrimSpace(req.Name)
	req.Username = strings.TrimSpace(req.Username)
	req.ClassroomID = strings.TrimSpace(req.ClassroomID)

	// ------------------------------------------------
	// 1. Basic validation
	// ------------------------------------------------

	if req.Name == "" {
		return nil, errors.New("student name is required")
	}

	if req.Username == "" {
		return nil, errors.New("username is required")
	}

	if len(req.Password) < 8 {
		return nil, errors.New(
			"password must contain at least 8 characters",
		)
	}

	if req.ClassroomID == "" {
		return nil, errors.New("classroom_id is required")
	}

	// ------------------------------------------------
	// 2. Check username already exists
	// ------------------------------------------------

	exists, err := s.userRepo.UsernameExists(
		ctx,
		req.Username,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to check username: %w",
			err,
		)
	}

	if exists {
		return nil, errors.New("username already exists")
	}

	// ------------------------------------------------
	// 3. Check classroom exists
	// ------------------------------------------------

	classroom, err := s.classroomRepo.GetByName(
		ctx,
		req.ClassroomID,
	)
	if err != nil {
		return nil, errors.New("invalid classroom_id")
	}

	// ------------------------------------------------
	// 4. Count students in classroom
	// ------------------------------------------------

	count, err := s.studentRepo.CountByClassroom(
		ctx,
		req.ClassroomID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to count classroom students: %w",
			err,
		)
	}

	// ------------------------------------------------
	// 5. Check classroom capacity
	// ------------------------------------------------

	if count >= int64(classroom.Capacity) {
		return nil, errors.New(
			"classroom seats are full, not able to enroll in this classroom",
		)
	}

	// ------------------------------------------------
	// 6. Create student object
	// ------------------------------------------------

	student := studentModel.Student{
		Name:        req.Name,
		Age:         req.Age,
		RollNumber:  req.RollNumber,
		Gender:      req.Gender,
		ClassroomID: req.ClassroomID,
		Address:     req.Address,
	}

	// ------------------------------------------------
	// 7. Save student
	// ------------------------------------------------

	createdStudent, err := s.studentRepo.Create(
		ctx,
		student,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create student: %w",
			err,
		)
	}

	// ------------------------------------------------
	// 8. Hash student password
	// ------------------------------------------------

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to hash password: %w",
			err,
		)
	}

	// ------------------------------------------------
	// 9. Create student user account
	// ------------------------------------------------

	user := userModel.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		Role:         "student",
		ReferenceID:  createdStudent.ID,
	}

	if err := s.userRepo.Create(
		ctx,
		user,
	); err != nil {

		return nil, fmt.Errorf(
			"failed to create student login account: %w",
			err,
		)
	}

	return createdStudent, nil
}

// Get all students
func (s *StudentService) GetAll(
	ctx context.Context,
) ([]studentModel.Student, error) {

	return s.studentRepo.GetAll(ctx)
}

// Get student by ID
func (s *StudentService) GetByID(
	ctx context.Context,
	id bson.ObjectID,
) (*studentModel.Student, error) {

	return s.studentRepo.GetByID(
		ctx,
		id,
	)
}

// Update student
func (s *StudentService) Update(
	ctx context.Context,
	id bson.ObjectID,
	student studentModel.Student,
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
