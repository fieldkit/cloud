package data

import (
	"testing"
)

func TestUserPassword(t *testing.T) {
	user, err := NewUser("test@ocr.nyc", "test", "password")
	if err != nil {
		t.Fatal(err)
	}

	if err := user.CheckPassword("password"); err != nil {
		t.Error(err)
	}

	if err := user.CheckPassword(""); err != IncorrectPasswordError {
		t.Error("empty password succeeded")
	}

	if err := user.CheckPassword("incorrect password"); err != IncorrectPasswordError {
		t.Error("incorrect password succeeded")
	}
}
