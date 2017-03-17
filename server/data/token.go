package data

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/base64"
	"errors"
)

var (
	invalidTokenError = errors.New("invalid token")
)

type Token []byte

func NewToken(length int) (Token, error) {
	t := make([]byte, length)
	if _, err := rand.Read(t); err != nil {
		return nil, err
	}

	return t, nil
}

// UnmarshalBinary sets the value of the Token based on a slice of bytes.
func (t *Token) UnmarshalBinary(data []byte) error {
	*t = make([]byte, len(data))
	copy(*t, data)
	return nil
}

// MarshalText returns a base64-encoded slice of bytes.
func (t Token) MarshalText() (text []byte, err error) {
	data := make([]byte, base64.URLEncoding.EncodedLen(len(t)))
	base64.URLEncoding.Encode(data, t)
	return data, nil
}

// UnmarshalText sets the value of the Token based on a base64-encoded slice of bytes.
func (t *Token) UnmarshalText(text []byte) error {
	data := make([]byte, base64.URLEncoding.DecodedLen(len(text)))
	n, err := base64.URLEncoding.Decode(data, text)
	if err != nil {
		return err
	}

	*t = data[:n]
	return nil
}

// String returns a base64-encoded string.
func (t Token) String() string {
	return base64.URLEncoding.EncodeToString(t)
}

// Scan sets the value of the Token based on an interface.
func (id *Token) Scan(src interface{}) error {
	data, ok := src.([]byte)
	if !ok {
		return invalidTokenError
	}

	return id.UnmarshalBinary(data)
}

// Value implements the driver Valuer interface.
func (t Token) Value() (driver.Value, error) {
	return []byte(t), nil
}
