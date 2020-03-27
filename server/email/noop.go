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
	override  []*string
}

func NewNoopEmailer(source, domain string, override []*string) (e Emailer, err error) {
	templates, err := NewEmailTemplates()
	if err != nil {
		return nil, err
	}

	e = &noopEmailer{
		templates: templates,
		source:    source,
		domain:    domain,
		override:  override,
	}

	return
}

type templateOptions struct {
	ValidationToken *data.ValidationToken
	RecoveryToken   *data.RecoveryToken
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

	if e.override != nil {
		fmt.Printf("Override: %v", e.override)
	}

	fmt.Printf("To: %s\nSubject: %s\n\n%s\n\n", person.Email, subjectBuffer.String(), bodyBuffer.String())

	return nil
}

func (e noopEmailer) SendRecoveryToken(person *data.User, recoveryToken *data.RecoveryToken) error {
	options := &templateOptions{
		RecoveryToken: recoveryToken,
		Source:        e.source,
		Domain:        e.domain,
	}

	subjectBuffer := bytes.NewBuffer([]byte{})
	if err := e.templates.Recovery.Subject.Execute(subjectBuffer, person); err != nil {
		return err
	}

	bodyBuffer := bytes.NewBuffer([]byte{})
	if err := e.templates.Recovery.BodyText.Execute(bodyBuffer, options); err != nil {
		return err
	}

	if e.override != nil {
		fmt.Printf("Override: %v", e.override)
	}

	fmt.Printf("To: %s\nSubject: %s\n\n%s\n\n", person.Email, subjectBuffer.String(), bodyBuffer.String())

	return nil
}
