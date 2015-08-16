package parsers

import (
	"bufio"
	"fmt"
	"net/http"
)

type GoogleParser struct {
}

func (parser *GoogleParser) Read(symbol string, symbolID string) ([]Tick, error) {
	var test []Tick

	resp, err := http.Get(fmt.Sprintf("https://www.google.com/finance/getprices?q=%s&i=60&p=1d", symbol))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		test = append(test, Tick{})
	}

	return test, nil
}
