package email

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/fieldkit/cloud/server/data"
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

func (a AWSSESEmailer) send(subject, body string, addresses []*string) error {
	sendEmailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String("jacob@conservify.org"),
				aws.String("shah@conservify.org"),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data:    aws.String(body),
					Charset: aws.String("utf8"),
				},
			},
			Subject: &ses.Content{
				Data:    aws.String(subject),
				Charset: aws.String("utf8"),
			},
		},
		Source:     a.sourceEmail,
		ReturnPath: a.sourceEmail,
		ReplyToAddresses: []*string{
			a.sourceEmail,
		},
	}

	if _, err := a.client.SendEmail(sendEmailInput); err != nil {
		return err
	}

	return nil
}

func (a AWSSESEmailer) SendSourceOfflineWarning(source *data.Source, age time.Duration) error {
	subject := fmt.Sprintf("FieldKit: Device %s is offline.", source.Name)
	body := fmt.Sprintf("Your device named '%s' is offline. The last reading from the device was %v ago.", source.Name, age)
	toAddresses := []*string{
		aws.String("jacob@conservify.org"),
		aws.String("shah@conservify.org"),
	}

	return a.send(subject, body, toAddresses)
}

func (a AWSSESEmailer) SendSourceOnlineWarning(source *data.Source, age time.Duration) error {
	subject := fmt.Sprintf("FieldKit: Device %s is back online.", source.Name)
	body := fmt.Sprintf("Your device named '%s' is online. The last reading from the device was %v ago.", source.Name, age)
	toAddresses := []*string{
		aws.String("jacob@conservify.org"),
		aws.String("shah@conservify.org"),
	}

	return a.send(subject, body, toAddresses)
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

	toAddresses := []*string{
		aws.String(person.Email),
	}

	return a.send(subjectBuffer.String(), bodyTextBuffer.String(), toAddresses)
}

func (a AWSSESEmailer) SendInvitation(email string) error {
	subject := fmt.Sprintf("You are invited to a FieldKit project!")
	body := fmt.Sprintf("Soon this will be a real invitation, but for now, just consider yourself invited.")
	toAddresses := []*string{
		aws.String(email),
	}

	return a.send(subject, body, toAddresses)
}
