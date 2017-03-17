// Code generated by goagen v1.1.0, command line:
// $ main
//
// API "fieldkit": Application User Types
//
// The content of this file is auto-generated, DO NOT MODIFY

package client

import (
	"github.com/goadesign/goa"
	"unicode/utf8"
)

// addExpeditionPayload user type.
type addExpeditionPayload struct {
	Description *string `form:"description,omitempty" json:"description,omitempty" xml:"description,omitempty"`
	Name        *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	Slug        *string `form:"slug,omitempty" json:"slug,omitempty" xml:"slug,omitempty"`
}

// Validate validates the addExpeditionPayload type instance.
func (ut *addExpeditionPayload) Validate() (err error) {
	if ut.Name == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if ut.Slug == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "slug"))
	}
	if ut.Description == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if ut.Slug != nil {
		if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, *ut.Slug); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`response.slug`, *ut.Slug, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
		}
	}
	if ut.Slug != nil {
		if utf8.RuneCountInString(*ut.Slug) > 40 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.slug`, *ut.Slug, utf8.RuneCountInString(*ut.Slug), 40, false))
		}
	}
	return
}

// Publicize creates AddExpeditionPayload from addExpeditionPayload
func (ut *addExpeditionPayload) Publicize() *AddExpeditionPayload {
	var pub AddExpeditionPayload
	if ut.Description != nil {
		pub.Description = *ut.Description
	}
	if ut.Name != nil {
		pub.Name = *ut.Name
	}
	if ut.Slug != nil {
		pub.Slug = *ut.Slug
	}
	return &pub
}

// AddExpeditionPayload user type.
type AddExpeditionPayload struct {
	Description string `form:"description" json:"description" xml:"description"`
	Name        string `form:"name" json:"name" xml:"name"`
	Slug        string `form:"slug" json:"slug" xml:"slug"`
}

// Validate validates the AddExpeditionPayload type instance.
func (ut *AddExpeditionPayload) Validate() (err error) {
	if ut.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if ut.Slug == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "slug"))
	}
	if ut.Description == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, ut.Slug); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.slug`, ut.Slug, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
	}
	if utf8.RuneCountInString(ut.Slug) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.slug`, ut.Slug, utf8.RuneCountInString(ut.Slug), 40, false))
	}
	return
}

// addProjectPayload user type.
type addProjectPayload struct {
	Description *string `form:"description,omitempty" json:"description,omitempty" xml:"description,omitempty"`
	Name        *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	Slug        *string `form:"slug,omitempty" json:"slug,omitempty" xml:"slug,omitempty"`
}

// Validate validates the addProjectPayload type instance.
func (ut *addProjectPayload) Validate() (err error) {
	if ut.Name == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if ut.Slug == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "slug"))
	}
	if ut.Description == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if ut.Slug != nil {
		if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, *ut.Slug); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`response.slug`, *ut.Slug, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
		}
	}
	if ut.Slug != nil {
		if utf8.RuneCountInString(*ut.Slug) > 40 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.slug`, *ut.Slug, utf8.RuneCountInString(*ut.Slug), 40, false))
		}
	}
	return
}

// Publicize creates AddProjectPayload from addProjectPayload
func (ut *addProjectPayload) Publicize() *AddProjectPayload {
	var pub AddProjectPayload
	if ut.Description != nil {
		pub.Description = *ut.Description
	}
	if ut.Name != nil {
		pub.Name = *ut.Name
	}
	if ut.Slug != nil {
		pub.Slug = *ut.Slug
	}
	return &pub
}

// AddProjectPayload user type.
type AddProjectPayload struct {
	Description string `form:"description" json:"description" xml:"description"`
	Name        string `form:"name" json:"name" xml:"name"`
	Slug        string `form:"slug" json:"slug" xml:"slug"`
}

// Validate validates the AddProjectPayload type instance.
func (ut *AddProjectPayload) Validate() (err error) {
	if ut.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if ut.Slug == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "slug"))
	}
	if ut.Description == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, ut.Slug); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.slug`, ut.Slug, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
	}
	if utf8.RuneCountInString(ut.Slug) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.slug`, ut.Slug, utf8.RuneCountInString(ut.Slug), 40, false))
	}
	return
}

// addTeamPayload user type.
type addTeamPayload struct {
	Description *string `form:"description,omitempty" json:"description,omitempty" xml:"description,omitempty"`
	Name        *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	Slug        *string `form:"slug,omitempty" json:"slug,omitempty" xml:"slug,omitempty"`
}

// Validate validates the addTeamPayload type instance.
func (ut *addTeamPayload) Validate() (err error) {
	if ut.Name == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if ut.Slug == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "slug"))
	}
	if ut.Description == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if ut.Slug != nil {
		if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, *ut.Slug); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`response.slug`, *ut.Slug, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
		}
	}
	if ut.Slug != nil {
		if utf8.RuneCountInString(*ut.Slug) > 40 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.slug`, *ut.Slug, utf8.RuneCountInString(*ut.Slug), 40, false))
		}
	}
	return
}

// Publicize creates AddTeamPayload from addTeamPayload
func (ut *addTeamPayload) Publicize() *AddTeamPayload {
	var pub AddTeamPayload
	if ut.Description != nil {
		pub.Description = *ut.Description
	}
	if ut.Name != nil {
		pub.Name = *ut.Name
	}
	if ut.Slug != nil {
		pub.Slug = *ut.Slug
	}
	return &pub
}

// AddTeamPayload user type.
type AddTeamPayload struct {
	Description string `form:"description" json:"description" xml:"description"`
	Name        string `form:"name" json:"name" xml:"name"`
	Slug        string `form:"slug" json:"slug" xml:"slug"`
}

// Validate validates the AddTeamPayload type instance.
func (ut *AddTeamPayload) Validate() (err error) {
	if ut.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if ut.Slug == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "slug"))
	}
	if ut.Description == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, ut.Slug); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.slug`, ut.Slug, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
	}
	if utf8.RuneCountInString(ut.Slug) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.slug`, ut.Slug, utf8.RuneCountInString(ut.Slug), 40, false))
	}
	return
}

// addUserPayload user type.
type addUserPayload struct {
	Email       *string `form:"email,omitempty" json:"email,omitempty" xml:"email,omitempty"`
	InviteToken *string `form:"invite_token,omitempty" json:"invite_token,omitempty" xml:"invite_token,omitempty"`
	Password    *string `form:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
	Username    *string `form:"username,omitempty" json:"username,omitempty" xml:"username,omitempty"`
}

