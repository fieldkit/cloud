package data

import (
	"errors"
	"time"

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
	ID       int32  `db:"id,omitempty" json:"id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Password []byte `db:"password" json:"-"`
	Valid    bool   `db:"valid" json:"valid"`
}

func (p *User) SetPassword(password string) error {
	hashedPassword, err := generateHashFromPassword(password)
	if err != nil {
		return err
	}

	p.Password = hashedPassword
	return nil
}

func (p *User) CheckPassword(password string) error {
	hashedPassword, err := compareHashAndPassword(p.Password, password)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return IncorrectPasswordError
	}

	if err != nil {
		return err
	}

	if hashedPassword != nil {
		p.Password = hashedPassword
	}

	return nil
}

type ValidationToken struct {
	Token   Token     `db:"token" json:"token"`
	UserID  int32     `db:"user_id" json:"user_id"`
	Expires time.Time `db:"expires" json:"expires"`
}

func NewValidationToken(userID int32, length int, expires time.Time) (*ValidationToken, error) {
	token, err := NewToken(length)
	if err != nil {
		return nil, err
	}

	return &ValidationToken{
		UserID:  userID,
		Token:   token,
		Expires: expires,
	}, nil
}

type RefreshToken struct {
	Token   Token     `db:"token" json:"token"`
	UserID  int32     `db:"user_id" json:"user_id"`
	Expires time.Time `db:"expires" json:"expires"`
}

func NewRefreshToken(userID int32, length int, expires time.Time) (*RefreshToken, error) {
	token, err := NewToken(length)
	if err != nil {
		return nil, err
	}

	return &RefreshToken{
		UserID:  userID,
		Token:   token,
		Expires: expires,
	}, nil
}
