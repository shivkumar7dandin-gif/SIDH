package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Address struct {
	HouseNo string `json:"house_no" bson:"house_no"`
	Street  string `json:"street" bson:"street"`
	Village string `json:"village" bson:"village"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	Pincode string `json:"pincode" bson:"pincode"`
}

type Guardian struct {
	Name          string `json:"name" bson:"name"`
	Relation      string `json:"relation" bson:"relation"`
	Phone         string `json:"phone" bson:"phone"`
	Email         string `json:"email" bson:"email"`
	Primary       bool   `json:"primary" bson:"primary"`
	ReceiveReport bool   `json:"receive_report" bson:"receive_report"`
}

type Student struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	CollegeID   bson.ObjectID `json:"college_id" bson:"college_id"`
	Name        string        `json:"name" bson:"name"`
	Age         int           `json:"age" bson:"age"`
	RollNumber  int           `json:"roll_number" bson:"roll_number"`
	Gender      string        `json:"gender" bson:"gender"`
	ClassroomID string        `json:"classroom_id" bson:"classroom_id"`
	Address     Address       `json:"address" bson:"address"`
	Guardians   []Guardian    `json:"guardians" bson:"guardians"`
}
