package email

import (
	"bytes"
	"fmt"
	html "html/template"
	text "text/template"
	"time"

	"github.com/fieldkit/cloud/server/data"
)

var (
	subjectTemplate  *text.Template
	bodyTextTemplate *text.Template
	bodyHTMLTemplate *html.Template
)

type templateOptions struct {
	ValidationToken *data.ValidationToken
	Source, Domain  string
}

func init() {
	var err error
	subjectTemplate, err = text.New("subject").Parse(subjectTemplateText)
	if err != nil {
		panic(err)
	}

	bodyTextTemplate, err = text.New("body").Parse(bodyTextTemplateText)
	if err != nil {
		panic(err)
	}

	bodyHTMLTemplate, err = html.New("body").Parse(bodyHTMLTemplateText)
	if err != nil {
		panic(err)
	}
}

type Emailer interface {
	SendValidationToken(person *data.User, validationToken *data.ValidationToken) error
	SendSourceSilenceWarning(source *data.Source, age time.Duration) error
}

type emailer struct {
	source string
	domain string
}

func (e emailer) SendSourceSilenceWarning(source *data.Source, age time.Duration) error {
	subject := fmt.Sprintf("FieldKit: Warning, source %s silent.", source.Name)
	body := fmt.Sprintf("The source '%s' has been silent for %v.", source.Name, age)
	fmt.Printf("To: %s\nSubject: %s\n\n%s\n\n", "jacob@conservify.org", subject, body)
	return nil
}

func (e emailer) SendValidationToken(person *data.User, validationToken *data.ValidationToken) error {
	options := &templateOptions{
		ValidationToken: validationToken,
		Source:          e.source,
		Domain:          e.domain,
	}

	subjectBuffer := bytes.NewBuffer([]byte{})
	if err := subjectTemplate.Execute(subjectBuffer, person); err != nil {
		return err
	}

	bodyBuffer := bytes.NewBuffer([]byte{})
	if err := bodyTextTemplate.Execute(bodyBuffer, options); err != nil {
		return err
	}

	fmt.Printf("To: %s\nSubject: %s\n\n%s\n\n", person.Email, subjectBuffer.String(), bodyBuffer.String())
	return nil
}

func NewEmailer(source, domain string) Emailer {
	return emailer{
		source: source,
		domain: domain,
	}
}

const (
	subjectTemplateText  = `Validate your Fieldkit account`
	bodyTextTemplateText = `To validate your Fieldkit account, navigate to:
https://api.{{.Domain}}/validate?token={{.ValidationToken.Token}}`
	bodyHTMLTemplateText = `<div style="font-family: 'Helvetica Neue','Helvetica',Arial,sans-serif;text-align:center;margin-top: 20px">
	<a href="https://api.{{.Domain}}/validate?token={{.ValidationToken.Token}}" style="text-decoration: none; color: rgb(0,0,0);">
		<div id="fieldkit" style="font-size: 24px; display: inline-block;clear:left; text-align: center; border: 1px solid rgb(0,0,0); padding: 20px 30px; margin-bottom:20px;border-radius: 3px;">
			<div>Validate your <strong style="color: #d0462c">Fieldkit</strong> account.</div>
		</div>
	</a>
	<div style="font-size: 10px;line-height: 1.5em">
		<div>Click above, or navigate to:</div>
		<div style="font-weight: bold">https://api.{{.Domain}}/validate?token={{.ValidationToken.Token}}</div>
	<div>
</div>`
)
