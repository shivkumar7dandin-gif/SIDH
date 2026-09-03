package model

type CreateStudentRequest struct {
	Name        string  `json:"name"`
	Age         int     `json:"age"`
	RollNumber  int     `json:"roll_number"`
	Gender      string  `json:"gender"`
	ClassroomID string  `json:"classroom_id"`
	Address     Address `json:"address"`

	Username string `json:"username"`
	Password string `json:"password"`
}
