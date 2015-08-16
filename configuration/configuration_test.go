package configuration

import "testing"

func TestMissingFile(t *testing.T) {
	configuration := Configuration{}
	err := configuration.Load("doesnotexist")

	if err == nil {
		t.Fatalf("should have failed to open missing file")
	}
}
