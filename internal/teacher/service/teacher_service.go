package service

import (
	"context"
	"errors"
	"strings"
	"time"

	teacherModel "github.com/shivkumar7dandin-gif/students-api/internal/teacher/model"
	teacherRepository "github.com/shivkumar7dandin-gif/students-api/internal/teacher/repository"

	userModel "github.com/shivkumar7dandin-gif/students-api/internal/user/model"
	userRepository "github.com/shivkumar7dandin-gif/students-api/internal/user/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

	req.EmployeeID = strings.TrimSpace(req.EmployeeID)
	req.Name = strings.TrimSpace(req.Name)
	req.Gender = strings.TrimSpace(req.Gender)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Subject = strings.TrimSpace(req.Subject)
	req.Qualification = strings.TrimSpace(req.Qualification)
	req.JoiningDate = strings.TrimSpace(req.JoiningDate)
	req.Username = strings.TrimSpace(req.Username)

	req.Address.HouseNo = strings.TrimSpace(req.Address.HouseNo)
	req.Address.Street = strings.TrimSpace(req.Address.Street)
	req.Address.Village = strings.TrimSpace(req.Address.Village)
	req.Address.City = strings.TrimSpace(req.Address.City)
	req.Address.State = strings.TrimSpace(req.Address.State)
	req.Address.Pincode = strings.TrimSpace(req.Address.Pincode)

	if req.Name == "" {
		return nil, errors.New("teacher name is required")
	}

	if req.EmployeeID == "" {
		return nil, errors.New("employee id is required")
	}

	if req.Qualification == "" {
		return nil, errors.New("teacher qualification is required")
	}

	if req.ExperienceYears < 0 {
		return nil, errors.New(
			"experience years cannot be negative",
		)
	}

	if req.JoiningDate == "" {
		return nil, errors.New(
			"joining date is required",
		)
	}

	if _, err := time.Parse(
		"2006-01-02",
		req.JoiningDate,
	); err != nil {
		return nil, errors.New(
			"joining date must be in YYYY-MM-DD format",
		)
	}

	if req.Address.City == "" {
		return nil, errors.New(
			"teacher address city is required",
		)
	}

	if req.Address.State == "" {
		return nil, errors.New(
			"teacher address state is required",
		)
	}

	if req.Address.Pincode == "" {
		return nil, errors.New(
			"teacher address pincode is required",
		)
	}

	if req.Age <= 0 {
		return nil, errors.New(
			"teacher age must be greater than 0",
		)
	}

	if req.Gender == "" {
		return nil, errors.New(
			"teacher gender is required",
		)
	}

	if req.Email == "" {
		return nil, errors.New(
			"teacher email is required",
		)
	}

	if req.Phone == "" {
		return nil, errors.New(
			"teacher phone is required",
		)
	}

	if req.Subject == "" {
		return nil, errors.New(
			"teacher subject is required",
		)
	}

	if req.Username == "" {
		return nil, errors.New(
			"username is required",
		)
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
		return nil, errors.New(
			"invalid college id",
		)
	}

	existingTeacher, err := s.teacherRepo.GetByEmployeeIDAndCollegeID(
		ctx,
		req.EmployeeID,
		collegeID,
	)

	if err == nil && existingTeacher != nil {
		return nil, errors.New(
			"employee id already exists in this school",
		)
	}

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
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
		CollegeID:       collegeID,
		EmployeeID:      req.EmployeeID,
		Name:            req.Name,
		Age:             req.Age,
		Gender:          req.Gender,
		Email:           req.Email,
		Phone:           req.Phone,
		Subject:         req.Subject,
		Qualification:   req.Qualification,
		ExperienceYears: req.ExperienceYears,
		JoiningDate:     req.JoiningDate,
		Address:         req.Address,
		Username:        req.Username,
		Status:          "active",
	}

	createdTeacher, err :=
		s.teacherRepo.Create(
			ctx,
			teacher,
		)

	if err != nil {
		return nil, err
	}

	passwordHash, err :=
		bcrypt.GenerateFromPassword(
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
		CollegeID:    createdTeacher.CollegeID,
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

	objectID, err :=
		bson.ObjectIDFromHex(id)

	if err != nil {
		return nil, errors.New(
			"invalid teacher id",
		)
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

	objectID, err :=
		bson.ObjectIDFromHex(id)

	if err != nil {
		return errors.New(
			"invalid teacher id",
		)
	}

	teacher.Name =
		strings.TrimSpace(
			teacher.Name,
		)

	teacher.Gender =
		strings.TrimSpace(
			teacher.Gender,
		)

	teacher.Email =
		strings.TrimSpace(
			teacher.Email,
		)

	teacher.Phone =
		strings.TrimSpace(
			teacher.Phone,
		)

	teacher.Subject =
		strings.TrimSpace(
			teacher.Subject,
		)

	if teacher.Name == "" {
		return errors.New(
			"teacher name is required",
		)
	}

	if teacher.Age <= 0 {
		return errors.New(
			"teacher age must be greater than 0",
		)
	}

	if teacher.Gender == "" {
		return errors.New(
			"teacher gender is required",
		)
	}

	if teacher.Subject == "" {
		return errors.New(
			"teacher subject is required",
		)
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

	objectID, err :=
		bson.ObjectIDFromHex(id)

	if err != nil {
		return errors.New(
			"invalid teacher id",
		)
	}

	return s.teacherRepo.Delete(
		ctx,
		objectID,
	)
}

func (s *TeacherService) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
) error {

	objectID, err :=
		bson.ObjectIDFromHex(id)

	if err != nil {
		return errors.New(
			"invalid teacher id",
		)
	}

	status = strings.ToLower(
		strings.TrimSpace(status),
	)

	if status != "active" && status != "left" {
		return errors.New(
			"status must be active or left",
		)
	}

	return s.teacherRepo.UpdateStatus(
		ctx,
		objectID,
		status,
	)
}

func (s *TeacherService) GetByCollegeID(
	ctx context.Context,
	collegeID bson.ObjectID,
) ([]teacherModel.Teacher, error) {

	return s.teacherRepo.GetByCollegeID(
		ctx,
		collegeID,
	)
}

func (s *TeacherService) GetByIDAndCollegeID(
	ctx context.Context,
	id string,
	collegeID bson.ObjectID,
) (*teacherModel.Teacher, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(
			"invalid teacher id",
		)
	}

	return s.teacherRepo.GetByIDAndCollegeID(
		ctx,
		objectID,
		collegeID,
	)
}

