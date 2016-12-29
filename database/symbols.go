package database

import redis "gopkg.in/redis.v5"

type Database struct {
	Client  *redis.Client
	Symbols []Symbol
}

type Symbol struct {
	Symbol   string
	Exchange string
}

func (db *Database) All() error {
	db.Symbols = make([]Symbol, 0)
	var err error

	err = db.fetchExchange("NASDAQ")
	err = db.fetchExchange("AMEX")
	err = db.fetchExchange("NYSE")

	return err
}

func (db *Database) fetchExchange(exchange string) error {
	var symbolsFromRedis []string
	var err error

	if symbolsFromRedis, err = db.Client.SMembers(exchange).Result(); err != nil {
		return err
	}

	for _, symbol := range symbolsFromRedis {
		db.Symbols = append(db.Symbols, Symbol{Symbol: symbol, Exchange: exchange})
	}

	return err
}
