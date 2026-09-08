package service

import (
	"context"
	"errors"
	"strings"
	"time"

	calendarModel "github.com/shivkumar7dandin-gif/students-api/internal/calendar/model"
	calendarRepository "github.com/shivkumar7dandin-gif/students-api/internal/calendar/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CalendarService struct {
	repo *calendarRepository.CalendarRepository
}

func NewCalendarService(
	repo *calendarRepository.CalendarRepository,
) *CalendarService {

	return &CalendarService{
		repo: repo,
	}
}

func (s *CalendarService) Create(
	ctx context.Context,
	req calendarModel.CreateAcademicCalendarRequest,
	collegeID bson.ObjectID,
) (*calendarModel.AcademicCalendar, error) {

	req.AcademicYear = strings.TrimSpace(req.AcademicYear)
	req.StartDate = strings.TrimSpace(req.StartDate)
	req.EndDate = strings.TrimSpace(req.EndDate)

	if req.AcademicYear == "" {
		return nil, errors.New("academic year is required")
	}

	if req.StartDate == "" {
		return nil, errors.New("start date is required")
	}

	if req.EndDate == "" {
		return nil, errors.New("end date is required")
	}

	startDate, err := time.Parse(
		"2006-01-02",
		req.StartDate,
	)

	if err != nil {
		return nil, errors.New(
			"invalid start date, use YYYY-MM-DD",
		)
	}

	endDate, err := time.Parse(
		"2006-01-02",
		req.EndDate,
	)

	if err != nil {
		return nil, errors.New(
			"invalid end date, use YYYY-MM-DD",
		)
	}

	if !endDate.After(startDate) {
		return nil, errors.New(
			"end date must be after start date",
		)
	}

	exists, err := s.repo.ExistsByAcademicYear(
		ctx,
		collegeID,
		req.AcademicYear,
	)

	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New(
			"academic calendar already exists for this year",
		)
	}

	holidays := make([]calendarModel.Holiday, 0)

	for _, holidayReq := range req.Holidays {

		holidayReq.Name = strings.TrimSpace(
			holidayReq.Name,
		)

		holidayReq.Type = strings.ToLower(
			strings.TrimSpace(
				holidayReq.Type,
			),
		)

		if holidayReq.Name == "" {
			return nil, errors.New(
				"holiday name is required",
			)
		}

		if holidayReq.Type == "" {
			return nil, errors.New(
				"holiday type is required",
			)
		}

		holidayDate, err := time.Parse(
			"2006-01-02",
			holidayReq.Date,
		)

		if err != nil {
			return nil, errors.New(
				"invalid holiday date, use YYYY-MM-DD",
			)
		}

		if holidayDate.Before(startDate) ||
			holidayDate.After(endDate) {

			return nil, errors.New(
				"holiday date must be inside academic year",
			)
		}

		holidays = append(
			holidays,
			calendarModel.Holiday{
				Date: holidayDate,
				Name: holidayReq.Name,
				Type: holidayReq.Type,
			},
		)
	}

	weeklyHolidays := make([]string, 0)

	for _, day := range req.WeeklyHolidays {

		day = strings.TrimSpace(day)

		if day == "" {
			continue
		}

		if !isValidWeekday(day) {
			return nil, errors.New(
				"invalid weekly holiday: " + day,
			)
		}

		weeklyHolidays = append(
			weeklyHolidays,
			day,
		)
	}

	calendar := calendarModel.AcademicCalendar{
		CollegeID:      collegeID,
		AcademicYear:   req.AcademicYear,
		StartDate:      startDate,
		EndDate:        endDate,
		WeeklyHolidays: weeklyHolidays,
		Holidays:       holidays,
	}

	return s.repo.Create(
		ctx,
		calendar,
	)
}

func (s *CalendarService) GetByCollegeID(
	ctx context.Context,
	collegeID bson.ObjectID,
) ([]calendarModel.AcademicCalendar, error) {

	return s.repo.GetByCollegeID(
		ctx,
		collegeID,
	)
}

func (s *CalendarService) GetByAcademicYear(
	ctx context.Context,
	collegeID bson.ObjectID,
	academicYear string,
) (*calendarModel.AcademicCalendar, error) {

	academicYear = strings.TrimSpace(
		academicYear,
	)

	if academicYear == "" {
		return nil, errors.New(
			"academic year is required",
		)
	}

	return s.repo.GetByAcademicYear(
		ctx,
		collegeID,
		academicYear,
	)
}

func isValidWeekday(day string) bool {

	validDays := map[string]bool{
		"Monday":    true,
		"Tuesday":   true,
		"Wednesday": true,
		"Thursday":  true,
		"Friday":    true,
		"Saturday":  true,
		"Sunday":    true,
	}

	return validDays[day]
}
