package main

import (
	"fmt"
	"time"

	"github.com/corybuecker/historicaldata/database"
	"github.com/corybuecker/historicaldata/parsers"
	"github.com/corybuecker/historicaldata/storage"
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
	rConfig.Get(&config)

	spew.Dump(config)

	var bucket = &storage.Bucket{}

	bucket.CreateSession(config.S3Id, config.S3Secret)

	symbolsFetcher := database.Database{Client: &database.RedisClient{Client: redis}}
	symbolsFetcher.LoadSymbolsNeedingUpdate()

	tradeFetcher := parsers.BuildTradierParser(config.TradierAPIKey)

	for _, symbol := range symbolsFetcher.Symbols {
		days, _ := tradeFetcher.FetchLastMonth(symbol.Symbol)

		for _, day := range days {
			bucket.Store(fmt.Sprintf("%s/%s/%s.json", symbol.Exchange, symbol.Symbol, day.Date.Time.Format(time.RFC3339)), day.Serialize(symbol.Symbol, symbol.Exchange))
			symbolsFetcher.IncrementStoreCount()
		}

		symbolsFetcher.UpdateSymbolFetched(symbol.Exchange, symbol.Symbol)
		symbolsFetcher.SetLastUpdated(symbol.Exchange, symbol.Symbol)
	}
}
