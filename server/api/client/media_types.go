// Code generated by goagen v1.1.0, command line:
// $ main
//
// API "fieldkit": Application Media Types
//
// The content of this file is auto-generated, DO NOT MODIFY

package client

import (
	"github.com/goadesign/goa"
	"net/http"
	"unicode/utf8"
)

// Expedition media type (default view)
//
// Identifier: application/vnd.app.expedition+json; view=default
type Expedition struct {
	Description string `form:"description" json:"description" xml:"description"`
	ID          int    `form:"id" json:"id" xml:"id"`
	Name        string `form:"name" json:"name" xml:"name"`
	Slug        string `form:"slug" json:"slug" xml:"slug"`
}

// Validate validates the Expedition media type instance.
func (mt *Expedition) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Slug == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "slug"))
	}
	if mt.Description == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, mt.Slug); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.slug`, mt.Slug, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
	}
	if utf8.RuneCountInString(mt.Slug) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.slug`, mt.Slug, utf8.RuneCountInString(mt.Slug), 40, false))
	}
	return
}

// DecodeExpedition decodes the Expedition instance encoded in resp body.
func (c *Client) DecodeExpedition(resp *http.Response) (*Expedition, error) {
	var decoded Expedition
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// ExpeditionCollection is the media type for an array of Expedition (default view)
//
// Identifier: application/vnd.app.expedition+json; type=collection; view=default
type ExpeditionCollection []*Expedition

// Validate validates the ExpeditionCollection media type instance.
func (mt ExpeditionCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// DecodeExpeditionCollection decodes the ExpeditionCollection instance encoded in resp body.
func (c *Client) DecodeExpeditionCollection(resp *http.Response) (ExpeditionCollection, error) {
	var decoded ExpeditionCollection
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return decoded, err
}

// Expeditions media type (default view)
//
// Identifier: application/vnd.app.expeditions+json; view=default
type Expeditions struct {
	Expeditions ExpeditionCollection `form:"expeditions" json:"expeditions" xml:"expeditions"`
}

// Validate validates the Expeditions media type instance.
func (mt *Expeditions) Validate() (err error) {
	if mt.Expeditions == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "expeditions"))
	}
	if err2 := mt.Expeditions.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// DecodeExpeditions decodes the Expeditions instance encoded in resp body.
func (c *Client) DecodeExpeditions(resp *http.Response) (*Expeditions, error) {
	var decoded Expeditions
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// Project media type (default view)
//
// Identifier: application/vnd.app.project+json; view=default
type Project struct {
	Description string `form:"description" json:"description" xml:"description"`
	ID          int    `form:"id" json:"id" xml:"id"`
	Name        string `form:"name" json:"name" xml:"name"`
	Slug        string `form:"slug" json:"slug" xml:"slug"`
}

// Validate validates the Project media type instance.
func (mt *Project) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Slug == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "slug"))
	}
	if mt.Description == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, mt.Slug); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.slug`, mt.Slug, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
	}
	if utf8.RuneCountInString(mt.Slug) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.slug`, mt.Slug, utf8.RuneCountInString(mt.Slug), 40, false))
	}
	return
}

// DecodeProject decodes the Project instance encoded in resp body.
func (c *Client) DecodeProject(resp *http.Response) (*Project, error) {
	var decoded Project
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// ProjectCollection is the media type for an array of Project (default view)
//
// Identifier: application/vnd.app.project+json; type=collection; view=default
type ProjectCollection []*Project

// Validate validates the ProjectCollection media type instance.
func (mt ProjectCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// DecodeProjectCollection decodes the ProjectCollection instance encoded in resp body.
func (c *Client) DecodeProjectCollection(resp *http.Response) (ProjectCollection, error) {
	var decoded ProjectCollection
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return decoded, err
}

// Projects media type (default view)
//
// Identifier: application/vnd.app.projects+json; view=default
type Projects struct {
	Projects ProjectCollection `form:"projects" json:"projects" xml:"projects"`
}

// Validate validates the Projects media type instance.
func (mt *Projects) Validate() (err error) {
	if mt.Projects == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "projects"))
	}
	if err2 := mt.Projects.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// DecodeProjects decodes the Projects instance encoded in resp body.
func (c *Client) DecodeProjects(resp *http.Response) (*Projects, error) {
	var decoded Projects
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// User media type (default view)
//
// Identifier: application/vnd.app.user+json; view=default
type User struct {
	ID       int    `form:"id" json:"id" xml:"id"`
	Username string `form:"username" json:"username" xml:"username"`
}

// Validate validates the User media type instance.
func (mt *User) Validate() (err error) {

	if mt.Username == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "username"))
	}
	if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, mt.Username); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.username`, mt.Username, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
	}
	if utf8.RuneCountInString(mt.Username) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.username`, mt.Username, utf8.RuneCountInString(mt.Username), 40, false))
	}
	return
}

// DecodeUser decodes the User instance encoded in resp body.
func (c *Client) DecodeUser(resp *http.Response) (*User, error) {
	var decoded User
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}

// UserCollection is the media type for an array of User (default view)
//
// Identifier: application/vnd.app.user+json; type=collection; view=default
type UserCollection []*User

// Validate validates the UserCollection media type instance.
func (mt UserCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// DecodeUserCollection decodes the UserCollection instance encoded in resp body.
func (c *Client) DecodeUserCollection(resp *http.Response) (UserCollection, error) {
	var decoded UserCollection
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return decoded, err
}

// Users media type (default view)
//
// Identifier: application/vnd.app.users+json; view=default
type Users struct {
	Users UserCollection `form:"users" json:"users" xml:"users"`
}

// Validate validates the Users media type instance.
func (mt *Users) Validate() (err error) {
	if mt.Users == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "users"))
	}
	if err2 := mt.Users.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// DecodeUsers decodes the Users instance encoded in resp body.
func (c *Client) DecodeUsers(resp *http.Response) (*Users, error) {
	var decoded Users
	err := c.Decoder.Decode(&decoded, resp.Body, resp.Header.Get("Content-Type"))
	return &decoded, err
}
