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

func (s *StudentService) Create(
	ctx context.Context,
	req studentModel.CreateStudentRequest,
	collegeID bson.ObjectID,
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
		return nil, errors.New(
			"classroom_id is required",
		)
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

	// Guardian validation
	if err := validateGuardians(req.Guardians); err != nil {
		return nil, err
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
	// 4. Check classroom belongs to logged-in college
	// ------------------------------------------------

	classroom, err :=
		s.classroomRepo.GetByIDAndCollegeID(
			ctx,
			classroomObjectID,
			collegeID,
		)

	if err != nil {
		return nil, errors.New(
			"classroom not found or does not belong to your school",
		)
	}

	// ------------------------------------------------
	// 5. Check duplicate roll number
	// ------------------------------------------------

	duplicateStudent, err :=
		s.studentRepo.ExistsByClassroomAndRollNumber(
			ctx,
			collegeID,
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
		collegeID,
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
		CollegeID:   collegeID,
		Name:        req.Name,
		Age:         req.Age,
		RollNumber:  req.RollNumber,
		Gender:      req.Gender,
		ClassroomID: req.ClassroomID,
		Address:     req.Address,
		Guardians:   req.Guardians,
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
		CollegeID:    collegeID,
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
				classroom.CollegeID,
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

func (s *StudentService) GetByIDAndCollegeID(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
) (*studentModel.Student, error) {

	return s.studentRepo.GetByIDAndCollegeID(
		ctx,
		id,
		collegeID,
	)
}

func (s *StudentService) UpdateByCollegeID(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
	student studentModel.Student,
) error {

	existingStudent, err := s.studentRepo.GetByIDAndCollegeID(
		ctx,
		id,
		collegeID,
	)
	if err != nil {
		return errors.New("student not found")
	}

	classroomObjectID, err := bson.ObjectIDFromHex(
		student.ClassroomID,
	)
	if err != nil {
		return errors.New("invalid classroom_id")
	}

	classroom, err := s.classroomRepo.GetByIDAndCollegeID(
		ctx,
		classroomObjectID,
		collegeID,
	)

	if err != nil {
		return errors.New(
			"classroom not found or does not belong to your school",
		)
	}

	exists, err :=
		s.studentRepo.ExistsByClassroomAndRollNumberExceptID(
			ctx,
			collegeID,
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

	if existingStudent.ClassroomID != student.ClassroomID {

		count, err := s.studentRepo.CountByClassroom(
			ctx,
			collegeID,
			student.ClassroomID,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to check classroom capacity: %w",
				err,
			)
		}

		if count >= int64(classroom.Capacity) {
			return fmt.Errorf(
				"%s - Section %s is full",
				classroom.Name,
				classroom.Section,
			)
		}
	}

	student.CollegeID = collegeID

	return s.studentRepo.UpdateByIDAndCollegeID(
		ctx,
		id,
		collegeID,
		student,
	)
}

func (s *StudentService) DeleteByCollegeID(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
) error {

	return s.studentRepo.DeleteByIDAndCollegeID(
		ctx,
		id,
		collegeID,
	)
}

func validateGuardians(guardians []studentModel.Guardian) error {

	if len(guardians) == 0 {
		return errors.New("at least one guardian is required")
	}

	primaryCount := 0
	reportReceiverCount := 0

	for i := range guardians {

		guardians[i].Name = strings.TrimSpace(guardians[i].Name)
		guardians[i].Relation = strings.TrimSpace(guardians[i].Relation)
		guardians[i].Phone = strings.TrimSpace(guardians[i].Phone)
		guardians[i].Email = strings.TrimSpace(guardians[i].Email)

		if guardians[i].Name == "" {
			return fmt.Errorf(
				"guardian %d name is required",
				i+1,
			)
		}

		if guardians[i].Relation == "" {
			return fmt.Errorf(
				"guardian %d relation is required",
				i+1,
			)
		}

		switch strings.ToLower(guardians[i].Relation) {
		case "father", "mother", "guardian", "other":
		default:
			return fmt.Errorf(
				"guardian %d relation must be Father, Mother, Guardian or Other",
				i+1,
			)
		}

		if guardians[i].Phone == "" &&
			guardians[i].Email == "" {

			return fmt.Errorf(
				"guardian %d must have phone or email",
				i+1,
			)
		}

		if guardians[i].Primary {
			primaryCount++
		}

		if guardians[i].ReceiveReport {
			reportReceiverCount++
		}
	}

	if primaryCount == 0 {
		return errors.New(
			"one primary guardian is required",
		)
	}

	if primaryCount > 1 {
		return errors.New(
			"only one guardian can be primary",
		)
	}

	if reportReceiverCount == 0 {
		return errors.New(
			"at least one guardian must receive reports",
		)
	}

	return nil
}
