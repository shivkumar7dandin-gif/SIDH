package repository

import (
	"context"

	"github.com/shivkumar7dandin-gif/students-api/internal/classroom/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ClassroomRepository struct {
	collection *mongo.Collection
}

func NewClassroomRepository(db *mongo.Database) *ClassroomRepository {
	return &ClassroomRepository{
		collection: db.Collection("classrooms"),
	}
}

// Create classroom
func (r *ClassroomRepository) Create(
	ctx context.Context,
	classroom model.Classroom,
) (*model.Classroom, error) {

	result, err := r.collection.InsertOne(ctx, classroom)
	if err != nil {
		return nil, err
	}

	// MongoDB generates the ObjectID
	classroom.ID = result.InsertedID.(bson.ObjectID)

	return &classroom, nil
}

// Get all classrooms
func (r *ClassroomRepository) GetAll(
	ctx context.Context,
) ([]model.Classroom, error) {

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var classrooms []model.Classroom

	if err := cursor.All(ctx, &classrooms); err != nil {
		return nil, err
	}

	return classrooms, nil
}

// Get classroom by ID
func (r *ClassroomRepository) GetByID(
	ctx context.Context,
	id bson.ObjectID,
) (*model.Classroom, error) {

	var classroom model.Classroom

	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&classroom)

	if err != nil {
		return nil, err
	}

	return &classroom, nil
}

// Update classroom
func (r *ClassroomRepository) Update(
	ctx context.Context,
	id bson.ObjectID,
	classroom model.Classroom,
) error {

	update := bson.M{
		"$set": bson.M{
			"name":     classroom.Name,
			"section":  classroom.Section,
			"capacity": classroom.Capacity,
		},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		update,
	)

	if err != nil {
		return err
	}

	// No classroom found
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// Delete classroom
func (r *ClassroomRepository) Delete(
	ctx context.Context,
	id bson.ObjectID,
) error {

	result, err := r.collection.DeleteOne(
		ctx,
		bson.M{"_id": id},
	)

	if err != nil {
		return err
	}

	// No classroom found
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *ClassroomRepository) GetByName(
	ctx context.Context,
	name string,
) (*model.Classroom, error) {

	var classroom model.Classroom

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"name": name,
		},
	).Decode(&classroom)

	if err != nil {
		return nil, err
	}

	return &classroom, nil
}

func (r *ClassroomRepository) GetByCollegeID(
	ctx context.Context,
	collegeID bson.ObjectID,
) ([]model.Classroom, error) {

	filter := bson.M{
		"college_id": collegeID,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var classrooms []model.Classroom

	if err := cursor.All(ctx, &classrooms); err != nil {
		return nil, err
	}

	return classrooms, nil
}
