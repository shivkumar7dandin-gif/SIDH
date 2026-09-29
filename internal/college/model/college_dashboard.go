package model

import (
	classroomModel "github.com/shivkumar7dandin-gif/students-api/internal/classroom/model"
	studentModel "github.com/shivkumar7dandin-gif/students-api/internal/student/model"
	teacherModel "github.com/shivkumar7dandin-gif/students-api/internal/teacher/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Today's overall student attendance
type TodayAttendance struct {
	TotalStudents int     `json:"total_students"`
	Present       int     `json:"present"`
	Absent        int     `json:"absent"`
	Marked        int     `json:"marked"`
	NotMarked     int     `json:"not_marked"`
	Percentage    float64 `json:"percentage"`

	PresentStudents   []studentModel.Student `json:"present_students"`
	AbsentStudents    []studentModel.Student `json:"absent_students"`
	NotMarkedStudents []studentModel.Student `json:"not_marked_students"`
}

// Today's attendance for each classroom
type ClassroomAttendance struct {
	ClassroomID   string  `json:"classroom_id"`
	ClassroomName string  `json:"classroom_name"`
	Section       string  `json:"section"`
	TotalStudents int     `json:"total_students"`
	Present       int     `json:"present"`
	Absent        int     `json:"absent"`
	Marked        int     `json:"marked"`
	NotMarked     int     `json:"not_marked"`
	Percentage    float64 `json:"percentage"`
	Status        string  `json:"status"`
}

type TeacherStatusSummary struct {
	Total  int `json:"total"`
	Active int `json:"active"`
	Left   int `json:"left"`

	ActiveTeachers []teacherModel.Teacher `json:"active_teachers"`
	LeftTeachers   []teacherModel.Teacher `json:"left_teachers"`
}

type TeacherTodayClass struct {
	TeacherID    bson.ObjectID `json:"teacher_id"`
	TeacherName  string        `json:"teacher_name"`
	Subject      string        `json:"subject"`
	TotalClasses int           `json:"total_classes"`
}

type CollegeDashboard struct {
	College    *College                   `json:"college"`
	Teachers   []teacherModel.Teacher     `json:"teachers"`
	Classrooms []classroomModel.Classroom `json:"classrooms"`
	Students   []studentModel.Student     `json:"students"`

	TotalTeachers   int `json:"total_teachers"`
	ActiveTeachers  int `json:"active_teachers"`
	LeftTeachers    int `json:"left_teachers"`
	TotalClassrooms int `json:"total_classrooms"`
	TotalStudents   int `json:"total_students"`

	TodayAttendance     TodayAttendance       `json:"today_attendance"`
	TeacherStatus       TeacherStatusSummary  `json:"teacher_status"`
	ClassroomAttendance []ClassroomAttendance `json:"classroom_attendance"`

	TeacherTodayClasses []TeacherTodayClass `json:"teacher_today_classes"`
}
