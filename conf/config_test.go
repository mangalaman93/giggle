package conf

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	if _, err := ReadConfig("../config.json.example"); err != nil {
		t.Fatalf("error in loading config :: %v", err)
	}
}
