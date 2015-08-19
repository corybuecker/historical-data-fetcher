package parsers

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GoogleParser struct {
}

func (parser *GoogleParser) fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (parser *GoogleParser) Read(symbol string, symbolID string) ([]Tick, error) {
	var test []Tick
	var body []byte
	var err error
	if body, err = parser.fetch(fmt.Sprintf("https://www.google.com/finance/getprices?q=%s&i=60&p=1d", symbol)); err != nil {
		return nil, err
	}
	reader := bytes.NewReader(body)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		test = append(test, Tick{})
	}

	return test, nil
}
