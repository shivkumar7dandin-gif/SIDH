package model

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID           bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string        `json:"username" bson:"username"`
	PasswordHash string        `json:"-" bson:"password_hash"`
	Role         string        `json:"role" bson:"role"`

	// Points to college / teacher / student record depending on role
	ReferenceID bson.ObjectID `json:"reference_id" bson:"reference_id"`

	// Always points to the school/college
	CollegeID bson.ObjectID `json:"college_id" bson:"college_id"`
}
