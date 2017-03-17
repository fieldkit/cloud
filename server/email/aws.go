package email

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

type AWSSESEmailer struct {
	client *ses.SES
	source *string
}

func NewAWSSESEmailer(client *ses.SES, source string) *AWSSESEmailer {
	return &AWSSESEmailer{
		client: client,
		source: aws.String(source),
	}
}

func (a *AWSSESEmailer) SendEmail(address, subject, message string) error {
	sendEmailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(address),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data:    aws.String(message),
					Charset: aws.String("utf8"),
				},
			},
			Subject: &ses.Content{
				Data:    aws.String(subject),
				Charset: aws.String("utf8"),
			},
		},
		Source: a.source,
		ReplyToAddresses: []*string{
			a.source,
		},
		ReturnPath: a.source,
	}

	if _, err := a.client.SendEmail(sendEmailInput); err != nil {
		return err
	}

	return nil
}
