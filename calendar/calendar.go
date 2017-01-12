package calendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/corybuecker/jsonfetcher"
)

type Calendar []Day

type Day struct {
	Date   time.Time
	Status string
}

type calendarResponse struct {
	Calendar struct {
		Days struct {
			Day []struct {
				Date   calendarDate
				Status string
			}
		}
	}
}

type calendarDate struct {
	time.Time
}

func (calendarDate *calendarDate) UnmarshalJSON(b []byte) (err error) {
	calendarDate.Time, err = time.Parse("2006-01-02", strings.Trim(string(b), "\""))
	return
}

type calendarFetcher struct {
	jsonFetcher jsonfetcher.Fetcher
	token       string
	headers     map[string]string
}

func (calendarFetcher *calendarFetcher) getCalendarForDate(month int, year int) (calendar Calendar, err error) {
	calendarResponse := new(calendarResponse)

	err = calendarFetcher.jsonFetcher.Get(fmt.Sprintf("https://sandbox.tradier.com/v1/markets/calendar?month=%d&year=%d", month, year), calendarFetcher.headers, calendarResponse)

	for _, day := range calendarResponse.Calendar.Days.Day {
		calendar = append(calendar, Day{Date: day.Date.Time, Status: day.Status})
	}

	return
}

func GetMostRecentOpenDay(token string) (mostRecentDay time.Time, err error) {
	var calendarFetcher = &calendarFetcher{
		jsonFetcher: &jsonfetcher.Jsonfetcher{},
		token:       token,
		headers:     map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token), "Accept": "application/json"},
	}
	var thisMonth, lastMonth Calendar

	thisMonth, err = calendarFetcher.getCalendarForDate(int(time.Now().UTC().Month()), int(time.Now().UTC().Year()))

	if err != nil {
		return
	}

	lastMonth, err = calendarFetcher.getCalendarForDate(int(lastMonthDate().Month()), int(lastMonthDate().Year()))

	if err != nil {
		return
	}

	for _, day := range append(lastMonth, thisMonth...) {
		if day.Status == "open" && day.Date.Before(time.Now().UTC().Truncate(time.Hour*24)) {
			mostRecentDay = day.Date.UTC()
		}
	}

	return
}

func lastMonthDate() time.Time {
	return time.Now().UTC().Truncate(time.Hour*24).AddDate(0, -1, 0)
}
