package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Subject struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	CollegeID bson.ObjectID `json:"college_id" bson:"college_id"`

	Name   string `json:"name" bson:"name"`
	Code   string `json:"code" bson:"code"`
	Status string `json:"status" bson:"status"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type CreateSubjectRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type UpdateSubjectRequest struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Status string `json:"status"`
}
