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

func (s *CalendarService) GetDayStatus(
	ctx context.Context,
	collegeID bson.ObjectID,
	academicYear string,
	date time.Time,
) (*calendarModel.DayStatus, error) {

	// Get calendar for this college and academic year
	calendar, err := s.repo.GetByAcademicYear(
		ctx,
		collegeID,
		academicYear,
	)
	if err != nil {
		return nil, err
	}

	// Convert date to date-only
	checkDate := dateOnlyUTC(date)

	startDate := dateOnlyUTC(calendar.StartDate)
	endDate := dateOnlyUTC(calendar.EndDate)

	// ----------------------------------------
	// 1. Check academic year range
	// ----------------------------------------

	if checkDate.Before(startDate) ||
		checkDate.After(endDate) {

		return &calendarModel.DayStatus{
			Date:         checkDate.Format("2006-01-02"),
			IsWorkingDay: false,
			Reason:       "outside academic year",
		}, nil
	}

	// ----------------------------------------
	// 2. Check specific holidays
	// ----------------------------------------

	for _, holiday := range calendar.Holidays {

		if sameDate(checkDate, holiday.Date) {

			holidayCopy := holiday

			return &calendarModel.DayStatus{
				Date:         checkDate.Format("2006-01-02"),
				IsWorkingDay: false,
				Reason:       holiday.Name,
				Holiday:      &holidayCopy,
			}, nil
		}
	}

	// ----------------------------------------
	// 3. Check weekly holidays
	// ----------------------------------------

	weekday := checkDate.Weekday().String()

	for _, weeklyHoliday := range calendar.WeeklyHolidays {

		if strings.EqualFold(
			weeklyHoliday,
			weekday,
		) {

			return &calendarModel.DayStatus{
				Date:         checkDate.Format("2006-01-02"),
				IsWorkingDay: false,
				Reason:       "weekly holiday",
			}, nil
		}
	}

	// ----------------------------------------
	// 4. Working day
	// ----------------------------------------

	return &calendarModel.DayStatus{
		Date:         checkDate.Format("2006-01-02"),
		IsWorkingDay: true,
		Reason:       "working day",
	}, nil
}

func sameDate(
	a time.Time,
	b time.Time,
) bool {

	ay, am, ad := a.Date()
	by, bm, bd := b.Date()

	return ay == by &&
		am == bm &&
		ad == bd
}

func dateOnlyUTC(
	t time.Time,
) time.Time {

	year, month, day := t.Date()

	return time.Date(
		year,
		month,
		day,
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

func (s *CalendarService) GetMonthlyWorkingDays(
	ctx context.Context,
	collegeID bson.ObjectID,
	academicYear string,
	year int,
	month int,
) (int, error) {

	if month < 1 || month > 12 {
		return 0, errors.New(
			"month must be between 1 and 12",
		)
	}

	calendar, err := s.repo.GetByAcademicYear(
		ctx,
		collegeID,
		academicYear,
	)
	if err != nil {
		return 0, err
	}

	// First day of requested month.
	currentDate := time.Date(
		year,
		time.Month(month),
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	// First day of next month.
	nextMonth := currentDate.AddDate(
		0,
		1,
		0,
	)

	workingDays := 0

	for currentDate.Before(nextMonth) {

		// Skip dates outside academic year.
		if currentDate.Before(
			dateOnlyUTC(calendar.StartDate),
		) ||
			currentDate.After(
				dateOnlyUTC(calendar.EndDate),
			) {

			currentDate =
				currentDate.AddDate(0, 0, 1)

			continue
		}

		isHoliday := false

		// Check explicit holidays.
		for _, holiday := range calendar.Holidays {

			if sameDate(
				currentDate,
				holiday.Date,
			) {

				isHoliday = true
				break
			}
		}

		if isHoliday {

			currentDate =
				currentDate.AddDate(0, 0, 1)

			continue
		}

		// Check weekly holidays.
		weekday :=
			currentDate.Weekday().String()

		for _, weeklyHoliday := range calendar.WeeklyHolidays {

			if strings.EqualFold(
				strings.TrimSpace(
					weeklyHoliday,
				),
				weekday,
			) {

				isHoliday = true
				break
			}
		}

		if !isHoliday {
			workingDays++
		}

		currentDate =
			currentDate.AddDate(0, 0, 1)
	}

	return workingDays, nil
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
