package api

import (
	"encoding/base64"
	"fmt"

	"github.com/goadesign/goa"
)

type ApiConfiguration struct {
	BucketName string
	ApiDomain  string
	ApiHost    string
	SessionKey string
	Emailer    string
	Domain     string
}

func (ac *ApiConfiguration) MakeApiUrl(f string, args ...interface{}) string {
	return fmt.Sprintf("%s%s", ac.ApiHost, fmt.Sprintf(f, args...))
}

func (ac *ApiConfiguration) NewJWTMiddleware() (goa.Middleware, error) {
	jwtHMACKey, err := base64.StdEncoding.DecodeString(ac.SessionKey)
	if err != nil {
		return nil, err
	}

	jwtMiddleware, err := NewJWTMiddleware(jwtHMACKey)
	if err != nil {
		return nil, err
	}

	return jwtMiddleware, nil
}
