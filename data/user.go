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
	ID        id.ID  `db:"id" json:"id"`
	Username  string `db:"username" json:"username"`
	Email     string `db:"email" json:"email"`
	Password  []byte `db:"password" json:"-"`
	FirstName string `db:"first_name" json:"first_name,omitempty"`
	LastName  string `db:"last_name" json:"last_name,omitempty"`
	Valid     bool   `db:"valid" json:"valid"`
}

func NewUser(email, username, password string) (*User, error) {
	userID, err := id.New()
	if err != nil {
		return nil, err
	}

	u := &User{
		ID:       userID,
		Username: username,
		Email:    email,
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

type ValidationToken struct {
	ID      id.ID     `db:"id"`
	UserID  id.ID     `db:"user_id"`
	Expires time.Time `db:"expires"`
}

func NewValidationToken(userID id.ID) (*ValidationToken, error) {
	validationTokenID, err := id.New()
	if err != nil {
		return nil, err
	}

	return &ValidationToken{
		ID:      validationTokenID,
		UserID:  userID,
		Expires: time.Now().Add(8 * time.Hour),
	}, nil
}
