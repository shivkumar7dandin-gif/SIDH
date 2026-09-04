package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/shivkumar7dandin-gif/students-api/internal/classroom/model"
	"github.com/shivkumar7dandin-gif/students-api/internal/classroom/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ClassroomService struct {
	repository *repository.ClassroomRepository
}

func NewClassroomService(
	repository *repository.ClassroomRepository,
) *ClassroomService {

	return &ClassroomService{
		repository: repository,
	}
}

func (s *ClassroomService) Create(
	ctx context.Context,
	req model.CreateClassroomRequest,
	collegeID bson.ObjectID,
) (*model.Classroom, error) {

	req.Name = strings.TrimSpace(req.Name)
	req.Section = strings.TrimSpace(req.Section)

	if req.Name == "" {
		return nil, fmt.Errorf("class name is required")
	}

	if req.Section == "" {
		return nil, fmt.Errorf("section is required")
	}

	if req.Capacity <= 0 {
		return nil, fmt.Errorf(
			"classroom capacity must be greater than 0",
		)
	}

	if req.Capacity > 60 {
		return nil, fmt.Errorf(
			"classroom capacity cannot be more than 60 students",
		)
	}

	exists, err := s.repository.ExistsByNameAndSection(
		ctx,
		collegeID,
		req.Name,
		req.Section,
	)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf(
			"%s - Section %s already exists",
			req.Name,
			req.Section,
		)
	}

	classroom := model.Classroom{
		CollegeID: collegeID,
		Name:      req.Name,
		Section:   req.Section,
		Capacity:  req.Capacity,
	}

	return s.repository.Create(ctx, classroom)
}

func (s *ClassroomService) GetAll(
	ctx context.Context,
) ([]model.Classroom, error) {

	return s.repository.GetAll(ctx)
}

func (s *ClassroomService) GetByCollegeID(
	ctx context.Context,
	collegeID bson.ObjectID,
) ([]model.Classroom, error) {

	return s.repository.GetByCollegeID(ctx, collegeID)
}

func (s *ClassroomService) GetByID(
	ctx context.Context,
	id bson.ObjectID,
) (*model.Classroom, error) {

	return s.repository.GetByID(ctx, id)
}

func (s *ClassroomService) Update(
	ctx context.Context,
	id bson.ObjectID,
	classroom model.Classroom,
) error {

	return s.repository.Update(ctx, id, classroom)
}

func (s *ClassroomService) Delete(
	ctx context.Context,
	id bson.ObjectID,
) error {

	return s.repository.Delete(ctx, id)
}
