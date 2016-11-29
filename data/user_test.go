package data

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestUserPassword(t *testing.T) {
	user, err := NewUser("test@ocr.nyc", "password")
	if err != nil {
		t.Fatal(err)
	}

	if err := user.CheckPassword("password"); err != nil {
		t.Error(err)
	}

	if err := user.CheckPassword(""); err != bcrypt.ErrMismatchedHashAndPassword {
		t.Error("empty password succeeded")
	}

	if err := user.CheckPassword("incorrect password"); err != bcrypt.ErrMismatchedHashAndPassword {
		t.Error("incorrect password succeeded")
	}
}
