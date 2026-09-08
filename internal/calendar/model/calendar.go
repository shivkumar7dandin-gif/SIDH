package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Holiday struct {
	Date time.Time `json:"date" bson:"date"`
	Name string    `json:"name" bson:"name"`
	Type string    `json:"type" bson:"type"`
}

type AcademicCalendar struct {
	ID bson.ObjectID `json:"id" bson:"_id,omitempty"`

	CollegeID bson.ObjectID `json:"college_id" bson:"college_id"`

	AcademicYear string `json:"academic_year" bson:"academic_year"`

	StartDate time.Time `json:"start_date" bson:"start_date"`
	EndDate   time.Time `json:"end_date" bson:"end_date"`

	WeeklyHolidays []string `json:"weekly_holidays" bson:"weekly_holidays"`

	Holidays []Holiday `json:"holidays" bson:"holidays"`
}

type CreateAcademicCalendarRequest struct {
	AcademicYear string `json:"academic_year"`

	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`

	WeeklyHolidays []string `json:"weekly_holidays"`

	Holidays []CreateHolidayRequest `json:"holidays"`
}

type CreateHolidayRequest struct {
	Date string `json:"date"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type DayStatus struct {
	Date         string   `json:"date"`
	IsWorkingDay bool     `json:"is_working_day"`
	Reason       string   `json:"reason"`
	Holiday      *Holiday `json:"holiday,omitempty"`
}
