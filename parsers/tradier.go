package parsers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	Client *http.Client
}

func (parser *TradierParser) Fetch(url string) ([]byte, error) {
	if parser.Client == nil {
		parser.Client = &http.Client{}
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer 8RroaVdaoWDGxthGzs0G6KJFgUUK")
	req.Header.Add("Accept", "application/json")
	resp, err := parser.Client.Do(req)
	if err != nil {
		return nil, err
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
	body, err := parser.Fetch(fmt.Sprintf("https://sandbox.tradier.com/v1/markets/timesales?symbol=%s&interval=1min", symbol))
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
