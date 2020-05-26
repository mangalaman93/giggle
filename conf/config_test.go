package conf

import (
	"testing"
)

func TestNewFrom(t *testing.T) {
	if _, err := newFrom("../config.json.example"); err != nil {
		t.Fatalf("error in loading config :: %v", err)
	}
}
