package repository

import (
	"context"

	"github.com/shivkumar7dandin-gif/students-api/internal/teacherassignment/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TeacherAssignmentRepository struct {
	collection *mongo.Collection
}

func NewTeacherAssignmentRepository(
	db *mongo.Database,
) *TeacherAssignmentRepository {

	return &TeacherAssignmentRepository{
		collection: db.Collection("teacher_class_assignments"),
	}
}

// =====================================================
// CREATE ASSIGNMENT
// =====================================================

func (r *TeacherAssignmentRepository) Create(
	ctx context.Context,
	assignment model.TeacherClassAssignment,
) (*model.TeacherClassAssignment, error) {

	result, err := r.collection.InsertOne(
		ctx,
		assignment,
	)

	if err != nil {
		return nil, err
	}

	assignment.ID = result.InsertedID.(bson.ObjectID)

	return &assignment, nil
}

// =====================================================
// GET ASSIGNMENT BY ID + COLLEGE ID
// TENANT SAFE
// =====================================================

func (r *TeacherAssignmentRepository) GetByIDAndCollegeID(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
) (*model.TeacherClassAssignment, error) {

	var assignment model.TeacherClassAssignment

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"_id":        id,
			"college_id": collegeID,
		},
	).Decode(&assignment)

	if err != nil {
		return nil, err
	}

	return &assignment, nil
}

// =====================================================
// GET ALL ASSIGNMENTS FOR A TEACHER
// TENANT SAFE
// =====================================================

func (r *TeacherAssignmentRepository) GetByTeacherIDAndCollegeID(
	ctx context.Context,
	teacherID bson.ObjectID,
	collegeID bson.ObjectID,
) ([]model.TeacherClassAssignment, error) {

	filter := bson.M{
		"teacher_id": teacherID,
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

	assignments := make(
		[]model.TeacherClassAssignment,
		0,
	)

	if err := cursor.All(
		ctx,
		&assignments,
	); err != nil {
		return nil, err
	}

	return assignments, nil
}

// =====================================================
// GET ALL ASSIGNMENTS FOR A CLASSROOM
// TENANT SAFE
// Useful later for timetable/subject management.
// =====================================================

func (r *TeacherAssignmentRepository) GetByClassroomIDAndCollegeID(
	ctx context.Context,
	classroomID bson.ObjectID,
	collegeID bson.ObjectID,
) ([]model.TeacherClassAssignment, error) {

	filter := bson.M{
		"classroom_id": classroomID,
		"college_id":   collegeID,
	}

	cursor, err := r.collection.Find(
		ctx,
		filter,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	assignments := make(
		[]model.TeacherClassAssignment,
		0,
	)

	if err := cursor.All(
		ctx,
		&assignments,
	); err != nil {
		return nil, err
	}

	return assignments, nil
}

// =====================================================
// CHECK DUPLICATE ASSIGNMENT
//
// Same school + same teacher + same classroom +
// same subject should not be inserted twice.
// =====================================================

func (r *TeacherAssignmentRepository) Exists(
	ctx context.Context,
	collegeID bson.ObjectID,
	teacherID bson.ObjectID,
	classroomID bson.ObjectID,
	subject string,
) (bool, error) {

	filter := bson.M{
		"college_id":   collegeID,
		"teacher_id":   teacherID,
		"classroom_id": classroomID,
		"subject":      subject,
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

// =====================================================
// UPDATE SUBJECT
// TENANT SAFE
// =====================================================

func (r *TeacherAssignmentRepository) UpdateSubject(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
	subject string,
) error {

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id":        id,
			"college_id": collegeID,
		},
		bson.M{
			"$set": bson.M{
				"subject": subject,
			},
		},
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// =====================================================
// DELETE ASSIGNMENT
// TENANT SAFE
// =====================================================

func (r *TeacherAssignmentRepository) DeleteByIDAndCollegeID(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
) error {

	result, err := r.collection.DeleteOne(
		ctx,
		bson.M{
			"_id":        id,
			"college_id": collegeID,
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

// =====================================================
// DELETE ALL ASSIGNMENTS FOR A TEACHER
// TENANT SAFE
//
// Useful later if teacher leaves or is permanently
// removed from the school.
// =====================================================

func (r *TeacherAssignmentRepository) DeleteByTeacherIDAndCollegeID(
	ctx context.Context,
	teacherID bson.ObjectID,
	collegeID bson.ObjectID,
) error {

	_, err := r.collection.DeleteMany(
		ctx,
		bson.M{
			"teacher_id": teacherID,
			"college_id": collegeID,
		},
	)

	return err
}
