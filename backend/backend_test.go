package backend

import (
	"testing"

	"github.com/O-C-R/fieldkit/data"
)

const (
	testURL = "mongodb://localhost/test"
)

func testBackend() (*Backend, error) {
	backend, err := newBackend(testURL, false)
	if err != nil {
		return nil, err
	}

	if err := backend.dropDatabase(); err != nil {
		return nil, err
	}

	if err := backend.init(); err != nil {
		return nil, err
	}

	return backend, nil
}

func TestBackend(t *testing.T) {
	backend, err := NewBackend(testURL, false)
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.Ping(); err != nil {
		t.Fatal(err)
	}
}

func TestBackendAddUser(t *testing.T) {
	backend, err := testBackend()
	if err != nil {
		t.Fatal(err)
	}

	user, err := data.NewUser("test@ocr.nyc", "password")
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.AddUser(user); err != nil {
		t.Error(err)
	}

	if err := backend.AddUser(user); err != DuplicateKeyError {
		t.Error("duplicate key succeeded")
	}
}
