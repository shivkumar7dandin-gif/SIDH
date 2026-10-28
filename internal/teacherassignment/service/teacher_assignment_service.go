package service

import (
	"context"
	"errors"
	"strings"
	"time"

	classroomRepository "github.com/shivkumar7dandin-gif/students-api/internal/classroom/repository"
	teacherRepository "github.com/shivkumar7dandin-gif/students-api/internal/teacher/repository"
	assignmentModel "github.com/shivkumar7dandin-gif/students-api/internal/teacherassignment/model"
	assignmentRepository "github.com/shivkumar7dandin-gif/students-api/internal/teacherassignment/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TeacherAssignmentService struct {
	assignmentRepo *assignmentRepository.TeacherAssignmentRepository
	teacherRepo    *teacherRepository.TeacherRepository
	classroomRepo  *classroomRepository.ClassroomRepository
}

func NewTeacherAssignmentService(
	assignmentRepo *assignmentRepository.TeacherAssignmentRepository,
	teacherRepo *teacherRepository.TeacherRepository,
	classroomRepo *classroomRepository.ClassroomRepository,
) *TeacherAssignmentService {

	return &TeacherAssignmentService{
		assignmentRepo: assignmentRepo,
		teacherRepo:    teacherRepo,
		classroomRepo:  classroomRepo,
	}
}

// =====================================================
// CREATE TEACHER CLASS ASSIGNMENT
// =====================================================

func (s *TeacherAssignmentService) Create(
	ctx context.Context,
	collegeID bson.ObjectID,
	req assignmentModel.CreateTeacherClassAssignmentRequest,
) (*assignmentModel.TeacherClassAssignment, error) {

	req.TeacherID = strings.TrimSpace(req.TeacherID)
	req.ClassroomID = strings.TrimSpace(req.ClassroomID)
	req.Subject = strings.TrimSpace(req.Subject)

	// -------------------------------------------------
	// VALIDATE REQUEST
	// -------------------------------------------------

	if req.TeacherID == "" {
		return nil, errors.New("teacher id is required")
	}

	if req.ClassroomID == "" {
		return nil, errors.New("classroom id is required")
	}

	if req.Subject == "" {
		return nil, errors.New("subject is required")
	}

	// -------------------------------------------------
	// CONVERT TEACHER ID
	// -------------------------------------------------

	teacherID, err := bson.ObjectIDFromHex(req.TeacherID)
	if err != nil {
		return nil, errors.New("invalid teacher id")
	}

	// -------------------------------------------------
	// CONVERT CLASSROOM ID
	// -------------------------------------------------

	classroomID, err := bson.ObjectIDFromHex(req.ClassroomID)
	if err != nil {
		return nil, errors.New("invalid classroom id")
	}

	// -------------------------------------------------
	// VERIFY TEACHER BELONGS TO JWT COLLEGE
	// -------------------------------------------------

	teacher, err := s.teacherRepo.GetByIDAndCollegeID(
		ctx,
		teacherID,
		collegeID,
	)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(
				"teacher not found in this school",
			)
		}

		return nil, err
	}

	if teacher.Status == "left" {
		return nil, errors.New(
			"cannot assign classroom to teacher who has left",
		)
	}

	// -------------------------------------------------
	// VERIFY CLASSROOM BELONGS TO JWT COLLEGE
	// -------------------------------------------------

	_, err = s.classroomRepo.GetByIDAndCollegeID(
		ctx,
		classroomID,
		collegeID,
	)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(
				"classroom not found in this school",
			)
		}

		return nil, err
	}

	// -------------------------------------------------
	// CHECK DUPLICATE ASSIGNMENT
	// -------------------------------------------------

	exists, err := s.assignmentRepo.Exists(
		ctx,
		collegeID,
		teacherID,
		classroomID,
		req.Subject,
	)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New(
			"teacher is already assigned to this classroom for this subject",
		)
	}

	// -------------------------------------------------
	// CREATE ASSIGNMENT
	// -------------------------------------------------

	now := time.Now()

	assignment := assignmentModel.TeacherClassAssignment{
		CollegeID:   collegeID,
		TeacherID:   teacherID,
		ClassroomID: classroomID,
		Subject:     req.Subject,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.assignmentRepo.Create(
		ctx,
		assignment,
	)
}

// =====================================================
// GET TEACHER ASSIGNMENTS
// TENANT SAFE
// =====================================================

func (s *TeacherAssignmentService) GetByTeacher(
	ctx context.Context,
	collegeID bson.ObjectID,
	teacherIDString string,
) ([]assignmentModel.TeacherClassAssignment, error) {

	teacherIDString = strings.TrimSpace(
		teacherIDString,
	)

	if teacherIDString == "" {
		return nil, errors.New(
			"teacher id is required",
		)
	}

	teacherID, err := bson.ObjectIDFromHex(
		teacherIDString,
	)

	if err != nil {
		return nil, errors.New(
			"invalid teacher id",
		)
	}

	// Verify teacher belongs to JWT school.
	_, err = s.teacherRepo.GetByIDAndCollegeID(
		ctx,
		teacherID,
		collegeID,
	)

	if err != nil {
		if errors.Is(
			err,
			mongo.ErrNoDocuments,
		) {
			return nil, errors.New(
				"teacher not found in this school",
			)
		}

		return nil, err
	}

	return s.assignmentRepo.
		GetByTeacherIDAndCollegeID(
			ctx,
			teacherID,
			collegeID,
		)
}

// =====================================================
// DELETE ASSIGNMENT
// TENANT SAFE
// =====================================================

func (s *TeacherAssignmentService) Delete(
	ctx context.Context,
	collegeID bson.ObjectID,
	assignmentIDString string,
) error {

	assignmentIDString =
		strings.TrimSpace(
			assignmentIDString,
		)

	if assignmentIDString == "" {
		return errors.New(
			"assignment id is required",
		)
	}

	assignmentID, err :=
		bson.ObjectIDFromHex(
			assignmentIDString,
		)

	if err != nil {
		return errors.New(
			"invalid assignment id",
		)
	}

	err = s.assignmentRepo.
		DeleteByIDAndCollegeID(
			ctx,
			assignmentID,
			collegeID,
		)

	if errors.Is(
		err,
		mongo.ErrNoDocuments,
	) {
		return errors.New(
			"teacher assignment not found",
		)
	}

	return err
}