func (s *TeacherService) UpdateByCollegeID(
	ctx context.Context,
	id string,
	collegeID bson.ObjectID,
	teacher teacherModel.Teacher,
) error {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid teacher id")
	}

	teacher.Name = strings.TrimSpace(teacher.Name)
	teacher.Gender = strings.TrimSpace(teacher.Gender)
	teacher.Email = strings.TrimSpace(teacher.Email)
	teacher.Phone = strings.TrimSpace(teacher.Phone)
	teacher.Subject = strings.TrimSpace(teacher.Subject)

	if teacher.Name == "" {
		return errors.New("teacher name is required")
	}

	if teacher.Age <= 0 {
		return errors.New("teacher age must be greater than 0")
	}

	if teacher.Gender == "" {
		return errors.New("teacher gender is required")
	}

	if teacher.Email == "" {
		return errors.New("teacher email is required")
	}

	if teacher.Phone == "" {
		return errors.New("teacher phone is required")
	}

	if teacher.Subject == "" {
		return errors.New("teacher subject is required")
	}

	return s.teacherRepo.UpdateByIDAndCollegeID(
		ctx,
		objectID,
		collegeID,
		teacher,
	)
}

func (s *TeacherService) DeleteByCollegeID(
	ctx context.Context,
	id string,
	collegeID bson.ObjectID,
) error {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid teacher id")
	}

	return s.teacherRepo.DeleteByIDAndCollegeID(
		ctx,
		objectID,
		collegeID,
	)
}

func (s *TeacherService) UpdateStatusByCollegeID(
	ctx context.Context,
	id string,
	collegeID bson.ObjectID,
	status string,
) error {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid teacher id")
	}

	status = strings.ToLower(
		strings.TrimSpace(status),
	)

	if status != "active" && status != "left" {
		return errors.New(
			"status must be active or left",
		)
	}

	return s.teacherRepo.UpdateStatusByIDAndCollegeID(
		ctx,
		objectID,
		collegeID,
		status,
	)
}
