package main

import (
	"fmt"
	"log"

	"github.com/corybuecker/trade-fetcher/configuration"
	"github.com/corybuecker/trade-fetcher/database"
	"github.com/corybuecker/trade-fetcher/parsers"
)

func main() {
	var config = new(configuration.Configuration)
	var database = &database.Database{}
	if err := config.Load("./config.json"); err != nil {
		log.Fatal(err)
	}
	var parser parsers.Parser
	if config.ParserType == "Tradier" {
		parser = &parsers.TradierParser{}
	} else {
		parser = &parsers.GoogleParser{}
	}

	if err := database.Connect(fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.DatabaseHost, config.DatabaseUser, config.DatabasePassword, config.DatabaseName)); err != nil {
		log.Fatal(err)
	}

	database.FetchSymbols()

	fmt.Println(database)
	for _, symbol := range database.Symbols {
		ticks, _ := parser.Read(symbol.Symbol, symbol.Id)
		for _, tick := range ticks {
			if err := database.InsertTick(tick); err != nil {
				log.Println(err)
			}
		}
	}
}
