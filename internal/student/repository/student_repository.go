package repository

import (
	"context"

	"github.com/shivkumar7dandin-gif/students-api/internal/student/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type StudentRepository struct {
	collection *mongo.Collection
}

func NewStudentRepository(db *mongo.Database) *StudentRepository {
	return &StudentRepository{
		collection: db.Collection("students"),
	}
}

// Create student
func (r *StudentRepository) Create(
	ctx context.Context,
	student model.Student,
) (*model.Student, error) {

	result, err := r.collection.InsertOne(ctx, student)

	if err != nil {
		return nil, err
	}

	student.ID = result.InsertedID.(bson.ObjectID)

	return &student, nil
}

// Get all students
func (r *StudentRepository) GetAll(
	ctx context.Context,
) ([]model.Student, error) {

	cursor, err := r.collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var students []model.Student

	if err := cursor.All(ctx, &students); err != nil {
		return nil, err
	}

	return students, nil
}

// Get student by ID
func (r *StudentRepository) GetByID(
	ctx context.Context,
	id bson.ObjectID,
) (*model.Student, error) {

	var student model.Student

	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&student)

	if err != nil {
		return nil, err
	}

	return &student, nil
}

// Update student
func (r *StudentRepository) Update(
	ctx context.Context,
	id bson.ObjectID,
	student model.Student,
) error {

	update := bson.M{
		"$set": bson.M{
			"name":         student.Name,
			"roll_number":  student.RollNumber,
			"age":          student.Age,
			"classroom_id": student.ClassroomID,
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

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// Delete student
func (r *StudentRepository) Delete(
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

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// Count students in a classroom
func (r *StudentRepository) CountByClassroom(
	ctx context.Context,
	classroomID string,
) (int64, error) {

	count, err := r.collection.CountDocuments(
		ctx,
		bson.M{
			"classroom_id": classroomID,
		},
	)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *StudentRepository) GetByCollegeID(
	ctx context.Context,
	collegeID bson.ObjectID,
) ([]model.Student, error) {

	filter := bson.M{
		"college_id": collegeID,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var students []model.Student

	if err := cursor.All(ctx, &students); err != nil {
		return nil, err
	}

	return students, nil
}
