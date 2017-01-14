package id

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/base32"
	"errors"
)

var (
	InvalidIDError = errors.New("invalid ID")
)

// ID is a unique identifier.
type ID [20]byte

// New returns a random ID value.
func New() (ID, error) {
	id := ID{}
	if _, err := rand.Read(id[:]); err != nil {
		return id, err
	}

	return id, nil
}

// MarshalBinary returns a slice of bytes.
func (id ID) MarshalBinary() ([]byte, error) {
	return id[:], nil
}

// UnmarshalText sets the value of the ID based on a slice of bytes.
func (id *ID) UnmarshalBinary(data []byte) error {
	if len(data) != len(id) {
		return InvalidIDError
	}

	copy(id[:], data)
	return nil
}

// MarshalText returns a base32-encoded slice of bytes.
func (id ID) MarshalText() (text []byte, err error) {
	data := make([]byte, base32.StdEncoding.EncodedLen(len(id)))
	base32.StdEncoding.Encode(data, id[:])
	return data, nil
}

// UnmarshalText sets the value of the ID based on a base32-encoded slice of bytes.
func (id *ID) UnmarshalText(text []byte) error {
	data := make([]byte, base32.StdEncoding.DecodedLen(len(text)))
	if _, err := base32.StdEncoding.Decode(data, text); err != nil {
		return err
	}

	return id.UnmarshalBinary(data)
}

// String returns a hex-encoded string.
func (id ID) String() string {
	return base32.StdEncoding.EncodeToString(id[:])
}

// Scan sets the value of the ID based on an interface.
func (id *ID) Scan(src interface{}) error {
	data, ok := src.([]byte)
	if !ok {
		return InvalidIDError
	}

	return id.UnmarshalBinary(data)
}

// Value implements the driver Valuer interface.
func (id ID) Value() (driver.Value, error) {
	return id[:], nil
}
