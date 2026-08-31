package model

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID           bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string        `json:"username" bson:"username"`
	PasswordHash string        `json:"-" bson:"password_hash"`
	Role         string        `json:"role" bson:"role"`
	ReferenceID  bson.ObjectID `json:"reference_id" bson:"reference_id"`
}
