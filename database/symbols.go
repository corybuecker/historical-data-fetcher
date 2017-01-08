package database

import (
	"fmt"
	"log"
	"time"
)

type Database struct {
	Client  DatabaseClient
	Symbols []Symbol
}

type Symbol struct {
	Symbol          string
	Exchange        string
	LastDateFetched time.Time
}

func yesterday() string {
	return time.Now().UTC().AddDate(0, 0, -1).Truncate(time.Hour * 24).Format(time.RFC3339)
}

func thirtydaysago() time.Time {
	return time.Now().UTC().Truncate(time.Hour*24).AddDate(0, -1, 0)
}

func (fetcher *Database) LoadSymbolsNeedingUpdate(mostRecentOpenDay time.Time) error {
	fetcher.Symbols = make([]Symbol, 0)

	var err error

	if err = fetcher.fetchExchange("NASDAQ"); err != nil {
		return err
	}

	if err = fetcher.fetchExchange("AMEX"); err != nil {
		return err
	}

	if err = fetcher.fetchExchange("NYSE"); err != nil {
		return err
	}

	if err = fetcher.fetchExchange("NYSEARCA"); err != nil {
		return err
	}
	var tempSymbols []Symbol

	for _, symbol := range fetcher.Symbols {
		if !mostRecentOpenDay.Equal(symbol.LastDateFetched) {
			tempSymbols = append(tempSymbols, symbol)
		}
	}

	fetcher.Symbols = tempSymbols

	return nil
}

func (symbol *Symbol) getLastDate(fetcher *Database) bool {
	var values map[string]string

	values, _ = fetcher.Client.HGetAll(fmt.Sprintf("%s:%s", symbol.Exchange, symbol.Symbol))

	if lastDate, err := time.Parse(time.RFC3339, values["last_date_fetched"]); err != nil {
		symbol.LastDateFetched = thirtydaysago()
	} else {
		symbol.LastDateFetched = lastDate
	}

	return true
}

func (fetcher *Database) fetchExchange(exchange string) error {
	var symbolsFromRedis []string
	var err error

	if symbolsFromRedis, err = fetcher.Client.SMembers(exchange); err != nil {
		return err
	}

	for _, symbol := range symbolsFromRedis {
		symbol := Symbol{Symbol: symbol, Exchange: exchange}
		symbol.getLastDate(fetcher)
		fetcher.Symbols = append(fetcher.Symbols, symbol)
	}

	return err
}

func (fetcher *Database) UpdateSymbolFetched(exchange, symbol string, mostRecentOpenDay time.Time) error {
	log.Printf("updating last fetched date for %s:%s to %s", exchange, symbol, mostRecentOpenDay.Format(time.RFC3339))
	return fetcher.Client.HSet(fmt.Sprintf("%s:%s", exchange, symbol), "last_date_fetched", mostRecentOpenDay.Format(time.RFC3339))
}

func (fetcher *Database) IncrementStoreCount() error {
	return fetcher.Client.HIncrBy("metrics", "s3_write", 1)
}

func (fetcher *Database) IncrementDateCount(date string) error {
	return fetcher.Client.HIncrBy("metrics", date, 1)
}

func (fetcher *Database) SetLastUpdated(exchange, symbol string) error {
	log.Printf("updating last updated for %s:%s", exchange, symbol)
	return fetcher.Client.HSet(fmt.Sprintf("%s:%s", exchange, symbol), "last_updated", time.Now().UTC().Format(time.RFC3339))
}

func (fetcher *Database) MarkPresentInWiki(exchange, symbol string) error {
	log.Printf("updating present in wiki for %s:%s", exchange, symbol)
	return fetcher.Client.HSet(fmt.Sprintf("%s:%s", exchange, symbol), "present_in_wiki", "true")
}
