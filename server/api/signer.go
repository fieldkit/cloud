package api

import (
	"errors"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type Signer struct {
	key []byte
}

func NewSigner(key []byte) (s *Signer) {
	return &Signer{
		key: key,
	}
}

func (s *Signer) Verify(token string) error {
	claims := make(jwtgo.MapClaims)
	_, err := jwtgo.ParseWithClaims(token, claims, func(t *jwtgo.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	if claims["exp"] == nil {
		return errors.New("auth missing expiration")
	}

	unix := claims["exp"].(float64)
	expiration := time.Unix(int64(unix), 0)
	now := time.Now()

	if now.After(expiration) {
		return errors.New("auth expired")
	}

	return nil
}

func (s *Signer) SignURL(url string) (string, error) {
	now := time.Now()
	token := jwtgo.New(jwtgo.SigningMethodHS512)
	token.Claims = jwtgo.MapClaims{
		"exp": now.Add(time.Hour * 1).Unix(),
	}

	signed, err := token.SignedString(s.key)
	if err != nil {
		return "", err
	}

	return url + "?auth=" + signed, nil
}

func (s *Signer) SignAndBustURL(url string, key *string) (*string, error) {
	if key == nil {
		return nil, nil
	}

	signed, err := s.SignURL(url)
	if err != nil {
		return nil, err
	}

	return &signed, nil
}
