package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Student struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Age         int           `json:"age" bson:"age"`
	RollNumber  int           `json:"roll_number" bson:"roll_number"`
	ClassroomID string        `json:"classroom_id" bson:"classroom_id"`
}
