package parsers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var server *httptest.Server

func init() {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{'test': true}")
	}))
}

func TestCreatedClient(t *testing.T) {
	parser := new(TradierParser)
	parser.fetch(server.URL)
	if parser.Client == nil {
		t.Fatalf("should have created a parser")
	}
}

func TestSuccessfulFetch(t *testing.T) {
	parser := new(TradierParser)
	if value, _ := parser.fetch(server.URL); string(value) != "{'test': true}" {
		t.Fatalf("should have received %s, but got %s", "{'test': true}", string(value))
	}
}

func TestErrorFromClient(t *testing.T) {
	parser := TradierParser{Client: &http.Client{Timeout: time.Nanosecond}}
	if _, err := parser.fetch(server.URL); err == nil {
		t.Fatalf("should have failed, %s", err)
	}
}

func TestNonTwoHundred(t *testing.T) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	parser := new(TradierParser)
	if _, err := parser.fetch(server.URL); err.Error() != "the call to the API failed with status code 500" {
		t.Fatalf("should have failed, %s", err.Error())
	}
}
