package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Attendance struct {
	ID           bson.ObjectID `json:"id" bson:"_id,omitempty"`
	CollegeID    bson.ObjectID `json:"college_id" bson:"college_id"`
	StudentID    bson.ObjectID `json:"student_id" bson:"student_id"`
	ClassroomID  string        `json:"classroom_id" bson:"classroom_id"`
	AcademicYear string        `json:"academic_year" bson:"academic_year"`
	Date         string        `json:"date" bson:"date"`
	Attendance   string        `json:"attendance" bson:"attendance"`
}

// ========================================
// STUDENT FULL ACADEMIC YEAR SUMMARY
// ========================================

type AttendanceSummary struct {
	StudentID            bson.ObjectID `json:"student_id"`
	AcademicYear         string        `json:"academic_year"`
	TotalDays            int           `json:"total_days"`
	Present              int           `json:"present"`
	Absent               int           `json:"absent"`
	AttendancePercentage float64       `json:"attendance_percentage"`
}

// ========================================
// STUDENT MONTHLY SUMMARY
// ========================================

type MonthlyAttendanceSummary struct {
	StudentID            bson.ObjectID `json:"student_id"`
	AcademicYear         string        `json:"academic_year"`
	Year                 int           `json:"year"`
	Month                int           `json:"month"`
	WorkingDays          int           `json:"working_days"`
	AttendanceMarkedDays int           `json:"attendance_marked_days"`
	Present              int           `json:"present"`
	Absent               int           `json:"absent"`
	AttendancePercentage float64       `json:"attendance_percentage"`
}

// ========================================
// CLASSROOM STUDENT MONTHLY SUMMARY
// ========================================

type ClassroomStudentAttendanceSummary struct {
	StudentID            bson.ObjectID `json:"student_id"`
	RollNumber           int           `json:"roll_number"`
	Name                 string        `json:"name"`
	Present              int           `json:"present"`
	Absent               int           `json:"absent"`
	AttendanceMarkedDays int           `json:"attendance_marked_days"`
	AttendancePercentage float64       `json:"attendance_percentage"`
}

// ========================================
// CLASSROOM MONTHLY SUMMARY
// ========================================

type ClassroomMonthlyAttendanceSummary struct {
	ClassroomID          string                              `json:"classroom_id"`
	AcademicYear         string                              `json:"academic_year"`
	Year                 int                                 `json:"year"`
	Month                int                                 `json:"month"`
	WorkingDays          int                                 `json:"working_days"`
	AttendanceMarkedDays int                                 `json:"attendance_marked_days"`
	TotalStudents        int                                 `json:"total_students"`
	Present              int                                 `json:"present"`
	Absent               int                                 `json:"absent"`
	AttendancePercentage float64                             `json:"attendance_percentage"`
	Students             []ClassroomStudentAttendanceSummary `json:"students"`
}

// ========================================
// BULK ATTENDANCE
// ========================================

type BulkAttendanceItem struct {
	StudentID  bson.ObjectID `json:"student_id" bson:"student_id"`
	Attendance string        `json:"attendance" bson:"attendance"`
}

type BulkAttendanceRequest struct {
	ClassroomID  string               `json:"classroom_id"`
	AcademicYear string               `json:"academic_year"`
	Date         string               `json:"date"`
	Attendances  []BulkAttendanceItem `json:"attendances"`
}
