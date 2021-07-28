package data

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwtgo "github.com/dgrijalva/jwt-go"
)

const (
	bcryptCost = bcrypt.DefaultCost
)

var (
	IncorrectPasswordError = errors.New("incorrect password")
	UnverifiedUserError    = errors.New("unverified user error")
)

var (
	// RefreshTokenTtl = 2 * time.Minute
	// LoginTokenTtl   = 1 * time.Minute

	RefreshTokenTtl = 72 * time.Hour
	LoginTokenTtl   = 168 * time.Hour

	RecoveryTokenTtl     = 1 * time.Hour
	ValidationTokenTtl   = 72 * time.Hour
	TransmissionTokenTtl = 2 * 24 * 365 * time.Hour
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
	ID               int32     `db:"id,omitempty"`
	Name             string    `db:"name"`
	Username         string    `db:"username"`
	Email            string    `db:"email"`
	Password         []byte    `db:"password"`
	Valid            bool      `db:"valid"`
	Bio              string    `db:"bio"`
	MediaURL         *string   `db:"media_url"`
	MediaContentType *string   `db:"media_content_type"`
	Admin            bool      `db:"admin"`
	FirmwareTester   bool      `db:"firmware_tester"`
	FirmwarePattern  *string   `db:"firmware_pattern"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
	TncDate          time.Time `db:"tnc_date"`
}

func (user *User) SetPassword(password string) error {
	hashedPassword, err := generateHashFromPassword(password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return nil
}

func (user *User) CheckPassword(password string) error {
	hashedPassword, err := compareHashAndPassword(user.Password, password)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return IncorrectPasswordError
	}

	if err != nil {
		return err
	}

	if hashedPassword != nil {
		user.Password = hashedPassword
	}

	return nil
}

func (user *User) NewToken(now time.Time, refreshToken *RefreshToken) *jwtgo.Token {
	scopes := []string{"api:access"}

	if user.Admin {
		scopes = []string{"api:access", "api:admin"}
	}

	token := jwtgo.New(jwtgo.SigningMethodHS512)
	token.Claims = jwtgo.MapClaims{
		"iat":           now.Unix(),
		"exp":           now.Add(LoginTokenTtl).Unix(),
		"sub":           user.ID,
		"email":         user.Email,
		"refresh_token": refreshToken.Token.String(),
		"scopes":        scopes,
	}

	return token
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

type RecoveryToken struct {
	Token   Token     `db:"token" json:"token"`
	UserID  int32     `db:"user_id" json:"user_id"`
	Expires time.Time `db:"expires" json:"expires"`
}

func NewRecoveryToken(userID int32, length int, expires time.Time) (*RecoveryToken, error) {
	token, err := NewToken(length)
	if err != nil {
		return nil, err
	}

	return &RecoveryToken{
		UserID:  userID,
		Token:   token,
		Expires: expires,
	}, nil
}
