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
	Symbol   string
	Exchange string
}

func yesterday() string {
	return time.Now().UTC().AddDate(0, 0, -1).Truncate(time.Hour * 24).Format(time.RFC3339)
}

func (fetcher *Database) LoadSymbolsNeedingUpdate() error {
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

	// fetcher.Symbols = append(fetcher.Symbols, Symbol{Symbol: "AOSL", Exchange: "NASDAQ"})
	var tempSymbols []Symbol

	for _, symbol := range fetcher.Symbols {
		if yesterday() != fetcher.getLastDate(symbol) {
			tempSymbols = append(tempSymbols, symbol)
		}
	}

	fetcher.Symbols = tempSymbols

	return nil
}

func (fetcher *Database) getLastDate(symbol Symbol) string {
	var values map[string]string
	var err error

	values, err = fetcher.Client.HGetAll(fmt.Sprintf("%s:%s", symbol.Exchange, symbol.Symbol))

	if err != nil {
		return ""
	}

	return values["last_date_fetched"]
}

func (fetcher *Database) fetchExchange(exchange string) error {
	var symbolsFromRedis []string
	var err error

	if symbolsFromRedis, err = fetcher.Client.SMembers(exchange); err != nil {
		return err
	}

	for _, symbol := range symbolsFromRedis {
		fetcher.Symbols = append(fetcher.Symbols, Symbol{Symbol: symbol, Exchange: exchange})
	}

	return err
}

func (fetcher *Database) UpdateSymbolFetched(exchange, symbol string) error {
	log.Printf("updating last fetched date for %s:%s to %s", exchange, symbol, yesterday())
	return fetcher.Client.HSet(fmt.Sprintf("%s:%s", exchange, symbol), "last_date_fetched", yesterday())
}

func (fetcher *Database) IncrementStoreCount() error {
	return fetcher.Client.HIncrBy("metrics", "s3_write", 1)
}

func (fetcher *Database) SetLastUpdated(exchange, symbol string) error {
	log.Printf("updating last updated for %s:%s", exchange, symbol)
	return fetcher.Client.HSet(fmt.Sprintf("%s:%s", exchange, symbol), "last_updated", time.Now().UTC().Format(time.RFC3339))
}
