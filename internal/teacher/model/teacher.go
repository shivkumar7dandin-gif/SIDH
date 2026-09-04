package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Teacher struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	CollegeID bson.ObjectID `json:"college_id" bson:"college_id"`
	Name      string        `json:"name" bson:"name"`
	Age       int           `json:"age" bson:"age"`
	Gender    string        `json:"gender" bson:"gender"`
	Email     string        `json:"email" bson:"email"`
	Phone     string        `json:"phone" bson:"phone"`
	Subject   string        `json:"subject" bson:"subject"`
	Username  string        `json:"username" bson:"username"`
}

type CreateTeacherRequest struct {
	CollegeID string `json:"college_id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Subject   string `json:"subject"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}
