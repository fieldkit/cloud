package api

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/email"
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

func createEmailer(awsSession *session.Session, config *ApiConfiguration) (emailer email.Emailer, err error) {
	switch config.Emailer {
	case "default":
		emailer = email.NewEmailer("admin", config.Domain)
	case "aws":
		emailer = email.NewAWSSESEmailer(ses.New(awsSession), "admin", config.Domain)
	default:
		panic("invalid emailer")
	}
	return
}
