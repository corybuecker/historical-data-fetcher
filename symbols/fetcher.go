package symbols

import redis "gopkg.in/redis.v5"

type Fetcher struct {
	Redis *redis.Client
}

func (fetcher *Fetcher) All() ([]string, error) {
	var nasdaq, amex, nyse, symbols []string
	var err error

	nasdaq, err = fetcher.Nasdaq()
	amex, err = fetcher.Nyse()
	nyse, err = fetcher.Amex()
	symbols = append(append(nasdaq, amex...), nyse...)

	return symbols, err
}

func (fetcher *Fetcher) Nasdaq() ([]string, error) {
	symbols, err := fetcher.Redis.SMembers("NASDAQ").Result()
	return symbols, err
}

func (fetcher *Fetcher) Nyse() ([]string, error) {
	symbols, err := fetcher.Redis.SMembers("NYSE").Result()
	return symbols, err
}

func (fetcher *Fetcher) Amex() ([]string, error) {
	symbols, err := fetcher.Redis.SMembers("AMEX").Result()
	return symbols, err
}
