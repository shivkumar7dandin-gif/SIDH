package repository

import (
	"context"

	"github.com/shivkumar7dandin-gif/students-api/internal/calendar/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CalendarRepository struct {
	collection *mongo.Collection
}

func NewCalendarRepository(
	db *mongo.Database,
) *CalendarRepository {

	return &CalendarRepository{
		collection: db.Collection("academic_calendars"),
	}
}

// Create academic calendar
func (r *CalendarRepository) Create(
	ctx context.Context,
	calendar model.AcademicCalendar,
) (*model.AcademicCalendar, error) {

	result, err := r.collection.InsertOne(
		ctx,
		calendar,
	)

	if err != nil {
		return nil, err
	}

	calendar.ID = result.InsertedID.(bson.ObjectID)

	return &calendar, nil
}

// Get calendar by college
func (r *CalendarRepository) GetByCollegeID(
	ctx context.Context,
	collegeID bson.ObjectID,
) ([]model.AcademicCalendar, error) {

	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"college_id": collegeID,
		},
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	calendars := make([]model.AcademicCalendar, 0)

	if err := cursor.All(
		ctx,
		&calendars,
	); err != nil {
		return nil, err
	}

	return calendars, nil
}

// Get calendar by academic year
func (r *CalendarRepository) GetByAcademicYear(
	ctx context.Context,
	collegeID bson.ObjectID,
	academicYear string,
) (*model.AcademicCalendar, error) {

	var calendar model.AcademicCalendar

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"college_id":    collegeID,
			"academic_year": academicYear,
		},
	).Decode(&calendar)

	if err != nil {
		return nil, err
	}

	return &calendar, nil
}

// Check duplicate academic year for same school
func (r *CalendarRepository) ExistsByAcademicYear(
	ctx context.Context,
	collegeID bson.ObjectID,
	academicYear string,
) (bool, error) {

	count, err := r.collection.CountDocuments(
		ctx,
		bson.M{
			"college_id":    collegeID,
			"academic_year": academicYear,
		},
	)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
