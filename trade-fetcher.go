package main

import (
	"log"

	"github.com/corybuecker/redisconfig"
	"github.com/corybuecker/trade-fetcher/parsers"
	"github.com/corybuecker/trade-fetcher/symbols"
	"github.com/davecgh/go-spew/spew"
	redis "gopkg.in/redis.v5"
)

type Config struct {
	TradierAPIKey string
}

func main() {
	redis := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "master",
		SentinelAddrs: []string{":26379"},
	})

	rConfig := redisconfig.Configuration{
		Client: redis,
		Key:    "trade_fetcher_configuration",
	}
	config := Config{}
	err := rConfig.Get(&config)
	log.Println(err)
	log.Printf("%v", config)
	symbolsFetcher := symbols.Fetcher{Redis: redis}
	symbols, _ := symbolsFetcher.All()
	tradeFetcher := parsers.TradierParser{Token: config.TradierAPIKey}

	for _, symbol := range symbols {
		spew.Dump(tradeFetcher.Read(symbol))

	}
}
