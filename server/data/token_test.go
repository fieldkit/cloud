package data

import (
	"bytes"
	"testing"
)

const (
	testTokenLength = 20
)

func TestToken(t *testing.T) {
	token, err := NewToken(testTokenLength)
	if err != nil {
		t.Fatal(err)
	}

	if l := len(token); l != testTokenLength {
		t.Fatalf("expected token length %d, found %d", testTokenLength, l)
	}

	tokenString := token.String()
	t.Log(tokenString)

	decodedToken := Token{}
	if err := decodedToken.UnmarshalText([]byte(tokenString)); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(decodedToken, token) {
		t.Errorf("decoded token did not match original\n\t%s\n\t%s\n", decodedToken, token)
	}
}
