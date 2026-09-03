package model

import (
	classroomModel "github.com/shivkumar7dandin-gif/students-api/internal/classroom/model"
	studentModel "github.com/shivkumar7dandin-gif/students-api/internal/student/model"
	teacherModel "github.com/shivkumar7dandin-gif/students-api/internal/teacher/model"
)

type CollegeDashboard struct {
	College    *College                   `json:"college"`
	Teachers   []teacherModel.Teacher     `json:"teachers"`
	Classrooms []classroomModel.Classroom `json:"classrooms"`
	Students   []studentModel.Student     `json:"students"`

	TotalTeachers   int `json:"total_teachers"`
	TotalClassrooms int `json:"total_classrooms"`
	TotalStudents   int `json:"total_students"`
}
