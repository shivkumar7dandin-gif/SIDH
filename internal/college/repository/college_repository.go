package repository

import (
	"context"

	"github.com/shivkumar7dandin-gif/students-api/internal/college/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CollegeRepository struct {
	collection *mongo.Collection
}

func NewCollegeRepository(db *mongo.Database) *CollegeRepository {
	return &CollegeRepository{
		collection: db.Collection("colleges"),
	}
}

func (r *CollegeRepository) Create(
	ctx context.Context,
	college model.College,
) (*model.College, error) {

	result, err := r.collection.InsertOne(ctx, college)
	if err != nil {
		return nil, err
	}

	college.ID = result.InsertedID.(bson.ObjectID)

	return &college, nil
}

func (r *CollegeRepository) GetAll(
	ctx context.Context,
) ([]model.College, error) {

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var colleges []model.College

	if err := cursor.All(ctx, &colleges); err != nil {
		return nil, err
	}

	return colleges, nil
}

func (r *CollegeRepository) GetByID(
	ctx context.Context,
	id bson.ObjectID,
) (*model.College, error) {

	var college model.College

	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&college)

	if err != nil {
		return nil, err
	}

	return &college, nil
}
