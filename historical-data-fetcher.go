package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/corybuecker/historical-data-fetcher/database"
	"github.com/corybuecker/historical-data-fetcher/parsers"
	"github.com/corybuecker/historical-data-fetcher/storage"
	"github.com/corybuecker/redisconfig"
	"github.com/davecgh/go-spew/spew"
	redis "gopkg.in/redis.v5"
)

type Config struct {
	TradierAPIKey string
	S3Id          string
	S3Secret      string
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
	var bucket = &storage.Bucket{}

	spew.Dump(bucket.CreateSession(config.S3Id, config.S3Secret))
	symbolsFetcher := database.Database{Client: redis}
	spew.Dump(symbolsFetcher.All())
	tradeFetcher := parsers.TradierParser{Token: config.TradierAPIKey}
	for _, symbol := range symbolsFetcher.Symbols {
		spew.Dump(symbol)
		day, err := tradeFetcher.Read(symbol.Symbol)
		spew.Dump(err)
		b, _ := json.Marshal(day)
		spew.Dump(bucket.Store(fmt.Sprintf("%s/%s/%s.json", symbol.Exchange, symbol.Symbol, day.Date.Time), string(b)))

	}
}
