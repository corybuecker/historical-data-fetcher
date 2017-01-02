package parsers

import (
	"fmt"
	"time"

	"github.com/corybuecker/jsonfetcher"
)

type CalendarResponse struct {
	Calendar Calendar
}

type Calendar struct {
	Month int
	Year  int
	Days  CalendarDays
}

type CalendarDays struct {
	Day []CalendarDay
}

type CalendarDay struct {
	Date   string
	Status string
}

func lastMonth() time.Time {
	return time.Now().UTC().Truncate(time.Hour*24).AddDate(0, -1, 0)
}

func getCombinedDays(token string) []CalendarDay {
	thisMonth, _ := getThisMonthCalendar(token)
	lastMonth, _ := getLastMonthCalendar(token)
	return append(lastMonth.Days.Day, thisMonth.Days.Day...)
}

func GetMostRecentOpenDay(token string) time.Time {
	var current time.Time
	for _, day := range getCombinedDays(token) {
		dayTime, _ := time.Parse("2006-01-02", day.Date)
		if day.Status == "open" && dayTime.UTC().Before(time.Now().UTC().Truncate(time.Hour*24)) {
			current = dayTime.UTC()
		}
	}
	return current
}

func getThisMonthCalendar(token string) (calendar *Calendar, err error) {
	calendarResponse := &CalendarResponse{}
	jsonFetcher := &jsonfetcher.Jsonfetcher{}
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token), "Accept": "application/json"}

	err = jsonFetcher.Get(fmt.Sprintf("https://sandbox.tradier.com/v1/markets/calendar?month=%d&year=%d", int(time.Now().UTC().Month()), time.Now().UTC().Year()), headers, calendarResponse)

	return &calendarResponse.Calendar, err
}

func getLastMonthCalendar(token string) (calendar *Calendar, err error) {
	calendarResponse := &CalendarResponse{}
	jsonFetcher := &jsonfetcher.Jsonfetcher{}
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token), "Accept": "application/json"}

	err = jsonFetcher.Get(fmt.Sprintf("https://sandbox.tradier.com/v1/markets/calendar?month=%d&year=%d", int(lastMonth().Month()), lastMonth().Year()), headers, calendarResponse)

	return &calendarResponse.Calendar, err
}
