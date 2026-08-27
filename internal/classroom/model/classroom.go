package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Classroom struct {
	ID       bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string        `json:"name" bson:"name"`
	Section  string        `json:"section" bson:"section"`
	Capacity int           `json:"capacity" bson:"capacity"`
}
