package database

import "time"

type Symbols struct {
	Symbols []Symbol
	client  DatabaseClient
}

func (symbols *Symbols) Initialize(mostRecentOpenDay time.Time) (err error) {
	symbols.Symbols = make([]Symbol, 0)

	var tempSymbols []Symbol

	if err = symbols.fetchExchange("NASDAQ"); err != nil {
		return
	}

	if err = symbols.fetchExchange("AMEX"); err != nil {
		return
	}

	if err = symbols.fetchExchange("NYSE"); err != nil {
		return
	}

	if err = symbols.fetchExchange("NYSEARCA"); err != nil {
		return
	}

	for _, symbol := range symbols.Symbols {
		if mostRecentOpenDay.After(symbol.LastDateFetched) {
			tempSymbols = append(tempSymbols, symbol)
		}
	}

	symbols.Symbols = tempSymbols

	return
}

func (symbols *Symbols) fetchExchange(exchange string) error {
	var symbolsFromRedis []string
	var err error

	if symbolsFromRedis, err = symbols.client.SMembers(exchange); err != nil {
		return err
	}

	for _, symbol := range symbolsFromRedis {
		symbol := Symbol{Symbol: symbol, Exchange: exchange, client: symbols.client}
		symbol.getLastDate()
		symbols.Symbols = append(symbols.Symbols, symbol)
	}

	return err
}
