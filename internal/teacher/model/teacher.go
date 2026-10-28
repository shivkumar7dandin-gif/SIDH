package model

import "go.mongodb.org/mongo-driver/v2/bson"

// Teacher represents a teacher belonging to one school.
type Teacher struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	CollegeID bson.ObjectID `json:"college_id" bson:"college_id"`

	// Basic details
	EmployeeID string `json:"employee_id" bson:"employee_id"`
	Name       string `json:"name" bson:"name"`
	Age        int    `json:"age" bson:"age"`
	Gender     string `json:"gender" bson:"gender"`

	// Contact details
	Email string `json:"email" bson:"email"`
	Phone string `json:"phone" bson:"phone"`

	// Professional details
	Subject         string `json:"subject" bson:"subject"`
	Qualification   string `json:"qualification" bson:"qualification"`
	ExperienceYears int    `json:"experience_years" bson:"experience_years"`
	JoiningDate     string `json:"joining_date" bson:"joining_date"`

	// Address
	Address TeacherAddress `json:"address" bson:"address"`

	// Login / account
	Username string `json:"username" bson:"username"`
	Status   string `json:"status" bson:"status"`
}

type TeacherAddress struct {
	HouseNo string `json:"house_no" bson:"house_no"`
	Street  string `json:"street" bson:"street"`
	Village string `json:"village" bson:"village"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	Pincode string `json:"pincode" bson:"pincode"`
}

type CreateTeacherRequest struct {
	CollegeID string `json:"college_id"`

	// Basic details
	EmployeeID string `json:"employee_id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`

	// Contact details
	Email string `json:"email"`
	Phone string `json:"phone"`

	// Professional details
	Subject         string `json:"subject"`
	Qualification   string `json:"qualification"`
	ExperienceYears int    `json:"experience_years"`
	JoiningDate     string `json:"joining_date"`

	// Address
	Address TeacherAddress `json:"address"`

	// Login
	Username string `json:"username"`
	Password string `json:"password"`
}
