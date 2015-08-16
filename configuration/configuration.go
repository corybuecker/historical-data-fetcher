package configuration

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	ParserToken      string
	ParserType       string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseHost     string
}

func (configuration *Configuration) Load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(configuration); err != nil {
		return err
	}
	return nil
}
