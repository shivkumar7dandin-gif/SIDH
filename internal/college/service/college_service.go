package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/shivkumar7dandin-gif/students-api/internal/college/model"
	"github.com/shivkumar7dandin-gif/students-api/internal/college/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	college model.College,
) (*model.College, error) {

	if strings.TrimSpace(college.Name) == "" {
		return nil, fmt.Errorf("college name is required")
	}

	if strings.TrimSpace(college.Code) == "" {
		return nil, fmt.Errorf("college code is required")
	}

	if strings.TrimSpace(college.City) == "" {
		return nil, fmt.Errorf("city is required")
	}

	if strings.TrimSpace(college.State) == "" {
		return nil, fmt.Errorf("state is required")
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
