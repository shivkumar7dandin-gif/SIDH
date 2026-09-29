package model

import "go.mongodb.org/mongo-driver/v2/bson"

// Timetable represents one recurring class period.
//
// Example:
// Shiva teaches Science
// Monday
// 8th Standard - A
// 09:00 to 10:00
type Timetable struct {
	ID bson.ObjectID `json:"id" bson:"_id,omitempty"`

	// Multi-school tenant
	CollegeID bson.ObjectID `json:"college_id" bson:"college_id"`

	// Teacher taking the class
	TeacherID bson.ObjectID `json:"teacher_id" bson:"teacher_id"`

	// Classroom
	ClassroomID bson.ObjectID `json:"classroom_id" bson:"classroom_id"`

	// Subject
	Subject string `json:"subject" bson:"subject"`

	// monday, tuesday, wednesday...
	DayOfWeek string `json:"day_of_week" bson:"day_of_week"`

	// HH:MM format
	// Example: 09:00
	StartTime string `json:"start_time" bson:"start_time"`

	// HH:MM format
	// Example: 10:00
	EndTime string `json:"end_time" bson:"end_time"`

	// Example: 2026-2027
	AcademicYear string `json:"academic_year" bson:"academic_year"`
}

// CreateTimetableRequest is used when college admin
// creates one timetable period.
type CreateTimetableRequest struct {
	TeacherID string `json:"teacher_id"`

	ClassroomID string `json:"classroom_id"`

	Subject string `json:"subject"`

	DayOfWeek string `json:"day_of_week"`

	StartTime string `json:"start_time"`

	EndTime string `json:"end_time"`

	AcademicYear string `json:"academic_year"`
}

type UpdateTimetableRequest struct {
	TeacherID    string `json:"teacher_id"`
	ClassroomID  string `json:"classroom_id"`
	Subject      string `json:"subject"`
	DayOfWeek    string `json:"day_of_week"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	AcademicYear string `json:"academic_year"`
}
