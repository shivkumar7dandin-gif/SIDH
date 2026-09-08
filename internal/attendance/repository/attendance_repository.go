package repository

import (
	"context"
	"fmt"

	"github.com/shivkumar7dandin-gif/students-api/internal/attendance/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AttendanceRepository struct {
	collection *mongo.Collection
}

func NewAttendanceRepository(
	db *mongo.Database,
) *AttendanceRepository {

	return &AttendanceRepository{
		collection: db.Collection("attendance"),
	}
}

// ========================================
// CREATE
// ========================================

func (r *AttendanceRepository) Create(
	ctx context.Context,
	attendance model.Attendance,
) (*model.Attendance, error) {

	result, err := r.collection.InsertOne(
		ctx,
		attendance,
	)
	if err != nil {
		return nil, err
	}

	attendance.ID = result.InsertedID.(bson.ObjectID)

	return &attendance, nil
}

// ========================================
// GET BY STUDENT
// ========================================

func (r *AttendanceRepository) GetByStudent(
	ctx context.Context,
	collegeID bson.ObjectID,
	studentID bson.ObjectID,
	academicYear string,
) ([]model.Attendance, error) {

	filter := bson.M{
		"college_id":    collegeID,
		"student_id":    studentID,
		"academic_year": academicYear,
	}

	cursor, err := r.collection.Find(
		ctx,
		filter,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var attendanceList []model.Attendance

	if err := cursor.All(
		ctx,
		&attendanceList,
	); err != nil {
		return nil, err
	}

	return attendanceList, nil
}

// ========================================
// GET ALL BY COLLEGE + ACADEMIC YEAR
// ========================================

func (r *AttendanceRepository) GetAll(
	ctx context.Context,
	collegeID bson.ObjectID,
	academicYear string,
) ([]model.Attendance, error) {

	filter := bson.M{
		"college_id":    collegeID,
		"academic_year": academicYear,
	}

	cursor, err := r.collection.Find(
		ctx,
		filter,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var attendanceList []model.Attendance

	if err := cursor.All(
		ctx,
		&attendanceList,
	); err != nil {
		return nil, err
	}

	return attendanceList, nil
}

// ========================================
// DUPLICATE CHECK
// ========================================

func (r *AttendanceRepository) GetByStudentAndDate(
	ctx context.Context,
	collegeID bson.ObjectID,
	studentID bson.ObjectID,
	academicYear string,
	date string,
) (*model.Attendance, error) {

	var attendance model.Attendance

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"college_id":    collegeID,
			"student_id":    studentID,
			"academic_year": academicYear,
			"date":          date,
		},
	).Decode(&attendance)

	if err != nil {
		return nil, err
	}

	return &attendance, nil
}

// ========================================
// GET MONTHLY ATTENDANCE BY STUDENT
// ========================================

func (r *AttendanceRepository) GetByStudentAndMonth(
	ctx context.Context,
	collegeID bson.ObjectID,
	studentID bson.ObjectID,
	academicYear string,
	year int,
	month int,
) ([]model.Attendance, error) {

	startDate := fmt.Sprintf(
		"%04d-%02d-01",
		year,
		month,
	)

	var endYear int
	var endMonth int

	if month == 12 {
		endYear = year + 1
		endMonth = 1
	} else {
		endYear = year
		endMonth = month + 1
	}

	endDate := fmt.Sprintf(
		"%04d-%02d-01",
		endYear,
		endMonth,
	)

	filter := bson.M{
		"college_id":    collegeID,
		"student_id":    studentID,
		"academic_year": academicYear,
		"date": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}

	cursor, err := r.collection.Find(
		ctx,
		filter,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var attendanceList []model.Attendance

	if err := cursor.All(
		ctx,
		&attendanceList,
	); err != nil {
		return nil, err
	}

	return attendanceList, nil
}
