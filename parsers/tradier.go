package parsers

import (
	"fmt"
	"strings"
	"time"

	"github.com/corybuecker/historicaldata/database"
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

func (parser *TradierParser) FetchIntoSlice(symbol *database.Symbol) (database.HistoricalData, error) {
	slice := make(database.HistoricalData, 0)
	temp := TradierMultiData{}

	err := parser.jsonFetcher.Get(fmt.Sprintf("https://sandbox.tradier.com/v1/markets/history?symbol=%s&start=%s&end=%s", symbol.Symbol, symbol.LastDateFetched.AddDate(0, 0, -7).Format(time.RFC3339), yesterday()), parser.headers, &temp)
	if err := parser.rateLimiter.ObeyRateLimit(parser.jsonFetcher.LastResponseHeaders()); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	for _, day := range temp.History.Day {

		slice = append(slice, database.HistoricalDatum{
			Date:     day.Date.Time,
			Open:     day.Open,
			High:     day.High,
			Low:      day.Low,
			Close:    day.Close,
			Volume:   day.Volume,
			Symbol:   symbol.Symbol,
			Exchange: symbol.Exchange,
		})
	}

	return slice, nil
}
