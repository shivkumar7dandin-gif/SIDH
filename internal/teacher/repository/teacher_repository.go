package repository

import (
	"context"

	"github.com/shivkumar7dandin-gif/students-api/internal/teacher/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TeacherRepository struct {
	collection *mongo.Collection
}

func NewTeacherRepository(
	db *mongo.Database,
) *TeacherRepository {

	return &TeacherRepository{
		collection: db.Collection("teachers"),
	}
}

func (r *TeacherRepository) Create(
	ctx context.Context,
	teacher model.Teacher,
) (*model.Teacher, error) {

	result, err := r.collection.InsertOne(
		ctx,
		teacher,
	)

	if err != nil {
		return nil, err
	}

	teacher.ID = result.InsertedID.(bson.ObjectID)

	return &teacher, nil
}

func (r *TeacherRepository) GetAll(
	ctx context.Context,
) ([]model.Teacher, error) {

	cursor, err := r.collection.Find(
		ctx,
		bson.M{},
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	teachers := make([]model.Teacher, 0)

	if err := cursor.All(
		ctx,
		&teachers,
	); err != nil {
		return nil, err
	}

	return teachers, nil
}

func (r *TeacherRepository) GetByID(
	ctx context.Context,
	id bson.ObjectID,
) (*model.Teacher, error) {

	var teacher model.Teacher

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"_id": id,
		},
	).Decode(&teacher)

	if err != nil {
		return nil, err
	}

	return &teacher, nil
}

func (r *TeacherRepository) Update(
	ctx context.Context,
	id bson.ObjectID,
	teacher model.Teacher,
) error {

	update := bson.M{
		"$set": bson.M{
			"name":    teacher.Name,
			"age":     teacher.Age,
			"gender":  teacher.Gender,
			"email":   teacher.Email,
			"phone":   teacher.Phone,
			"subject": teacher.Subject,
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

func (r *TeacherRepository) Delete(
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

func (r *TeacherRepository) GetByCollegeID(
	ctx context.Context,
	collegeID bson.ObjectID,
) ([]model.Teacher, error) {

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

	teachers := make([]model.Teacher, 0)

	if err := cursor.All(
		ctx,
		&teachers,
	); err != nil {
		return nil, err
	}

	return teachers, nil
}
