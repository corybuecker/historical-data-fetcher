package database

import (
	"fmt"
	"time"
)

type Symbol struct {
	Symbol          string
	Exchange        string
	LastDateFetched *time.Time
	client          DatabaseClient
}

func (symbol *Symbol) UpdateFetched(mostRecentOpenDay time.Time) error {
	return symbol.client.HSet(fmt.Sprintf("%s:%s", symbol.Exchange, symbol.Symbol), "last_date_fetched", mostRecentOpenDay.Format(time.RFC3339))
}

func (symbol *Symbol) SetLastUpdated() error {
	return symbol.client.HSet(fmt.Sprintf("%s:%s", symbol.Exchange, symbol.Symbol), "last_updated", time.Now().UTC().Format(time.RFC3339))
}

func (symbol *Symbol) MarkPresentInWiki() error {
	return symbol.client.HSet(fmt.Sprintf("%s:%s", symbol.Exchange, symbol.Symbol), "present_in_wiki", "true")
}
func (symbol *Symbol) getLastDate() bool {
	var values map[string]string
	var err error

	symbol.LastDateFetched = nil

	values, err = symbol.client.HGetAll(fmt.Sprintf("%s:%s", symbol.Exchange, symbol.Symbol))

	if err != nil {
		return false
	}

	if lastDate, err := time.Parse(time.RFC3339, values["last_date_fetched"]); err != nil {
		return false
	} else {
		symbol.LastDateFetched = &lastDate
	}

	return true
}
