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

	if req.RollNumber <= 0 {
		return nil, errors.New(
			"roll number must be greater than 0",
		)
	}

	if req.Age <= 0 {
		return nil, errors.New(
			"student age must be greater than 0",
		)
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
		return nil, errors.New(
			"username already exists",
		)
	}

	// ------------------------------------------------
	// 3. Convert classroom ID
	// ------------------------------------------------

	classroomObjectID, err := bson.ObjectIDFromHex(
		req.ClassroomID,
	)

	if err != nil {
		return nil, errors.New(
			"invalid classroom_id",
		)
	}

	// ------------------------------------------------
	// 4. Check classroom exists
	// ------------------------------------------------

	classroom, err := s.classroomRepo.GetByID(
		ctx,
		classroomObjectID,
	)

	if err != nil {
		return nil, errors.New(
			"classroom not found",
		)
	}

	// ------------------------------------------------
	// 5. Check duplicate roll number
	// ------------------------------------------------

	duplicateStudent, err :=
		s.studentRepo.ExistsByClassroomAndRollNumber(
			ctx,
			classroom.CollegeID,
			req.ClassroomID,
			req.RollNumber,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to check duplicate student: %w",
			err,
		)
	}

	if duplicateStudent {
		return nil, fmt.Errorf(
			"roll number %d already exists in %s - Section %s",
			req.RollNumber,
			classroom.Name,
			classroom.Section,
		)
	}

	// ------------------------------------------------
	// 6. Count students in classroom
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
	// 7. Check classroom capacity
	// ------------------------------------------------

	if count >= int64(classroom.Capacity) {
		return nil, fmt.Errorf(
			"%s - Section %s is full",
			classroom.Name,
			classroom.Section,
		)
	}

	// ------------------------------------------------
	// 8. Create student object
	// ------------------------------------------------

	student := studentModel.Student{
		CollegeID:   classroom.CollegeID,
		Name:        req.Name,
		Age:         req.Age,
		RollNumber:  req.RollNumber,
		Gender:      req.Gender,
		ClassroomID: req.ClassroomID,
		Address:     req.Address,
	}

	// ------------------------------------------------
	// 9. Save student
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
	// 10. Hash password
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
	// 11. Create student login account
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

	// ------------------------------------------
	// 1. Get existing student
	// ------------------------------------------

	existingStudent, err :=
		s.studentRepo.GetByID(
			ctx,
			id,
		)

	if err != nil {
		return errors.New(
			"student not found",
		)
	}

	// ------------------------------------------
	// 2. Validate classroom ID
	// ------------------------------------------

	classroomObjectID, err :=
		bson.ObjectIDFromHex(
			student.ClassroomID,
		)

	if err != nil {
		return errors.New(
			"invalid classroom_id",
		)
	}

	// ------------------------------------------
	// 3. Get classroom
	// ------------------------------------------

	classroom, err :=
		s.classroomRepo.GetByID(
			ctx,
			classroomObjectID,
		)

	if err != nil {
		return errors.New(
			"classroom not found",
		)
	}

	// ------------------------------------------
	// 4. Prevent duplicate roll number
	// ------------------------------------------

	exists, err :=
		s.studentRepo.
			ExistsByClassroomAndRollNumberExceptID(
				ctx,
				classroom.CollegeID,
				student.ClassroomID,
				student.RollNumber,
				id,
			)

	if err != nil {
		return fmt.Errorf(
			"failed to check duplicate student: %w",
			err,
		)
	}

	if exists {
		return fmt.Errorf(
			"roll number %d already exists in %s - Section %s",
			student.RollNumber,
			classroom.Name,
			classroom.Section,
		)
	}

	// ------------------------------------------
	// 5. Check classroom capacity
	// ------------------------------------------

	if existingStudent.ClassroomID !=
		student.ClassroomID {

		count, err :=
			s.studentRepo.CountByClassroom(
				ctx,
				student.ClassroomID,
			)

		if err != nil {
			return fmt.Errorf(
				"failed to check classroom capacity: %w",
				err,
			)
		}

		if count >= int64(
			classroom.Capacity,
		) {
			return fmt.Errorf(
				"%s - Section %s is full",
				classroom.Name,
				classroom.Section,
			)
		}
	}

	// ------------------------------------------
	// 6. Keep correct college ID
	// ------------------------------------------

	student.CollegeID =
		classroom.CollegeID

	// ------------------------------------------
	// 7. Update student
	// ------------------------------------------

	err = s.studentRepo.Update(
		ctx,
		id,
		student,
	)

	if err != nil {
		return err
	}

	return nil
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

// Get students by college
func (s *StudentService) GetByCollegeID(
	ctx context.Context,
	collegeID bson.ObjectID,
) ([]studentModel.Student, error) {

	return s.studentRepo.GetByCollegeID(
		ctx,
		collegeID,
	)
}
