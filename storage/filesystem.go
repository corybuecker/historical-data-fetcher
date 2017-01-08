package storage

import (
	"fmt"
	"io/ioutil"
	"log"
)

func WriteFile(key string, contents string) error {
	log.Println(fmt.Sprintf("writing /Users/stantonbuecker/Desktop/historicaldata/%s", key))
	return ioutil.WriteFile(fmt.Sprintf("/Users/stantonbuecker/Desktop/historicaldata/%s", key), []byte(contents), 0644)
}
