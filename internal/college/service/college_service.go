package service

import (
	"context"
	"fmt"
	"strings"

	collegeModel "github.com/shivkumar7dandin-gif/students-api/internal/college/model"
	collegeRepository "github.com/shivkumar7dandin-gif/students-api/internal/college/repository"

	classroomRepository "github.com/shivkumar7dandin-gif/students-api/internal/classroom/repository"
	studentRepository "github.com/shivkumar7dandin-gif/students-api/internal/student/repository"
	teacherRepository "github.com/shivkumar7dandin-gif/students-api/internal/teacher/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type CollegeService struct {
	collegeRepo   *collegeRepository.CollegeRepository
	teacherRepo   *teacherRepository.TeacherRepository
	classroomRepo *classroomRepository.ClassroomRepository
	studentRepo   *studentRepository.StudentRepository
}

func NewCollegeService(
	collegeRepo *collegeRepository.CollegeRepository,
	teacherRepo *teacherRepository.TeacherRepository,
	classroomRepo *classroomRepository.ClassroomRepository,
	studentRepo *studentRepository.StudentRepository,
) *CollegeService {

	return &CollegeService{
		collegeRepo:   collegeRepo,
		teacherRepo:   teacherRepo,
		classroomRepo: classroomRepo,
		studentRepo:   studentRepo,
	}
}

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

func (s *CollegeService) GetAll(
	ctx context.Context,
) ([]collegeModel.College, error) {

	return s.collegeRepo.GetAll(ctx)
}

func (s *CollegeService) GetByID(
	ctx context.Context,
	id bson.ObjectID,
) (*collegeModel.College, error) {

	return s.collegeRepo.GetByID(ctx, id)
}

func (s *CollegeService) GetDashboard(
	ctx context.Context,
	collegeID bson.ObjectID,
) (*collegeModel.CollegeDashboard, error) {

	// Get college
	college, err := s.collegeRepo.GetByID(ctx, collegeID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get college: %w",
			err,
		)
	}

	// Get teachers
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

	// Get classrooms
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

	// Get students
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

	return &collegeModel.CollegeDashboard{
		College:         college,
		Teachers:        teachers,
		Classrooms:      classrooms,
		Students:        students,
		TotalTeachers:   len(teachers),
		TotalClassrooms: len(classrooms),
		TotalStudents:   len(students),
	}, nil
}
