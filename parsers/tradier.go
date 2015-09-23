package parsers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/corybuecker/trade-fetcher/ratelimiters"
)

type TradierDatum struct {
	Time   string
	Price  float32
	Volume uint32
}
type TradierData struct {
	Series struct {
		Data []TradierDatum
	}
}

type clock struct{}

type TradierParser struct {
	Client      *http.Client
	Token       string
	RateLimiter *ratelimiters.TradierRateLimiter
}

func (parser *TradierParser) fetch(url string) ([]byte, error) {
	if parser.Client == nil {
		parser.Client = &http.Client{}
	}

	if parser.RateLimiter == nil {
		parser.RateLimiter = &ratelimiters.TradierRateLimiter{}
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", parser.Token))
	req.Header.Add("Accept", "application/json")
	resp, err := parser.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the call to the API failed with status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	if err := parser.RateLimiter.ObeyRateLimit(resp.Header); err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	return body, nil
}

func (parser *TradierParser) Read(symbol string, symbolID string) ([]Tick, error) {
	var test []Tick
	temp := TradierData{}
	body, err := parser.fetch(fmt.Sprintf("https://sandbox.tradier.com/v1/markets/timesales?symbol=%s&interval=1min", symbol))

	if err != nil {
		return nil, err
	}
	json.Unmarshal(body, &temp)
	location, _ := time.LoadLocation("America/New_York")
	for _, datum := range temp.Series.Data {
		time, _ := time.ParseInLocation("2006-01-02T15:04:05", datum.Time, location)
		test = append(test, Tick{Time: time, Volume: datum.Volume, Price: datum.Price, SymbolID: symbolID})
	}
	return test, nil
}
