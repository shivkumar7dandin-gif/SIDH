package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/shivkumar7dandin-gif/students-api/internal/college/model"
	"github.com/shivkumar7dandin-gif/students-api/internal/college/repository"
	"go.mongodb.org/mongo-driver/v2/bson"

	"golang.org/x/crypto/bcrypt"
)

type CollegeService struct {
	repo *repository.CollegeRepository
}

func NewCollegeService(
	repo *repository.CollegeRepository,
) *CollegeService {
	return &CollegeService{
		repo: repo,
	}
}

func (s *CollegeService) Create(
	ctx context.Context,
	req model.CreateCollegeRequest,
) (*model.College, error) {

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

	college := model.College{
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

	return s.repo.Create(ctx, college)
}

func (s *CollegeService) GetAll(
	ctx context.Context,
) ([]model.College, error) {
	return s.repo.GetAll(ctx)
}

func (s *CollegeService) GetByID(
	ctx context.Context,
	id bson.ObjectID,
) (*model.College, error) {
	return s.repo.GetByID(ctx, id)
}
