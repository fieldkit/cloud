package api

import (
	"fmt"
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
