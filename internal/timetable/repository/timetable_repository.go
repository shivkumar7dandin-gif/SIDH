package repository

import (
	"context"
	"fmt"

	timetableModel "github.com/shivkumar7dandin-gif/students-api/internal/timetable/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TimetableRepository struct {
	collection *mongo.Collection
}

func NewTimetableRepository(
	db *mongo.Database,
) *TimetableRepository {

	return &TimetableRepository{
		collection: db.Collection("timetables"),
	}
}

// ========================================
// CREATE TIMETABLE
// ========================================

func (r *TimetableRepository) Create(
	ctx context.Context,
	timetable timetableModel.Timetable,
) (*timetableModel.Timetable, error) {

	result, err := r.collection.InsertOne(
		ctx,
		timetable,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to create timetable: %w",
			err,
		)
	}

	id, ok := result.InsertedID.(bson.ObjectID)

	if !ok {
		return nil, fmt.Errorf(
			"failed to get timetable id",
		)
	}

	timetable.ID = id

	return &timetable, nil
}

// ========================================
// GET BY COLLEGE
// ========================================

func (r *TimetableRepository) GetByCollegeID(
	ctx context.Context,
	collegeID bson.ObjectID,
) ([]timetableModel.Timetable, error) {

	filter := bson.M{
		"college_id": collegeID,
	}

	cursor, err := r.collection.Find(
		ctx,
		filter,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get timetables: %w",
			err,
		)
	}

	defer cursor.Close(ctx)

	timetables := make(
		[]timetableModel.Timetable,
		0,
	)

	for cursor.Next(ctx) {

		var timetable timetableModel.Timetable

		if err := cursor.Decode(&timetable); err != nil {
			return nil, fmt.Errorf(
				"failed to decode timetable: %w",
				err,
			)
		}

		timetables = append(
			timetables,
			timetable,
		)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf(
			"timetable cursor error: %w",
			err,
		)
	}

	return timetables, nil
}

// ========================================
// GET BY COLLEGE + ACADEMIC YEAR
// ========================================

func (r *TimetableRepository) GetByCollegeAndAcademicYear(
	ctx context.Context,
	collegeID bson.ObjectID,
	academicYear string,
) ([]timetableModel.Timetable, error) {

	filter := bson.M{
		"college_id":    collegeID,
		"academic_year": academicYear,
	}

	cursor, err := r.collection.Find(
		ctx,
		filter,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get timetables: %w",
			err,
		)
	}

	defer cursor.Close(ctx)

	timetables := make(
		[]timetableModel.Timetable,
		0,
	)

	for cursor.Next(ctx) {

		var timetable timetableModel.Timetable

		if err := cursor.Decode(&timetable); err != nil {
			return nil, fmt.Errorf(
				"failed to decode timetable: %w",
				err,
			)
		}

		timetables = append(
			timetables,
			timetable,
		)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf(
			"timetable cursor error: %w",
			err,
		)
	}

	return timetables, nil
}

// ========================================
// GET BY COLLEGE + DAY
// ========================================

func (r *TimetableRepository) GetByCollegeAndDay(
	ctx context.Context,
	collegeID bson.ObjectID,
	dayOfWeek string,
	academicYear string,
) ([]timetableModel.Timetable, error) {

	filter := bson.M{
		"college_id":    collegeID,
		"day_of_week":   dayOfWeek,
		"academic_year": academicYear,
	}

	cursor, err := r.collection.Find(
		ctx,
		filter,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get timetable by day: %w",
			err,
		)
	}

	defer cursor.Close(ctx)

	timetables := make(
		[]timetableModel.Timetable,
		0,
	)

	for cursor.Next(ctx) {

		var timetable timetableModel.Timetable

		if err := cursor.Decode(&timetable); err != nil {
			return nil, fmt.Errorf(
				"failed to decode timetable: %w",
				err,
			)
		}

		timetables = append(
			timetables,
			timetable,
		)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf(
			"timetable cursor error: %w",
			err,
		)
	}

	return timetables, nil
}

// ========================================
// GET BY TEACHER + DAY
// ========================================

