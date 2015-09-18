package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/corybuecker/trade-fetcher/parsers"
	"github.com/erikstmartin/go-testdb"
)

var testDB *sql.DB
var database Database

func init() {
	testDB, _ = sql.Open("testdb", "")
	database = Database{DB: testDB}
}

func TestInvalidConnection(t *testing.T) {
	database := Database{}
	err := database.Connect("database")
	if err == nil {
		t.Fatalf("should have failed the database connection")
	}
}

func TestFetchSymbols(t *testing.T) {
	sql := "select id, symbol from symbols"
	testdb.StubQuery(sql, testdb.RowsFromCSVString([]string{"id", "symbol"}, "1,QQQ"))

	if err := database.FetchSymbols(); err != nil {
		t.Fatalf("should not have failed to fetch symbols")
	}
}

func TestFetchSymbolsError(t *testing.T) {
	sql := "select id, symbol from symbols"
	testdb.StubQueryError(sql, errors.New("test"))

	if err := database.FetchSymbols(); err == nil {
		t.Fatalf("should have failed to fetch symbols")
	}
}

func TestInsertTick(t *testing.T) {
	testdb.StubExec("insert into trades (symbol_id, time, price, volume) values ($1, $2, $3, $4)", testdb.NewResult(0, nil, 1, nil))
	tick := parsers.Tick{}
	err := database.InsertTick(tick)

	if err != nil {
		t.Fatalf("should not have errored")
	}
}

func TestInsertTickError(t *testing.T) {
	testdb.StubExecError("insert into trades (symbol_id, time, price, volume) values ($1, $2, $3, $4)", errors.New("error"))
	tick := parsers.Tick{}
	err := database.InsertTick(tick)
	if err == nil {
		t.Fatalf("should not have errored")
	}
}
