package email

import (
	"fmt"
)

type Emailer interface {
	SendEmail(address, subject, message string) error
}

type emailer struct{}

func (e emailer) SendEmail(address, subject, message string) error {
	fmt.Printf("To: %s\nSubject: %s\n\n%s\n", address, subject, message)
	return nil
}

func NewEmailer() Emailer {
	return emailer{}
}
