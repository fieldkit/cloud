package email

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/O-C-R/fieldkit/server/data"
)

type AWSSESEmailer struct {
	client      *ses.SES
	source      string
	domain      string
	sourceEmail *string
}

func NewAWSSESEmailer(client *ses.SES, source, domain string) *AWSSESEmailer {
	return &AWSSESEmailer{
		client:      client,
		source:      source,
		domain:      domain,
		sourceEmail: aws.String(source + "@" + domain),
	}
}

func (a *AWSSESEmailer) SendValidationToken(person *data.User, validationToken *data.ValidationToken) error {
	options := &templateOptions{
		ValidationToken: validationToken,
		Source:          a.source,
		Domain:          a.domain,
	}

	subjectBuffer := bytes.NewBuffer([]byte{})
	if err := subjectTemplate.Execute(subjectBuffer, person); err != nil {
		return err
	}

	bodyTextBuffer := bytes.NewBuffer([]byte{})
	if err := bodyTextTemplate.Execute(bodyTextBuffer, options); err != nil {
		return err
	}

	bodyHTMLBuffer := bytes.NewBuffer([]byte{})
	if err := bodyHTMLTemplate.Execute(bodyHTMLBuffer, options); err != nil {
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
		Source: a.sourceEmail,
		ReplyToAddresses: []*string{
			a.sourceEmail,
		},
		ReturnPath: a.sourceEmail,
	}

	if _, err := a.client.SendEmail(sendEmailInput); err != nil {
		return err
	}

	return nil
}
