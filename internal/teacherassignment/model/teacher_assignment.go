package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TeacherClassAssignment struct {
	ID bson.ObjectID `json:"id" bson:"_id,omitempty"`

	CollegeID bson.ObjectID `json:"college_id" bson:"college_id"`

	TeacherID bson.ObjectID `json:"teacher_id" bson:"teacher_id"`

	ClassroomID bson.ObjectID `json:"classroom_id" bson:"classroom_id"`

	Subject string `json:"subject" bson:"subject"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`

	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type CreateTeacherClassAssignmentRequest struct {
	TeacherID string `json:"teacher_id"`

	ClassroomID string `json:"classroom_id"`

	Subject string `json:"subject"`
}
