package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Attendance struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	StudentID   bson.ObjectID `json:"student_id" bson:"student_id"`
	ClassroomID string        `json:"classroom_id" bson:"classroom_id"`
	Date        string        `json:"date" bson:"date"`
	Attendance  string        `json:"attendance" bson:"attendance"`
}
