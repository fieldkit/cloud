package data

import (
	"errors"
	"time"

	"github.com/O-C-R/auth/id"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = bcrypt.DefaultCost
)

var (
	IncorrectPasswordError = errors.New("incorrect password")
)

func generateHashFromPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
}

func compareHashAndPassword(hashedPassword []byte, password string) ([]byte, error) {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return nil, err
	}

	passwordBcryptCost, err := bcrypt.Cost(hashedPassword)
	if err != nil {
		return nil, err
	}

	if passwordBcryptCost != bcryptCost {
		return generateHashFromPassword(password)
	}

	return nil, nil
}

type User struct {
	ID                    id.ID     `bson:"id" json:"-"`
	Email                 string    `bson:"email" json:"email"`
	Password              []byte    `bson:"password" json:"-"`
	Valid                 bool      `bson:"valid" json:"valid"`
	ValidationToken       id.ID     `bson:"validation_token" json:"-"`
	ValidationTokenExpire time.Time `bson:"validation_token_expire" json:"-"`
}

func NewUser(email, password string) (*User, error) {
	userID, err := id.New()
	if err != nil {
		return nil, err
	}

	validationToken, err := id.New()
	if err != nil {
		return nil, err
	}

	u := &User{
		ID:                    userID,
		Email:                 email,
		ValidationToken:       validationToken,
		ValidationTokenExpire: time.Now().Add(8 * time.Hour),
	}

	if err := u.SetPassword(password); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := generateHashFromPassword(password)
	if err != nil {
		return err
	}

	u.Password = hashedPassword
	return nil
}

func (u *User) CheckPassword(password string) error {
	hashedPassword, err := compareHashAndPassword(u.Password, password)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return IncorrectPasswordError
	}

	if err != nil {
		return err
	}

	if hashedPassword != nil {
		u.Password = hashedPassword
	}

	return nil
}
