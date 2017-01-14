package database

import (
	"fmt"
	"time"
)

type Symbol struct {
	Symbol          string
	Exchange        string
	LastDateFetched time.Time
	client          DatabaseClient
}

func (symbol *Symbol) UpdateFetched(mostRecentOpenDay time.Time) error {
	return symbol.client.HSet(fmt.Sprintf("%s:%s", symbol.Exchange, symbol.Symbol), "last_date_fetched", mostRecentOpenDay.Format(time.RFC3339))
}

func (symbol *Symbol) IncrementStoreCount() error {
	return symbol.client.HIncrBy("metrics", "s3_write", 1)
}

func (symbol *Symbol) IncrementDateCount(date string) error {
	return symbol.client.HIncrBy("metrics", date, 1)
}

func (symbol *Symbol) SetLastUpdated() error {
	return symbol.client.HSet(fmt.Sprintf("%s:%s", symbol.Exchange, symbol.Symbol), "last_updated", time.Now().UTC().Format(time.RFC3339))
}

func (symbol *Symbol) MarkPresentInWiki() error {
	return symbol.client.HSet(fmt.Sprintf("%s:%s", symbol.Exchange, symbol.Symbol), "present_in_wiki", "true")
}
func (symbol *Symbol) getLastDate() bool {
	var values map[string]string

	values, _ = symbol.client.HGetAll(fmt.Sprintf("%s:%s", symbol.Exchange, symbol.Symbol))

	if lastDate, err := time.Parse(time.RFC3339, values["last_date_fetched"]); err != nil {
		symbol.LastDateFetched = time.Now().UTC().Truncate(time.Hour*24).AddDate(0, -1, 0)
	} else {
		symbol.LastDateFetched = lastDate
	}

	return true
}
