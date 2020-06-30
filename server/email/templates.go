package email

import (
	html "html/template"
	text "text/template"

	"github.com/fieldkit/cloud/server/data"
)

type templateOptions struct {
	Sender          *data.User
	Invite          *data.ProjectInvite
	ValidationToken *data.ValidationToken
	RecoveryToken   *data.RecoveryToken
	Source, Domain  string
}

type EmailTemplates struct {
	Validation        *EmailTemplate
	Recovery          *EmailTemplate
	ProjectInvitation *EmailTemplate
}

func NewEmailTemplates() (e *EmailTemplates, err error) {
	validation, err := NewValidationEmailTemplate()
	if err != nil {
		return nil, err
	}

	recovery, err := NewRecoveryEmailTemplate()
	if err != nil {
		return nil, err
	}

	projectInvitation, err := NewProjectInvitationEmailTemplate()
	if err != nil {
		return nil, err
	}

	e = &EmailTemplates{
		Validation:        validation,
		Recovery:          recovery,
		ProjectInvitation: projectInvitation,
	}

	return
}

type EmailTemplate struct {
	Subject  *text.Template
	BodyText *text.Template
	BodyHTML *html.Template
}

func NewValidationEmailTemplate() (et *EmailTemplate, err error) {
	subjectText := `Validate your Fieldkit account`

	bodyTextText := `To validate your Fieldkit account, navigate to:
https://api.{{.Domain}}/validate?token={{.ValidationToken.Token}}`

	bodyHTMLText := `<div style="font-family: 'Helvetica Neue','Helvetica',Arial,sans-serif;text-align:center;margin-top: 20px">
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

	subject, err := text.New("subject").Parse(subjectText)
	if err != nil {
		return nil, err
	}

	bodyText, err := text.New("body").Parse(bodyTextText)
	if err != nil {
		return nil, err
	}

	bodyHTML, err := html.New("body").Parse(bodyHTMLText)
	if err != nil {
		return nil, err
	}

	et = &EmailTemplate{
		Subject:  subject,
		BodyText: bodyText,
		BodyHTML: bodyHTML,
	}

	return
}

func NewRecoveryEmailTemplate() (et *EmailTemplate, err error) {
	subjectText := `Recover your Fieldkit account`

	bodyTextText := `To recover your Fieldkit account, navigate to:
https://portal.{{.Domain}}/recover/complete?token={{.RecoveryToken.Token}}`

	bodyHTMLText := bodyTextText

	subject, err := text.New("subject").Parse(subjectText)
	if err != nil {
		return nil, err
	}

	bodyText, err := text.New("body").Parse(bodyTextText)
	if err != nil {
		return nil, err
	}

	bodyHTML, err := html.New("body").Parse(bodyHTMLText)
	if err != nil {
		return nil, err
	}

	et = &EmailTemplate{
		Subject:  subject,
		BodyText: bodyText,
		BodyHTML: bodyHTML,
	}

	return
}

func NewProjectInvitationEmailTemplate() (et *EmailTemplate, err error) {
	subjectText := `You're invited!`

	bodyTextText := `You have been invited to new FieldKit project. To view the invitation, click here:
https://portal.{{.Domain}}/projects/invitation?token={{.Invite.Token}}`

	bodyHTMLText := bodyTextText

	subject, err := text.New("subject").Parse(subjectText)
	if err != nil {
		return nil, err
	}

	bodyText, err := text.New("body").Parse(bodyTextText)
	if err != nil {
		return nil, err
	}

	bodyHTML, err := html.New("body").Parse(bodyHTMLText)
	if err != nil {
		return nil, err
	}

	et = &EmailTemplate{
		Subject:  subject,
		BodyText: bodyText,
		BodyHTML: bodyHTML,
	}

	return
}
