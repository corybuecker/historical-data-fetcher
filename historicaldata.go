package main

import (
	"fmt"
	"log"
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

	var bucket = storage.CreateBucket(config.S3Id, config.S3Secret)

	mostRecentOpenDay := parsers.GetMostRecentOpenDay(config.TradierAPIKey)

	symbolsFetcher := database.Database{Client: &database.RedisClient{Client: redis}}
	symbolsFetcher.LoadSymbolsNeedingUpdate(mostRecentOpenDay)

	wikiFetcher := parsers.BuildWikiParser("y5cY91y4m99oEwPjfKBK")
	tradeFetcher := parsers.BuildTradierParser(config.TradierAPIKey)

	for _, symbol := range symbolsFetcher.Symbols {
		days, err := wikiFetcher.FetchLastMonth(symbol.Symbol)

		if days == nil {
			days, err = tradeFetcher.FetchLastMonth(symbol.Symbol)
		} else {
			symbolsFetcher.MarkPresentInWiki(symbol.Exchange, symbol.Symbol)
		}

		if err != nil {
			log.Println(err)
		} else {
			for _, day := range days {
				bucket.Store(fmt.Sprintf("%s/%s/%s.json", symbol.Exchange, symbol.Symbol, day.Date.Time.Format(time.RFC3339)), day.Serialize(symbol.Symbol, symbol.Exchange))
				symbolsFetcher.IncrementStoreCount()
			}
		}
		symbolsFetcher.UpdateSymbolFetched(symbol.Exchange, symbol.Symbol, mostRecentOpenDay)
		symbolsFetcher.SetLastUpdated(symbol.Exchange, symbol.Symbol)
	}
}
