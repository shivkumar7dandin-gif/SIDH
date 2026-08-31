package model

type CreateCollegeRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	Pincode     string `json:"pincode"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`

	Username string `json:"username"`
	Password string `json:"password"`
}
