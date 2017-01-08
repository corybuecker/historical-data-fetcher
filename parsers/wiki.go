package parsers

import (
	"errors"
	"fmt"
	"time"

	"github.com/corybuecker/historicaldata/database"
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

func (parser *WikiParser) FetchIntoSlice(symbol *database.Symbol) (database.HistoricalData, error) {
	temp := WikiResponse{}
	slice := make(database.HistoricalData, 0)

	parser.jsonFetcher.Get(fmt.Sprintf("https://www.quandl.com/api/v3/datatables/WIKI/PRICES.json?ticker=%s&date.gt=%s&api_key=%s", symbol.Symbol, symbol.LastDateFetched.Format(time.RFC3339), parser.token), parser.headers, &temp)

	if len(temp.Datatable.Data) == 0 {
		return nil, errors.New("there was no data in wiki")
	}

	for _, arr := range temp.Datatable.Data {
		time, err := time.Parse("2006-01-02", arr.([]interface{})[1].(string))
		if err != nil {
			return nil, err
		}
		slice = append(slice, database.HistoricalDatum{
			Date:     time,
			Open:     float32(arr.([]interface{})[2].(float64)),
			High:     float32(arr.([]interface{})[3].(float64)),
			Low:      float32(arr.([]interface{})[4].(float64)),
			Close:    float32(arr.([]interface{})[5].(float64)),
			Volume:   uint32(arr.([]interface{})[6].(float64)),
			Symbol:   symbol.Symbol,
			Exchange: symbol.Exchange,
		})
	}

	return slice, nil
}
