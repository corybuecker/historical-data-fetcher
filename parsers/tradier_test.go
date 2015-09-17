package parsers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

var server *httptest.Server

func buildServer(headers http.Header) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for header, value := range headers {
			w.Header().Set(header, value[0])
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprint(w, "{'test': true}")
	}))
}

func init() {
	log.Println("test")

	headers := http.Header{}
	headers.Add("X-Ratelimit-Available", "59")
	headers.Add("X-Ratelimit-Expiry", strconv.FormatInt(time.Now().Add(time.Second*1).Unix()*1000, 10))
	buildServer(headers)
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

func TestRateLimitDivideZero(t *testing.T) {
	headers := http.Header{}
	headers.Add("X-Ratelimit-Available", "0")
	headers.Add("X-Ratelimit-Expiry", strconv.FormatInt(time.Now().Add(time.Second*10).Unix()*1000, 10))
	buildServer(headers)
	parser := new(TradierParser)

	if _, err := parser.fetch(server.URL); err != nil {
		t.Fatalf("should not have failed, %s", err.Error())
	}
}
