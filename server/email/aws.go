package email

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/fieldkit/cloud/server/data"
)

type AWSSESEmailer struct {
	source      string
	domain      string
	sourceEmail *string
	override    []*string
	templates   *EmailTemplates
	sesc        *ses.SES
}

func NewAWSSESEmailer(sesc *ses.SES, source, domain string, override []*string) (e *AWSSESEmailer, err error) {
	templates, err := NewEmailTemplates()
	if err != nil {
		return nil, err
	}

	e = &AWSSESEmailer{
		source:      source,
		domain:      domain,
		sourceEmail: aws.String(source + "@" + domain),
		override:    override,
		templates:   templates,
		sesc:        sesc,
	}

	return
}

func (a AWSSESEmailer) send(subject, body string, addresses []*string) error {
	destination := addresses
	if a.override != nil {
		destination = a.override
	}

	sendEmailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: destination,
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

	if _, err := a.sesc.SendEmail(sendEmailInput); err != nil {
		return err
	}

	return nil
}

func (a *AWSSESEmailer) SendValidationToken(person *data.User, validationToken *data.ValidationToken) error {
	options := &templateOptions{
		ValidationToken: validationToken,
		Source:          a.source,
		Domain:          a.domain,
	}

	subjectBuffer := bytes.NewBuffer([]byte{})
	if err := a.templates.Validation.Subject.Execute(subjectBuffer, person); err != nil {
		return err
	}

	bodyTextBuffer := bytes.NewBuffer([]byte{})
	if err := a.templates.Validation.BodyText.Execute(bodyTextBuffer, options); err != nil {
		return err
	}

	toAddresses := []*string{
		aws.String(person.Email),
	}

	return a.send(subjectBuffer.String(), bodyTextBuffer.String(), toAddresses)
}
