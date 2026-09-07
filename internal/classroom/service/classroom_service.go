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

// =====================================================
// CREATE CLASSROOM
// =====================================================

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
		return nil, fmt.Errorf(
			"failed to check classroom: %w",
			err,
		)
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

	return s.repository.Create(
		ctx,
		classroom,
	)
}

// =====================================================
// GET ALL CLASSROOMS BY COLLEGE
// TENANT SAFE
// =====================================================

func (s *ClassroomService) GetByCollegeID(
	ctx context.Context,
	collegeID bson.ObjectID,
) ([]model.Classroom, error) {

	return s.repository.GetByCollegeID(
		ctx,
		collegeID,
	)
}

// =====================================================
// GET CLASSROOM BY ID + COLLEGE ID
// TENANT SAFE
// =====================================================

func (s *ClassroomService) GetByIDAndCollegeID(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
) (*model.Classroom, error) {

	return s.repository.GetByIDAndCollegeID(
		ctx,
		id,
		collegeID,
	)
}

// =====================================================
// UPDATE CLASSROOM
// TENANT SAFE
// =====================================================

func (s *ClassroomService) UpdateByCollegeID(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
	classroom model.Classroom,
) error {

	classroom.Name = strings.TrimSpace(
		classroom.Name,
	)

	classroom.Section = strings.TrimSpace(
		classroom.Section,
	)

	if classroom.Name == "" {
		return fmt.Errorf("class name is required")
	}

	if classroom.Section == "" {
		return fmt.Errorf("section is required")
	}

	if classroom.Capacity <= 0 {
		return fmt.Errorf(
			"classroom capacity must be greater than 0",
		)
	}

	if classroom.Capacity > 60 {
		return fmt.Errorf(
			"classroom capacity cannot be more than 60 students",
		)
	}

	// Check classroom belongs to this college
	_, err := s.repository.GetByIDAndCollegeID(
		ctx,
		id,
		collegeID,
	)

	if err != nil {
		return fmt.Errorf("classroom not found")
	}

	return s.repository.UpdateByIDAndCollegeID(
		ctx,
		id,
		collegeID,
		classroom,
	)
}

// =====================================================
// DELETE CLASSROOM
// TENANT SAFE
// =====================================================

func (s *ClassroomService) DeleteByCollegeID(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
) error {

	return s.repository.DeleteByIDAndCollegeID(
		ctx,
		id,
		collegeID,
	)
}
