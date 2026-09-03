package model

import (
	assessmentModel "github.com/shivkumar7dandin-gif/students-api/internal/assessment/model"
	attendanceModel "github.com/shivkumar7dandin-gif/students-api/internal/attendance/model"
)

type StudentDetailsResponse struct {
	Student           *Student                           `json:"student"`
	AttendanceSummary *attendanceModel.AttendanceSummary `json:"attendance_summary"`
	Assessments       []assessmentModel.Assessment       `json:"assessments"`
}
