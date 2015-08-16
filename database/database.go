package database

import (
	"database/sql"

	"github.com/corybuecker/trade-fetcher/parsers"
	_ "github.com/lib/pq" // blank import for postgresql driver
)

type Symbol struct {
	ID     string
	Symbol string
}
type Database struct {
	DB      *sql.DB
	Symbols []Symbol
}

func (database *Database) Connect(dataSourceName string) error {
	// Note that the error returned here isn't used. The Ping call below is used to confirm the connection.
	database.DB, _ = sql.Open("postgres", dataSourceName)
	if err := database.DB.Ping(); err != nil {
		return err
	}
	return nil
}

func (database *Database) FetchSymbols() error {
	var rows *sql.Rows
	var err error
	if rows, err = database.DB.Query("select id, symbol from symbols"); err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id, symbol string
		rows.Scan(&id, &symbol)
		database.Symbols = append(database.Symbols, Symbol{ID: id, Symbol: symbol})
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func (database *Database) InsertTick(tick parsers.Tick) error {
	if _, err := database.DB.Exec("insert into trades (symbol_id, time, price, volume) values ($1, $2, $3, $4)", tick.SymbolID, tick.Time.UTC(), tick.Price, tick.Volume); err != nil {
		return err
	}
	return nil
}
