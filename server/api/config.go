package api

import (
	"fmt"
)

type ApiConfiguration struct {
	ApiDomain string
	ApiHost   string
}

func (ac *ApiConfiguration) MakeApiUrl(f string, args ...interface{}) string {
	return fmt.Sprintf("%s%s", ac.ApiHost, fmt.Sprintf(f, args...))
}
