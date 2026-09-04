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

	students := make([]model.Student, 0)

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
		bson.M{
			"_id": id,
		},
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
			"gender":       student.Gender,
			"classroom_id": student.ClassroomID,
			"address":      student.Address,
		},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": id,
		},
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
		bson.M{
			"_id": id,
		},
	)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// Count students inside one classroom
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

// Check duplicate roll number inside same classroom
func (r *StudentRepository) ExistsByClassroomAndRollNumber(
	ctx context.Context,
	collegeID bson.ObjectID,
	classroomID string,
	rollNumber int,
) (bool, error) {

	filter := bson.M{
		"college_id":   collegeID,
		"classroom_id": classroomID,
		"roll_number":  rollNumber,
	}

	count, err := r.collection.CountDocuments(
		ctx,
		filter,
	)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Get students belonging to one college
func (r *StudentRepository) GetByCollegeID(
	ctx context.Context,
	collegeID bson.ObjectID,
) ([]model.Student, error) {

	filter := bson.M{
		"college_id": collegeID,
	}

	cursor, err := r.collection.Find(
		ctx,
		filter,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	students := make([]model.Student, 0)

	if err := cursor.All(
		ctx,
		&students,
	); err != nil {
		return nil, err
	}

	return students, nil
}

func (r *StudentRepository) ExistsByClassroomAndRollNumberExceptID(
	ctx context.Context,
	collegeID bson.ObjectID,
	classroomID string,
	rollNumber int,
	studentID bson.ObjectID,
) (bool, error) {

	filter := bson.M{
		"college_id":   collegeID,
		"classroom_id": classroomID,
		"roll_number":  rollNumber,

		// Ignore the student currently being edited
		"_id": bson.M{
			"$ne": studentID,
		},
	}

	count, err := r.collection.CountDocuments(
		ctx,
		filter,
	)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
