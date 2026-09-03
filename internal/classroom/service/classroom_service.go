package service

import (
	"context"

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

// func (s *ClassroomService) Create(
// 	ctx context.Context,
// 	classroom model.Classroom,
// ) (*model.Classroom, error) {

// 	return s.repository.Create(ctx, classroom)
// }

func (s *ClassroomService) Create(
	ctx context.Context,
	req model.CreateClassroomRequest,
	collegeID bson.ObjectID,
) (*model.Classroom, error) {

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