// Validate validates the addUserPayload type instance.
func (ut *addUserPayload) Validate() (err error) {
	if ut.Email == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "email"))
	}
	if ut.Username == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "username"))
	}
	if ut.Password == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "password"))
	}
	if ut.InviteToken == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "invite_token"))
	}
	if ut.Email != nil {
		if err2 := goa.ValidateFormat(goa.FormatEmail, *ut.Email); err2 != nil {
			err = goa.MergeErrors(err, goa.InvalidFormatError(`response.email`, *ut.Email, goa.FormatEmail, err2))
		}
	}
	if ut.Password != nil {
		if utf8.RuneCountInString(*ut.Password) < 10 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.password`, *ut.Password, utf8.RuneCountInString(*ut.Password), 10, true))
		}
	}
	if ut.Username != nil {
		if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, *ut.Username); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`response.username`, *ut.Username, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
		}
	}
	if ut.Username != nil {
		if utf8.RuneCountInString(*ut.Username) > 40 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.username`, *ut.Username, utf8.RuneCountInString(*ut.Username), 40, false))
		}
	}
	return
}

// Publicize creates AddUserPayload from addUserPayload
func (ut *addUserPayload) Publicize() *AddUserPayload {
	var pub AddUserPayload
	if ut.Email != nil {
		pub.Email = *ut.Email
	}
	if ut.InviteToken != nil {
		pub.InviteToken = *ut.InviteToken
	}
	if ut.Password != nil {
		pub.Password = *ut.Password
	}
	if ut.Username != nil {
		pub.Username = *ut.Username
	}
	return &pub
}

// AddUserPayload user type.
type AddUserPayload struct {
	Email       string `form:"email" json:"email" xml:"email"`
	InviteToken string `form:"invite_token" json:"invite_token" xml:"invite_token"`
	Password    string `form:"password" json:"password" xml:"password"`
	Username    string `form:"username" json:"username" xml:"username"`
}

// Validate validates the AddUserPayload type instance.
func (ut *AddUserPayload) Validate() (err error) {
	if ut.Email == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "email"))
	}
	if ut.Username == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "username"))
	}
	if ut.Password == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "password"))
	}
	if ut.InviteToken == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "invite_token"))
	}
	if err2 := goa.ValidateFormat(goa.FormatEmail, ut.Email); err2 != nil {
		err = goa.MergeErrors(err, goa.InvalidFormatError(`response.email`, ut.Email, goa.FormatEmail, err2))
	}
	if utf8.RuneCountInString(ut.Password) < 10 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.password`, ut.Password, utf8.RuneCountInString(ut.Password), 10, true))
	}
	if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, ut.Username); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.username`, ut.Username, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
	}
	if utf8.RuneCountInString(ut.Username) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.username`, ut.Username, utf8.RuneCountInString(ut.Username), 40, false))
	}
	return
}

// loginPayload user type.
type loginPayload struct {
	Password *string `form:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
	Username *string `form:"username,omitempty" json:"username,omitempty" xml:"username,omitempty"`
}

// Validate validates the loginPayload type instance.
func (ut *loginPayload) Validate() (err error) {
	if ut.Username == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "username"))
	}
	if ut.Password == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "password"))
	}
	if ut.Password != nil {
		if utf8.RuneCountInString(*ut.Password) < 10 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.password`, *ut.Password, utf8.RuneCountInString(*ut.Password), 10, true))
		}
	}
	if ut.Username != nil {
		if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, *ut.Username); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`response.username`, *ut.Username, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
		}
	}
	if ut.Username != nil {
		if utf8.RuneCountInString(*ut.Username) > 40 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`response.username`, *ut.Username, utf8.RuneCountInString(*ut.Username), 40, false))
		}
	}
	return
}

// Publicize creates LoginPayload from loginPayload
func (ut *loginPayload) Publicize() *LoginPayload {
	var pub LoginPayload
	if ut.Password != nil {
		pub.Password = *ut.Password
	}
	if ut.Username != nil {
		pub.Username = *ut.Username
	}
	return &pub
}

// LoginPayload user type.
type LoginPayload struct {
	Password string `form:"password" json:"password" xml:"password"`
	Username string `form:"username" json:"username" xml:"username"`
}

// Validate validates the LoginPayload type instance.
func (ut *LoginPayload) Validate() (err error) {
	if ut.Username == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "username"))
	}
	if ut.Password == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "password"))
	}
	if utf8.RuneCountInString(ut.Password) < 10 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.password`, ut.Password, utf8.RuneCountInString(ut.Password), 10, true))
	}
	if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, ut.Username); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.username`, ut.Username, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
	}
	if utf8.RuneCountInString(ut.Username) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.username`, ut.Username, utf8.RuneCountInString(ut.Username), 40, false))
	}
	return
}
