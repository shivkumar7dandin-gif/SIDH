package repository

import (
	"context"

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
