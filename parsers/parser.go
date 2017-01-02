package parsers

import (
	"encoding/json"
	"time"
)

type History struct {
	Date   TradierDate
	Open   float32
	High   float32
	Low    float32
	Close  float32
	Volume uint32
}

func (day *History) Serialize(symbol, exchange string) string {
	bytes, _ := json.Marshal(map[string]interface{}{
		"date":     day.Date.Time,
		"open":     day.Open,
		"high":     day.High,
		"low":      day.Low,
		"close":    day.Close,
		"volume":   day.Volume,
		"symbol":   symbol,
		"exchange": exchange,
	})
	return string(bytes)
}

type Parser interface {
	Read(string) (History, error)
	fetch(url string) ([]byte, error)
}

func yesterday() string {
	return time.Now().AddDate(0, 0, -1).Format("2006-01-02")
}

func fourteenDaysAgo() string {
	return time.Now().AddDate(0, 0, -14).Format("2006-01-02")
}

func thirtydaysago() string {
	return time.Now().AddDate(0, 0, -30).Format("2006-01-02")
}
