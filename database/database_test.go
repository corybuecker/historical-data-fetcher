package database

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/erikstmartin/go-testdb"
)

var testDB *sql.DB

func init() {
	testDB, _ = sql.Open("testdb", "")
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
	database := Database{DB: testDB}

	sql := "select id, symbol from symbols"
	testdb.StubQuery(sql, testdb.RowsFromCSVString([]string{"id", "symbol"}, "1,QQQ"))

	if err := database.FetchSymbols(); err != nil {
		t.Fatalf("should not have failed to fetch symbols")
	}
}
