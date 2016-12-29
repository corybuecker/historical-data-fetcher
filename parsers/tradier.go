package parsers

import (
	"fmt"
	"strings"
	"time"

	"github.com/corybuecker/jsonfetcher"
	"github.com/corybuecker/trade-fetcher/ratelimiters"
)

type TradierData struct {
	History struct {
		Day History
	}
}

type TradierDate struct {
	Time time.Time
}

func (td *TradierDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")

	td.Time, err = time.Parse("2006-01-02", s)
	return
}

type clock struct{}

type TradierParser struct {
	Token       string
	RateLimiter *ratelimiters.TradierRateLimiter
}

func (parser *TradierParser) Read(symbol string) (History, error) {
	if parser.RateLimiter == nil {
		parser.RateLimiter = &ratelimiters.TradierRateLimiter{}
	}
	var test History
	temp := TradierData{}
	jsonfetcher := jsonfetcher.Jsonfetcher{}
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", parser.Token), "Accept": "application/json"}

	err := jsonfetcher.Get(fmt.Sprintf("https://sandbox.tradier.com/v1/markets/history?symbol=%s&start=%s&end=%s", symbol, yesterday(), yesterday()), headers, &temp)

	if err != nil {
		return test, err
	}

	if err := parser.RateLimiter.ObeyRateLimit(jsonfetcher.LastResponseHeaders()); err != nil {
		return test, err
	}

	test = temp.History.Day

	return test, nil
}
