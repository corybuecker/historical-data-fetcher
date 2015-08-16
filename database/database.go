package database

import (
	"database/sql"
	"log"

	"github.com/corybuecker/trade-fetcher/parsers"
	_ "github.com/lib/pq"
)

type Symbol struct {
	Id     string
	Symbol string
}
type Database struct {
	DB      *sql.DB
	Symbols []Symbol
}

func (database *Database) Connect(dataSourceName string) error {
	var err error
	database.DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) FetchSymbols() error {
	rows, err := database.DB.Query("SELECT id, symbol FROM symbols")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, symbol string
		if err := rows.Scan(&id, &symbol); err != nil {
			log.Fatal(err)
		}
		database.Symbols = append(database.Symbols, Symbol{Id: id, Symbol: symbol})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (database *Database) InsertTick(tick parsers.Tick) error {
	if _, err := database.DB.Exec("insert into trades (symbol_id, time, price, volume) values ($1, $2, $3, $4)", tick.SymbolID, tick.Time.UTC(), tick.Price, tick.Volume); err != nil {
		return err
	}
	return nil
}
