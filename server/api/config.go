package api

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/fieldkit/cloud/server/email"
)

type BucketNames struct {
	Media   []string
	Streams []string
}

type ApiConfiguration struct {
	ApiDomain     string
	ApiHost       string
	SessionKey    string
	MapboxToken   string
	Emailer       string
	Domain        string
	PortalDomain  string
	EmailOverride string
	Buckets       *BucketNames
}

func (ac *ApiConfiguration) MakeApiUrl(f string, args ...interface{}) string {
	return fmt.Sprintf("%s%s", ac.ApiHost, fmt.Sprintf(f, args...))
}

func createEmailer(awsSession *session.Session, config *ApiConfiguration) (emailer email.Emailer, err error) {
	var overrides []*string
	if len(config.EmailOverride) > 0 {
		overrides = []*string{
			&config.EmailOverride,
		}
	}
	switch config.Emailer {
	case "default":
		return email.NewNoopEmailer("admin", config.Domain, overrides)
	case "aws":
		return email.NewAWSSESEmailer(ses.New(awsSession), "admin", config.Domain, overrides)
	default:
		panic("invalid emailer")
	}
}
