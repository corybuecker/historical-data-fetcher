package parsers

import "time"

type History struct {
	Date   TradierDate
	Open   float32
	High   float32
	Low    float32
	Close  float32
	Volume uint32
}

type Parser interface {
	Read(string) (History, error)
	fetch(url string) ([]byte, error)
}

func yesterday() string {
	return time.Now().AddDate(0, 0, -1).Format("2006-01-02")
}
