package database

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

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

func TestValidConnection(t *testing.T) {
	database := Database{}
	err := database.Connect(fmt.Sprintf("host=localhost user=ubuntu dbname=circle_test sslmode=disable"))
	if err != nil {
		t.Fatalf("should have successfully connected to the database")
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
