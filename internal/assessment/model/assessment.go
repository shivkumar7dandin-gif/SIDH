package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Assessment struct {
	ID            bson.ObjectID `json:"id" bson:"_id,omitempty"`
	StudentID     bson.ObjectID `json:"student_id" bson:"student_id"`
	ClassroomID   string        `json:"classroom_id" bson:"classroom_id"`
	AssessmentNo  int           `json:"assessment_no" bson:"assessment_no"`
	Subject       string        `json:"subject" bson:"subject"`
	TotalMarks    float64       `json:"total_marks" bson:"total_marks"`
	ObtainedMarks float64       `json:"obtained_marks" bson:"obtained_marks"`
	Percentage    float64       `json:"percentage" bson:"percentage"`
	Result        string        `json:"result" bson:"result"`
}
