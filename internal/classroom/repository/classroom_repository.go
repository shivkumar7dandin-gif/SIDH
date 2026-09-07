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

func NewClassroomRepository(
	db *mongo.Database,
) *ClassroomRepository {

	return &ClassroomRepository{
		collection: db.Collection("classrooms"),
	}
}

// =====================================================
// CREATE CLASSROOM
// =====================================================

func (r *ClassroomRepository) Create(
	ctx context.Context,
	classroom model.Classroom,
) (*model.Classroom, error) {

	result, err := r.collection.InsertOne(
		ctx,
		classroom,
	)

	if err != nil {
		return nil, err
	}

	classroom.ID = result.InsertedID.(bson.ObjectID)

	return &classroom, nil
}

// =====================================================
// GET ALL CLASSROOMS
// NOTE: Not tenant safe.
// Keep only if some internal/admin logic needs it.
// =====================================================

func (r *ClassroomRepository) GetAll(
	ctx context.Context,
) ([]model.Classroom, error) {

	cursor, err := r.collection.Find(
		ctx,
		bson.M{},
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	classrooms := make([]model.Classroom, 0)

	if err := cursor.All(
		ctx,
		&classrooms,
	); err != nil {

		return nil, err
	}

	return classrooms, nil
}

// =====================================================
// GET CLASSROOM BY ID
// NOTE: Not tenant safe.
// Keep because some internal service code may still use it.
// =====================================================

func (r *ClassroomRepository) GetByID(
	ctx context.Context,
	id bson.ObjectID,
) (*model.Classroom, error) {

	var classroom model.Classroom

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"_id": id,
		},
	).Decode(&classroom)

	if err != nil {
		return nil, err
	}

	return &classroom, nil
}

// =====================================================
// GET CLASSROOM BY ID + COLLEGE ID
// TENANT SAFE
// =====================================================

func (r *ClassroomRepository) GetByIDAndCollegeID(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
) (*model.Classroom, error) {

	var classroom model.Classroom

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"_id":        id,
			"college_id": collegeID,
		},
	).Decode(&classroom)

	if err != nil {
		return nil, err
	}

	return &classroom, nil
}

// =====================================================
// UPDATE CLASSROOM
// NOTE: Old non-tenant-safe function.
// Keep temporarily if some older service uses it.
// =====================================================

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

// =====================================================
// UPDATE CLASSROOM BY ID + COLLEGE ID
// TENANT SAFE
// =====================================================

func (r *ClassroomRepository) UpdateByIDAndCollegeID(
	ctx context.Context,
	id bson.ObjectID,
	collegeID bson.ObjectID,
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
		bson.M{
			"_id":        id,
			"college_id": collegeID,
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

// =====================================================
// DELETE CLASSROOM
// NOTE: Old non-tenant-safe function.
// =====================================================

func (r *ClassroomRepository) Delete(
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

// =====================================================
// DELETE CLASSROOM BY ID + COLLEGE ID
// TENANT SAFE
// =====================================================

func (r *ClassroomRepository) DeleteByIDAndCollegeID(
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
// GET CLASSROOM BY NAME
// NOTE: This is not tenant safe.
// Avoid using it for school-specific operations.
// =====================================================

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

// =====================================================
// GET ALL CLASSROOMS BY COLLEGE
// TENANT SAFE
// =====================================================

func (r *ClassroomRepository) GetByCollegeID(
	ctx context.Context,
	collegeID bson.ObjectID,
) ([]model.Classroom, error) {

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

	classrooms := make([]model.Classroom, 0)

	if err := cursor.All(
		ctx,
		&classrooms,
	); err != nil {

		return nil, err
	}

	return classrooms, nil
}

// =====================================================
// CHECK DUPLICATE CLASSROOM
// SAME SCHOOL + NAME + SECTION
// =====================================================

func (r *ClassroomRepository) ExistsByNameAndSection(
	ctx context.Context,
	collegeID bson.ObjectID,
	name string,
	section string,
) (bool, error) {

	filter := bson.M{
		"college_id": collegeID,
		"name":       name,
		"section":    section,
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