func (r *TimetableRepository) GetByTeacherAndDay(
	ctx context.Context,
	collegeID bson.ObjectID,
	teacherID bson.ObjectID,
	dayOfWeek string,
	academicYear string,
) ([]timetableModel.Timetable, error) {

	filter := bson.M{
		"college_id":    collegeID,
		"teacher_id":    teacherID,
		"day_of_week":   dayOfWeek,
		"academic_year": academicYear,
	}

	cursor, err := r.collection.Find(
		ctx,
		filter,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get teacher timetable: %w",
			err,
		)
	}

	defer cursor.Close(ctx)

	timetables := make(
		[]timetableModel.Timetable,
		0,
	)

	for cursor.Next(ctx) {

		var timetable timetableModel.Timetable

		if err := cursor.Decode(&timetable); err != nil {
			return nil, fmt.Errorf(
				"failed to decode timetable: %w",
				err,
			)
		}

		timetables = append(
			timetables,
			timetable,
		)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf(
			"timetable cursor error: %w",
			err,
		)
	}

	return timetables, nil
}

// ========================================
// GET BY CLASSROOM + DAY
// ========================================

func (r *TimetableRepository) GetByClassroomAndDay(
	ctx context.Context,
	collegeID bson.ObjectID,
	classroomID bson.ObjectID,
	dayOfWeek string,
	academicYear string,
) ([]timetableModel.Timetable, error) {

	filter := bson.M{
		"college_id":    collegeID,
		"classroom_id":  classroomID,
		"day_of_week":   dayOfWeek,
		"academic_year": academicYear,
	}

	cursor, err := r.collection.Find(
		ctx,
		filter,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get classroom timetable: %w",
			err,
		)
	}

	defer cursor.Close(ctx)

	timetables := make(
		[]timetableModel.Timetable,
		0,
	)

	for cursor.Next(ctx) {

		var timetable timetableModel.Timetable

		if err := cursor.Decode(&timetable); err != nil {
			return nil, fmt.Errorf(
				"failed to decode timetable: %w",
				err,
			)
		}

		timetables = append(
			timetables,
			timetable,
		)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf(
			"timetable cursor error: %w",
			err,
		)
	}

	return timetables, nil
}

// ========================================
// GET TIMETABLE BY ID
// ========================================

func (r *TimetableRepository) GetByID(
	ctx context.Context,
	collegeID bson.ObjectID,
	timetableID bson.ObjectID,
) (*timetableModel.Timetable, error) {

	filter := bson.M{
		"_id":        timetableID,
		"college_id": collegeID,
	}

	var timetable timetableModel.Timetable

	err := r.collection.FindOne(
		ctx,
		filter,
	).Decode(&timetable)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf(
				"timetable not found",
			)
		}

		return nil, fmt.Errorf(
			"failed to get timetable: %w",
			err,
		)
	}

	return &timetable, nil
}

// ========================================
// UPDATE TIMETABLE
// ========================================

func (r *TimetableRepository) Update(
	ctx context.Context,
	collegeID bson.ObjectID,
	timetableID bson.ObjectID,
	timetable timetableModel.Timetable,
) (*timetableModel.Timetable, error) {

	filter := bson.M{
		"_id":        timetableID,
		"college_id": collegeID,
	}

	update := bson.M{
		"$set": bson.M{
			"teacher_id":    timetable.TeacherID,
			"classroom_id":  timetable.ClassroomID,
			"subject":       timetable.Subject,
			"day_of_week":   timetable.DayOfWeek,
			"start_time":    timetable.StartTime,
			"end_time":      timetable.EndTime,
			"academic_year": timetable.AcademicYear,
		},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		filter,
		update,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to update timetable: %w",
			err,
		)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf(
			"timetable not found",
		)
	}

	updatedTimetable, err := r.GetByID(
		ctx,
		collegeID,
		timetableID,
	)

	if err != nil {
		return nil, err
	}

	return updatedTimetable, nil
}

// ========================================
// DELETE TIMETABLE
// ========================================

func (r *TimetableRepository) Delete(
	ctx context.Context,
	collegeID bson.ObjectID,
	timetableID bson.ObjectID,
) error {

	filter := bson.M{
		"_id":        timetableID,
		"college_id": collegeID,
	}

	result, err := r.collection.DeleteOne(
		ctx,
		filter,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to delete timetable: %w",
			err,
		)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf(
			"timetable not found",
		)
	}

	return nil
}
