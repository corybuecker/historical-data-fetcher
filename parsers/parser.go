package parsers

import "time"

type Tick struct {
	SymbolID string
	Time     time.Time
	Volume   uint32
	Price    float32
}

type Parser interface {
	Read(string, string) ([]Tick, error)
	fetch(url string) ([]byte, error)
}
