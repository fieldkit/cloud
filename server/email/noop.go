package email

import (
	"bytes"
	"fmt"

	"github.com/fieldkit/cloud/server/data"
)

type noopEmailer struct {
	templates *EmailTemplates
	source    string
	domain    string
}

func NewNoopEmailer(source, domain string) (e Emailer, err error) {
	templates, err := NewEmailTemplates()
	if err != nil {
		return nil, err
	}

	e = &noopEmailer{
		templates: templates,
		source:    source,
		domain:    domain,
	}

	return
}

type templateOptions struct {
	ValidationToken *data.ValidationToken
	Source, Domain  string
}

func (e noopEmailer) SendValidationToken(person *data.User, validationToken *data.ValidationToken) error {
	options := &templateOptions{
		ValidationToken: validationToken,
		Source:          e.source,
		Domain:          e.domain,
	}

	subjectBuffer := bytes.NewBuffer([]byte{})
	if err := e.templates.Validation.Subject.Execute(subjectBuffer, person); err != nil {
		return err
	}

	bodyBuffer := bytes.NewBuffer([]byte{})
	if err := e.templates.Validation.BodyText.Execute(bodyBuffer, options); err != nil {
		return err
	}

	fmt.Printf("To: %s\nSubject: %s\n\n%s\n\n", person.Email, subjectBuffer.String(), bodyBuffer.String())

	return nil
}
