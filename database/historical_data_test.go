package database

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var historicalData HistoricalData
var now time.Time

func init() {
	historicalData = HistoricalData{}
	now = time.Now()
}

func TestMostRecentDayWithoutAnyDays(t *testing.T) {
	assert.Equal(t, time.Now().UTC().Truncate(time.Hour*24).AddDate(0, -1, 0), historicalData.MostRecentDay())
}

func TestSort(t *testing.T) {
	historicalData = append(historicalData, HistoricalDatum{Date: now, Close: 1.0})
	historicalData = append(historicalData, HistoricalDatum{Date: now.AddDate(-1, 0, 0), Close: 2.0})

	sort.Sort(historicalData)

	assert.Equal(t, float32(2.0), historicalData[0].Close)
}

func TestMostRecentDayWithDays(t *testing.T) {
	assert.Equal(t, now, historicalData.MostRecentDay())
}
