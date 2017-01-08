package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/corybuecker/historicaldata/calendar"
	"github.com/corybuecker/historicaldata/database"
	"github.com/corybuecker/historicaldata/parsers"
	"github.com/corybuecker/historicaldata/storage"
	"github.com/corybuecker/redisconfig"
	"github.com/davecgh/go-spew/spew"
	redis "gopkg.in/redis.v5"
)

type Config struct {
	TradierAPIKey string
	QuandlAPIKey  string
	S3Id          string
	S3Secret      string
}

func removeDuplicates(elements database.HistoricalData) database.HistoricalData {
	encountered := map[time.Time]bool{}
	result := make(database.HistoricalData, 0)

	for v := range elements {
		if encountered[elements[v].Date] == true {
		} else {
			encountered[elements[v].Date] = true
			result = append(result, elements[v])
		}
	}
	return result
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

	mostRecentOpenDay := calendar.GetMostRecentOpenDay(config.TradierAPIKey)
	log.Printf("the most recent open market day is: %s", mostRecentOpenDay.Format(time.RFC3339))

	symbolsFetcher := database.Database{Client: &database.RedisClient{Client: redis}}
	symbolsFetcher.LoadSymbolsNeedingUpdate(mostRecentOpenDay)

	wikiFetcher := parsers.BuildWikiParser(config.QuandlAPIKey)
	tradeFetcher := parsers.BuildTradierParser(config.TradierAPIKey)
	log.Printf("fetching %d symbols", len(symbolsFetcher.Symbols))

	for i, symbol := range symbolsFetcher.Symbols {
		var historicalData = make(database.HistoricalData, 0)

		existingSymbolBytes, exists := bucket.GetExistingSymbolBytes(fmt.Sprintf("%s/%s.json", symbol.Exchange, symbol.Symbol))

		if exists {
			json.Unmarshal(existingSymbolBytes, &historicalData)
		}

		var fetchedEntries database.HistoricalData
		var err error

		fetchedEntries, err = wikiFetcher.FetchIntoSlice(&symbol)

		if len(fetchedEntries) == 0 {
			fetchedEntries, err = tradeFetcher.FetchIntoSlice(&symbol)
		} else {
			symbolsFetcher.MarkPresentInWiki(symbol.Exchange, symbol.Symbol)
		}

		historicalData = append(historicalData, fetchedEntries...)

		historicalData = removeDuplicates(historicalData)
		sort.Sort(historicalData)

		if err != nil {
			log.Println(err)
		}

		newSymbolBytes, _ := json.Marshal(historicalData)
		bucket.Store(fmt.Sprintf("%s/%s.json", symbol.Exchange, symbol.Symbol), string(newSymbolBytes))

		symbolsFetcher.IncrementDateCount(mostRecentOpenDay.Format(time.RFC3339))
		symbolsFetcher.UpdateSymbolFetched(symbol.Exchange, symbol.Symbol, historicalData.MostRecentDay())
		symbolsFetcher.SetLastUpdated(symbol.Exchange, symbol.Symbol)

		log.Printf("remaining symbols %d", len(symbolsFetcher.Symbols)-i)
	}
}
