package service

import (
	"context"
	"errors"
	"strings"
	"time"

	subjectModel "github.com/shivkumar7dandin-gif/students-api/internal/subject/model"
	subjectRepository "github.com/shivkumar7dandin-gif/students-api/internal/subject/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SubjectService struct {
	repository *subjectRepository.SubjectRepository
}

func NewSubjectService(
	repository *subjectRepository.SubjectRepository,
) *SubjectService {
	return &SubjectService{
		repository: repository,
	}
}

func (s *SubjectService) Create(
	ctx context.Context,
	collegeID bson.ObjectID,
	req subjectModel.CreateSubjectRequest,
) (*subjectModel.Subject, error) {

	name := strings.TrimSpace(req.Name)
	code := strings.ToUpper(strings.TrimSpace(req.Code))

	if name == "" {
		return nil, errors.New("subject name is required")
	}

	if code == "" {
		return nil, errors.New("subject code is required")
	}

	exists, err := s.repository.ExistsByName(
		ctx,
		collegeID,
		name,
	)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("subject already exists")
	}

	now := time.Now()

	subject := &subjectModel.Subject{
		ID:        bson.NewObjectID(),
		CollegeID: collegeID,
		Name:      name,
		Code:      code,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repository.Create(ctx, subject); err != nil {
		return nil, err
	}

	return subject, nil
}

func (s *SubjectService) GetAll(
	ctx context.Context,
	collegeID bson.ObjectID,
) ([]subjectModel.Subject, error) {

	subjects, err := s.repository.GetAllByCollegeID(
		ctx,
		collegeID,
	)

	if err != nil {
		return nil, err
	}

	if subjects == nil {
		return []subjectModel.Subject{}, nil
	}

	return subjects, nil
}

func (s *SubjectService) GetByID(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
) (*subjectModel.Subject, error) {

	subject, err := s.repository.GetByIDAndCollegeID(
		ctx,
		id,
		collegeID,
	)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("subject not found")
		}

		return nil, err
	}

	return subject, nil
}

func (s *SubjectService) Delete(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
) error {

	err := s.repository.Delete(
		ctx,
		id,
		collegeID,
	)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("subject not found")
		}

		return err
	}

	return nil
}
