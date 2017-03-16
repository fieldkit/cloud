package email

import (
	"fmt"

	"github.com/O-C-R/fieldkit/server/data"
)

type Emailer interface {
	SendValidationToken(person *data.User, validationToken *data.ValidationToken) error
}

type emailer struct{}

func (e emailer) SendValidationToken(person *data.User, validationToken *data.ValidationToken) error {
	fmt.Printf("To: %s\nSubject: %s\n\n%s\n", person.Email)
	return nil
}

func NewEmailer() Emailer {
	return emailer{}
}

const (
	validationTokenTextTemplate = `

`
	validationTokenHTMLTemplate = `
<html>

</html>
`
)
