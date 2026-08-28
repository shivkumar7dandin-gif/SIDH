package repository

import (
	"context"

	"github.com/shivkumar7dandin-gif/students-api/internal/attendance/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AttendanceRepository struct {
	collection *mongo.Collection
}

func NewAttendanceRepository(db *mongo.Database) *AttendanceRepository {
	return &AttendanceRepository{
		collection: db.Collection("attendance"),
	}
}

func (r *AttendanceRepository) Create(
	ctx context.Context,
	attendance model.Attendance,
) (*model.Attendance, error) {

	result, err := r.collection.InsertOne(ctx, attendance)
	if err != nil {
		return nil, err
	}

	attendance.ID = result.InsertedID.(bson.ObjectID)

	return &attendance, nil
}

func (r *AttendanceRepository) GetByStudent(
	ctx context.Context,
	studentID bson.ObjectID,
) ([]model.Attendance, error) {

	cursor, err := r.collection.Find(
		ctx,
		bson.M{"student_id": studentID},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var attendance []model.Attendance

	if err := cursor.All(ctx, &attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (r *AttendanceRepository) GetAll(
	ctx context.Context,
) ([]model.Attendance, error) {

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var attendance []model.Attendance

	if err := cursor.All(ctx, &attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (r *AttendanceRepository) GetByStudentAndDate(
	ctx context.Context,
	studentID bson.ObjectID,
	date string,
) (*model.Attendance, error) {

	var attendance model.Attendance

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"student_id": studentID,
			"date":       date,
		},
	).Decode(&attendance)

	if err != nil {
		return nil, err
	}

	return &attendance, nil
}
