package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Assessment struct {
	ID             bson.ObjectID `json:"id" bson:"_id,omitempty"`
	CollegeID      bson.ObjectID `json:"college_id" bson:"college_id"`
	StudentID      bson.ObjectID `json:"student_id" bson:"student_id"`
	ClassroomID    string        `json:"classroom_id" bson:"classroom_id"`
	AcademicYear   string        `json:"academic_year" bson:"academic_year"`
	AssessmentDate string        `json:"assessment_date" bson:"assessment_date"`
	AssessmentNo   int           `json:"assessment_no" bson:"assessment_no"`
	Subject        string        `json:"subject" bson:"subject"`
	TotalMarks     float64       `json:"total_marks" bson:"total_marks"`
	ObtainedMarks  float64       `json:"obtained_marks" bson:"obtained_marks"`
	Percentage     float64       `json:"percentage" bson:"percentage"`
	Result         string        `json:"result" bson:"result"`
}

type MonthlyAssessmentSummary struct {
	StudentID        bson.ObjectID `json:"student_id"`
	AcademicYear     string        `json:"academic_year"`
	Year             int           `json:"year"`
	Month            int           `json:"month"`
	TotalAssessments int           `json:"total_assessments"`
	Passed           int           `json:"passed"`
	Failed           int           `json:"failed"`
	TotalMarks       float64       `json:"total_marks"`
	ObtainedMarks    float64       `json:"obtained_marks"`
	Percentage       float64       `json:"percentage"`
	Result           string        `json:"result"`
}
