package parsers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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

type TradierParser struct {
	Client    *http.Client
	Token     string
	RateLimit time.Duration
}

func (parser *TradierParser) obeyRateLimits(headers http.Header) error {
	var rateLimitAvailable int64
	var rateLimitExpires time.Duration
	var err error
	if ratelimitAvailableHeader, ok := headers["X-Ratelimit-Available"]; ok {
		var ratelimitAvailableTemp int64
		if ratelimitAvailableTemp, err = strconv.ParseInt(ratelimitAvailableHeader[0], 10, 8); err != nil {
			return err
		}
		rateLimitAvailable = ratelimitAvailableTemp
	}

	if ratelimitExpiresHeader, ok := headers["X-Ratelimit-Expiry"]; ok {
		var ratelimitExpiresTemp int64
		if ratelimitExpiresTemp, err = strconv.ParseInt(ratelimitExpiresHeader[0], 10, 64); err != nil {
			return err
		}
		rateLimitExpires = time.Unix(ratelimitExpiresTemp/1000, 0).Sub(time.Now())
	}
	parser.RateLimit = time.Duration(int64(rateLimitExpires)/rateLimitAvailable) + time.Millisecond*100
	time.Sleep(parser.RateLimit)
	return nil
}
func (parser *TradierParser) fetch(url string) ([]byte, error) {
	if parser.Client == nil {
		parser.Client = &http.Client{}
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
