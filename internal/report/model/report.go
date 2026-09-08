package model

import (
	assessmentModel "github.com/shivkumar7dandin-gif/students-api/internal/assessment/model"
	attendanceModel "github.com/shivkumar7dandin-gif/students-api/internal/attendance/model"
	studentModel "github.com/shivkumar7dandin-gif/students-api/internal/student/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ReportRecipient struct {
	Name     string `json:"name"`
	Relation string `json:"relation"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
}

type MonthlyParentReport struct {
	StudentID    bson.ObjectID `json:"student_id"`
	StudentName  string        `json:"student_name"`
	RollNumber   int           `json:"roll_number"`
	ClassroomID  string        `json:"classroom_id"`
	AcademicYear string        `json:"academic_year"`
	Year         int           `json:"year"`
	Month        int           `json:"month"`

	ReportRecipients []ReportRecipient `json:"report_recipients"`

	Attendance *attendanceModel.MonthlyAttendanceSummary `json:"attendance"`
	Assessment *assessmentModel.MonthlyAssessmentSummary `json:"assessment"`

	Address studentModel.Address `json:"address"`
}
