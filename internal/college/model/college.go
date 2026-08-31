package model

import "go.mongodb.org/mongo-driver/v2/bson"

type College struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Code        string        `json:"code" bson:"code"`
	Address     string        `json:"address" bson:"address"`
	City        string        `json:"city" bson:"city"`
	State       string        `json:"state" bson:"state"`
	Pincode     string        `json:"pincode" bson:"pincode"`
	PhoneNumber string        `json:"phone_number" bson:"phone_number"`
	Email       string        `json:"email" bson:"email"`

	Username     string `json:"username" bson:"username"`
	PasswordHash string `json:"-" bson:"password_hash"`
	Role         string `json:"role" bson:"role"`
}
