package service

import (
	"context"
	"errors"
	"strings"

	teacherModel "github.com/shivkumar7dandin-gif/students-api/internal/teacher/model"
	teacherRepository "github.com/shivkumar7dandin-gif/students-api/internal/teacher/repository"

	userModel "github.com/shivkumar7dandin-gif/students-api/internal/user/model"
	userRepository "github.com/shivkumar7dandin-gif/students-api/internal/user/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type TeacherService struct {
	teacherRepo *teacherRepository.TeacherRepository
	userRepo    *userRepository.UserRepository
}

func NewTeacherService(
	teacherRepo *teacherRepository.TeacherRepository,
	userRepo *userRepository.UserRepository,
) *TeacherService {

	return &TeacherService{
		teacherRepo: teacherRepo,
		userRepo:    userRepo,
	}
}

func (s *TeacherService) Create(
	ctx context.Context,
	req teacherModel.CreateTeacherRequest,
) (*teacherModel.Teacher, error) {

	req.Name = strings.TrimSpace(req.Name)
	req.Username = strings.TrimSpace(req.Username)
	req.Subject = strings.TrimSpace(req.Subject)

	if req.Name == "" {
		return nil, errors.New("teacher name is required")
	}

	if req.Username == "" {
		return nil, errors.New("username is required")
	}

	if len(req.Password) < 8 {
		return nil, errors.New(
			"password must contain at least 8 characters",
		)
	}

	collegeID, err := bson.ObjectIDFromHex(
		req.CollegeID,
	)
	if err != nil {
		return nil, errors.New("invalid college id")
	}

	exists, err := s.userRepo.UsernameExists(
		ctx,
		req.Username,
	)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New(
			"username already exists",
		)
	}

	teacher := teacherModel.Teacher{
		CollegeID: collegeID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Subject:   req.Subject,
		Username:  req.Username,
	}

	createdTeacher, err := s.teacherRepo.Create(
		ctx,
		teacher,
	)
	if err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := userModel.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		Role:         "teacher",
		ReferenceID:  createdTeacher.ID,
	}

	err = s.userRepo.Create(
		ctx,
		user,
	)
	if err != nil {
		return nil, err
	}

	return createdTeacher, nil
}

func (s *TeacherService) GetAll(
	ctx context.Context,
) ([]teacherModel.Teacher, error) {

	return s.teacherRepo.GetAll(ctx)
}

func (s *TeacherService) GetByID(
	ctx context.Context,
	id string,
) (*teacherModel.Teacher, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid teacher id")
	}

	return s.teacherRepo.GetByID(
		ctx,
		objectID,
	)
}

func (s *TeacherService) Update(
	ctx context.Context,
	id string,
	teacher teacherModel.Teacher,
) error {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid teacher id")
	}

	return s.teacherRepo.Update(
		ctx,
		objectID,
		teacher,
	)
}

func (s *TeacherService) Delete(
	ctx context.Context,
	id string,
) error {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid teacher id")
	}

	return s.teacherRepo.Delete(
		ctx,
		objectID,
	)
}
