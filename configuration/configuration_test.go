package configuration

import "testing"

func TestMissingFile(t *testing.T) {
	configuration := Configuration{}
	err := configuration.Load("doesnotexist")

	if err == nil {
		t.Fatalf("should have failed to open missing file")
	}
}
func TestInvalidFile(t *testing.T) {
	configuration := Configuration{}
	err := configuration.Load("../__test/invalid_config.json")

	if err == nil {
		t.Fatalf("should have failed to parse bad file")
	}
}

func TestValidFile(t *testing.T) {
	configuration := Configuration{}
	err := configuration.Load("../__test/valid_config.json")

	if err != nil {
		t.Fatalf("should have parsed valid file")
	}

	if configuration.DatabaseName != "database" {
		t.Fatalf("should have parsed valid file")
	}
}
