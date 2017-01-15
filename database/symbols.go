package database

import "time"

type Symbols struct {
	Symbols []Symbol
	Client  DatabaseClient
}

func (symbols *Symbols) Initialize() (err error) {
	symbols.Symbols = make([]Symbol, 0)

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

	return
}

func (symbols *Symbols) Filter(mostRecentOpenDay time.Time) {
	var tempSymbols []Symbol

	for _, symbol := range symbols.Symbols {
		if !symbol.LastDateFetched.IsZero() && mostRecentOpenDay.After(symbol.LastDateFetched) {
			tempSymbols = append(tempSymbols, symbol)
		}
	}

	symbols.Symbols = tempSymbols

	return
}

func (symbols *Symbols) fetchExchange(exchange string) (err error) {
	var symbolsFromRedis []string

	if symbolsFromRedis, err = symbols.Client.SMembers(exchange); err != nil {
		return
	}

	for _, symbol := range symbolsFromRedis {
		symbol := Symbol{Symbol: symbol, Exchange: exchange, client: symbols.Client}
		symbol.getLastDate()
		symbols.Symbols = append(symbols.Symbols, symbol)
	}

	return
}
