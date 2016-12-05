package backend

import (
	"testing"

	"github.com/O-C-R/fieldkit/data"
)

func TestBackend(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.Ping(); err != nil {
		t.Fatal(err)
	}
}

func TestBackendAddUser(t *testing.T) {
	backend, err := NewTestBackend()
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

func TestBackendUserByID(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	user, err := data.NewUser("test@ocr.nyc", "password")
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.AddUser(user); err != nil {
		t.Fatal(err)
	}

	returnedUser, err := backend.UserByID(user.ID)
	if err != nil {
		t.Fatal(err)
	}

	if returnedUser.ID != user.ID {
		t.Error("incorrect returned user ID")
	}
}

func TestBackendUserByValidationToken(t *testing.T) {
	backend, err := NewTestBackend()
	if err != nil {
		t.Fatal(err)
	}

	user, err := data.NewUser("test@ocr.nyc", "password")
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.AddUser(user); err != nil {
		t.Fatal(err)
	}

	returnedUser, err := backend.UserByValidationToken(user.ValidationToken)
	if err != nil {
		t.Fatal(err)
	}

	if returnedUser.ID != user.ID {
		t.Error("incorrect returned user ID")
	}
}
