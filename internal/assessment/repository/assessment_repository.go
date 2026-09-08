package repository

import (
	"context"
	"fmt"

	"github.com/shivkumar7dandin-gif/students-api/internal/assessment/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AssessmentRepository struct {
	collection *mongo.Collection
}

func NewAssessmentRepository(db *mongo.Database) *AssessmentRepository {
	return &AssessmentRepository{
		collection: db.Collection("assessments"),
	}
}

func (r *AssessmentRepository) Create(
	ctx context.Context,
	assessment model.Assessment,
) (*model.Assessment, error) {

	result, err := r.collection.InsertOne(ctx, assessment)
	if err != nil {
		return nil, err
	}

	assessment.ID = result.InsertedID.(bson.ObjectID)

	return &assessment, nil
}

func (r *AssessmentRepository) GetAll(
	ctx context.Context,
) ([]model.Assessment, error) {

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assessments []model.Assessment

	if err := cursor.All(ctx, &assessments); err != nil {
		return nil, err
	}

	return assessments, nil
}

func (r *AssessmentRepository) GetByStudent(
	ctx context.Context,
	studentID bson.ObjectID,
) ([]model.Assessment, error) {

	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"student_id": studentID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assessments []model.Assessment

	if err := cursor.All(ctx, &assessments); err != nil {
		return nil, err
	}

	return assessments, nil
}

// ========================================
// GET ALL BY COLLEGE + ACADEMIC YEAR
// ========================================

func (r *AssessmentRepository) GetAllByCollegeAndYear(
	ctx context.Context,
	collegeID bson.ObjectID,
	academicYear string,
) ([]model.Assessment, error) {

	filter := bson.M{
		"college_id":    collegeID,
		"academic_year": academicYear,
	}

	cursor, err := r.collection.Find(
		ctx,
		filter,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var assessments []model.Assessment

	if err := cursor.All(
		ctx,
		&assessments,
	); err != nil {
		return nil, err
	}

	return assessments, nil
}

// ========================================
// GET BY STUDENT + COLLEGE + YEAR
// ========================================

func (r *AssessmentRepository) GetByStudentAndYear(
	ctx context.Context,
	collegeID bson.ObjectID,
	studentID bson.ObjectID,
	academicYear string,
) ([]model.Assessment, error) {

	filter := bson.M{
		"college_id":    collegeID,
		"student_id":    studentID,
		"academic_year": academicYear,
	}

	cursor, err := r.collection.Find(
		ctx,
		filter,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var assessments []model.Assessment

	if err := cursor.All(
		ctx,
		&assessments,
	); err != nil {
		return nil, err
	}

	return assessments, nil
}

func (r *AssessmentRepository) GetByStudentAndMonth(
	ctx context.Context,
	collegeID bson.ObjectID,
	studentID bson.ObjectID,
	academicYear string,
	year int,
	month int,
) ([]model.Assessment, error) {

	startDate := fmt.Sprintf(
		"%04d-%02d-01",
		year,
		month,
	)

	endYear := year
	endMonth := month + 1

	if month == 12 {
		endYear = year + 1
		endMonth = 1
	}

	endDate := fmt.Sprintf(
		"%04d-%02d-01",
		endYear,
		endMonth,
	)

	filter := bson.M{
		"college_id":    collegeID,
		"student_id":    studentID,
		"academic_year": academicYear,
		"assessment_date": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}

	cursor, err := r.collection.Find(
		ctx,
		filter,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var assessments []model.Assessment

	if err := cursor.All(
		ctx,
		&assessments,
	); err != nil {
		return nil, err
	}

	return assessments, nil
}
