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

type Student struct {
	CollegeID   bson.ObjectID `json:"college_id" bson:"college_id"`
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Age         int           `json:"age" bson:"age"`
	RollNumber  int           `json:"roll_number" bson:"roll_number"`
	Gender      string        `json:"gender" bson:"gender"`
	ClassroomID string        `json:"classroom_id" bson:"classroom_id"`
	Address     Address       `json:"address" bson:"address"`

	Username string `json:"username"`
	Password string `json:"password"`
}
