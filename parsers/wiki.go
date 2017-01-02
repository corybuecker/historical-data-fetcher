package parsers

import (
	"fmt"
	"time"

	"github.com/corybuecker/jsonfetcher"
)

type WikiResponse struct {
	Datatable struct {
		Data []interface{}
	}
}

type WikiParser struct {
	jsonFetcher *jsonfetcher.Jsonfetcher
	token       string
	headers     map[string]string
}

func BuildWikiParser(token string) *WikiParser {
	parser := &WikiParser{
		jsonFetcher: &jsonfetcher.Jsonfetcher{},
		token:       token,
		headers:     make(map[string]string),
	}
	return parser
}

func (parser *WikiParser) FetchLastMonth(symbol string) ([]History, error) {
	var test []History
	temp := WikiResponse{}

	parser.jsonFetcher.Get(fmt.Sprintf("https://www.quandl.com/api/v3/datatables/WIKI/PRICES.json?ticker=%s&date.gt=%s&api_key=%s", symbol, fourteenDaysAgo(), parser.token), parser.headers, &temp)
	for _, arr := range temp.Datatable.Data {
		time, _ := time.Parse("2006-01-02", arr.([]interface{})[1].(string))
		test = append(test, History{
			Date: TradierDate{
				Time: time,
			},
			Open:   float32(arr.([]interface{})[2].(float64)),
			High:   float32(arr.([]interface{})[3].(float64)),
			Low:    float32(arr.([]interface{})[4].(float64)),
			Close:  float32(arr.([]interface{})[5].(float64)),
			Volume: uint32(arr.([]interface{})[6].(float64)),
		})
	}
	return test, nil
}
