package backend

import (
	"testing"
)

const (
	testURL = "mongodb://localhost/test"
)

func TestBackend(t *testing.T) {
	backend, err := NewBackend(testURL, false)
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.Ping(); err != nil {
		t.Fatal(err)
	}
}
