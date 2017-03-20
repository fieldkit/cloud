package email

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/O-C-R/fieldkit/server/data"
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

func (a *AWSSESEmailer) SendValidationToken(person *data.User, validationToken *data.ValidationToken) error {
	subjectBuffer := bytes.NewBuffer([]byte{})
	if err := subjectTemplate.Execute(subjectBuffer, person); err != nil {
		return err
	}

	bodyTextBuffer := bytes.NewBuffer([]byte{})
	if err := bodyTextTemplate.Execute(bodyTextBuffer, validationToken); err != nil {
		return err
	}

	bodyHTMLBuffer := bytes.NewBuffer([]byte{})
	if err := bodyHTMLTemplate.Execute(bodyHTMLBuffer, validationToken); err != nil {
		return err
	}

	sendEmailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(person.Email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data:    aws.String(bodyTextBuffer.String()),
					Charset: aws.String("utf8"),
				},
				Html: &ses.Content{
					Data:    aws.String(bodyHTMLBuffer.String()),
					Charset: aws.String("utf8"),
				},
			},
			Subject: &ses.Content{
				Data:    aws.String(subjectBuffer.String()),
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
