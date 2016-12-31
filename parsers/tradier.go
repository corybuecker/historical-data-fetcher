package parsers

import (
	"fmt"
	"strings"
	"time"

	"github.com/corybuecker/historicaldata/ratelimiters"
	"github.com/corybuecker/jsonfetcher"
)

type TradierData struct {
	History struct {
		Day History
	}
}

type TradierMultiData struct {
	History struct {
		Day []History
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
	rateLimiter *ratelimiters.TradierRateLimiter
	headers     map[string]string
	jsonFetcher *jsonfetcher.Jsonfetcher
}

func BuildTradierParser(token string) *TradierParser {
	parser := &TradierParser{
		rateLimiter: &ratelimiters.TradierRateLimiter{},
		headers:     map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token), "Accept": "application/json"},
		jsonFetcher: &jsonfetcher.Jsonfetcher{},
	}
	return parser
}

func (parser *TradierParser) FetchLastMonth(symbol string) ([]History, error) {
	var test []History
	temp := TradierMultiData{}

	err := parser.jsonFetcher.Get(fmt.Sprintf("https://sandbox.tradier.com/v1/markets/history?symbol=%s&start=%s&end=%s", symbol, threeDaysAgo(), yesterday()), parser.headers, &temp)

	if err != nil {
		return test, err
	}

	if err := parser.rateLimiter.ObeyRateLimit(parser.jsonFetcher.LastResponseHeaders()); err != nil {
		return test, err
	}

	return temp.History.Day, nil
}

func (parser *TradierParser) FetchYesterday(symbol string) (History, error) {
	var test History
	temp := TradierData{}

	err := parser.jsonFetcher.Get(fmt.Sprintf("https://sandbox.tradier.com/v1/markets/history?symbol=%s&start=%s&end=%s", symbol, yesterday(), yesterday()), parser.headers, &temp)

	if err != nil {
		return test, err
	}

	if err := parser.rateLimiter.ObeyRateLimit(parser.jsonFetcher.LastResponseHeaders()); err != nil {
		return test, err
	}

	return temp.History.Day, nil
}
