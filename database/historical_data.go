package database

import (
	"sort"
	"time"
)

type HistoricalData []HistoricalDatum

type HistoricalDatum struct {
	Date     time.Time
	Open     float32
	High     float32
	Low      float32
	Close    float32
	Volume   uint32
	Symbol   string
	Exchange string
}

func (slice HistoricalData) Len() int {
	return len(slice)
}

func (slice HistoricalData) Less(i, j int) bool {
	return slice[i].Date.Before(slice[j].Date)
}

func (slice HistoricalData) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (historicalData HistoricalData) MostRecentDay() time.Time {
	sort.Sort(historicalData)

	if len(historicalData) > 0 {
		return historicalData[len(historicalData)-1].Date
	} else {
		return time.Now().UTC().Truncate(time.Hour*24).AddDate(0, -1, 0)
	}
}
