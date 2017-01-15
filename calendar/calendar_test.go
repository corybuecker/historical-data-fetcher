package calendar

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var response = `{"calendar":{"month":1,"year":2017,"days":{"day":[{"date":"%s","status":"open","description":"Market is open","premarket":{"start":"08:00","end":"09:30"},"open":{"start":"09:30","end":"16:00"},"postmarket":{"start":"16:00","end":"20:00"}}]}}}`

func TestRunner(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("test bad calendar date for this month", testBadCalendarDateThisMonth)
	httpmock.Reset()
	t.Run("test bad calendar date for last month", testBadCalendarDateLastMonth)
	httpmock.Reset()
	t.Run("test valid calendar date", testValidCalendarDate)

}

func testBadCalendarDateThisMonth(t *testing.T) {
	httpmock.RegisterResponder("GET",
		fmt.Sprintf("https://sandbox.tradier.com/v1/markets/calendar?month=%d&year=%d", lastMonthDate().Month(), int(lastMonthDate().Year())),
		httpmock.NewStringResponder(200, fmt.Sprintf(response, "2016-01-01")),
	)

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("https://sandbox.tradier.com/v1/markets/calendar?month=%d&year=%d", int(time.Now().UTC().Month()), int(time.Now().UTC().Year())),
		httpmock.NewStringResponder(200, fmt.Sprintf(response, "junk")),
	)

	mostRecentOpenDay, err := GetMostRecentOpenDay("token")

	assert.True(t, mostRecentOpenDay.IsZero())
	assert.NotNil(t, err)
}

func testBadCalendarDateLastMonth(t *testing.T) {
	httpmock.RegisterResponder("GET",
		fmt.Sprintf("https://sandbox.tradier.com/v1/markets/calendar?month=%d&year=%d", lastMonthDate().Month(), int(lastMonthDate().Year())),
		httpmock.NewStringResponder(200, fmt.Sprintf(response, "junk")),
	)

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("https://sandbox.tradier.com/v1/markets/calendar?month=%d&year=%d", int(time.Now().UTC().Month()), int(time.Now().UTC().Year())),
		httpmock.NewStringResponder(200, fmt.Sprintf(response, "2016-01-01")),
	)

	mostRecentOpenDay, err := GetMostRecentOpenDay("token")

	assert.True(t, mostRecentOpenDay.IsZero())
	assert.NotNil(t, err)
}

func testValidCalendarDate(t *testing.T) {
	httpmock.RegisterResponder("GET",
		fmt.Sprintf("https://sandbox.tradier.com/v1/markets/calendar?month=%d&year=%d", lastMonthDate().Month(), int(lastMonthDate().Year())),
		httpmock.NewStringResponder(200, fmt.Sprintf(response, "2015-01-01")),
	)

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("https://sandbox.tradier.com/v1/markets/calendar?month=%d&year=%d", int(time.Now().UTC().Month()), int(time.Now().UTC().Year())),
		httpmock.NewStringResponder(200, fmt.Sprintf(response, "2016-01-01")),
	)

	mostRecentOpenDay, err := GetMostRecentOpenDay("token")

	assert.Equal(t, mostRecentOpenDay, time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.Nil(t, err)
}

//
// type Calendar []Day
//
// type Day struct {
// 	Date   time.Time
// 	Status string
// }
//
// type calendarResponse struct {
// 	Calendar struct {
// 		Days struct {
// 			Day []struct {
// 				Date   calendarDate
// 				Status string
// 			}
// 		}
// 	}
// }
//
// type calendarDate struct {
// 	time.Time
// }
//
// func (calendarDate *calendarDate) UnmarshalJSON(b []byte) (err error) {
// 	calendarDate.Time, err = time.Parse("2006-01-02", strings.Trim(string(b), "\""))
// 	return
// }
//
// type calendarFetcher struct {
// 	jsonFetcher jsonfetcher.Fetcher
// 	token       string
// 	headers     map[string]string
// }
//
// func (calendarFetcher *calendarFetcher) getCalendarForDate(month int, year int) (calendar Calendar, err error) {
// 	calendarResponse := new(calendarResponse)
//
// 	err = calendarFetcher.jsonFetcher.Get(fmt.Sprintf("https://sandbox.tradier.com/v1/markets/calendar?month=%d&year=%d", month, year), calendarFetcher.headers, calendarResponse)
//
// 	for _, day := range calendarResponse.Calendar.Days.Day {
// 		calendar = append(calendar, Day{Date: day.Date.Time, Status: day.Status})
// 	}
//
// 	return
// }
//
// func GetMostRecentOpenDay(token string) (mostRecentDay time.Time, err error) {
// 	var calendarFetcher = &calendarFetcher{
// 		jsonFetcher: &jsonfetcher.Jsonfetcher{},
// 		token:       token,
// 		headers:     map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token), "Accept": "application/json"},
// 	}
// 	var thisMonth, lastMonth Calendar
//
// 	thisMonth, err = calendarFetcher.getCalendarForDate(int(time.Now().UTC().Month()), int(time.Now().UTC().Year()))
//
// 	if err != nil {
// 		return
// 	}
//
// 	lastMonth, err = calendarFetcher.getCalendarForDate(int(lastMonthDate().Month()), int(lastMonthDate().Year()))
//
// 	if err != nil {
// 		return
// 	}
//
// 	for _, day := range append(lastMonth, thisMonth...) {
// 		if day.Status == "open" && day.Date.Before(time.Now().UTC().Truncate(time.Hour*24)) {
// 			mostRecentDay = day.Date.UTC()
// 		}
// 	}
//
// 	return
// }
//
// func lastMonthDate() time.Time {
// 	return time.Now().UTC().Truncate(time.Hour*24).AddDate(0, -1, 0)
// }
