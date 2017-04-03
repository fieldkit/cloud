// Code generated by goagen v1.1.0, command line:
// $ main
//
// API "fieldkit": Application Media Types
//
// The content of this file is auto-generated, DO NOT MODIFY

package app

import (
	"github.com/goadesign/goa"
	"unicode/utf8"
)

// ProjectAdministrator media type (default view)
//
// Identifier: application/vnd.app.administrator+json; view=default
type ProjectAdministrator struct {
	ProjectID int `form:"project_id" json:"project_id" xml:"project_id"`
	UserID    int `form:"user_id" json:"user_id" xml:"user_id"`
}

// Validate validates the ProjectAdministrator media type instance.
func (mt *ProjectAdministrator) Validate() (err error) {

	return
}

// ProjectAdministratorCollection is the media type for an array of ProjectAdministrator (default view)
//
// Identifier: application/vnd.app.administrator+json; type=collection; view=default
type ProjectAdministratorCollection []*ProjectAdministrator

// Validate validates the ProjectAdministratorCollection media type instance.
func (mt ProjectAdministratorCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ProjectAdministrators media type (default view)
//
// Identifier: application/vnd.app.administrators+json; view=default
type ProjectAdministrators struct {
	Administrators ProjectAdministratorCollection `form:"administrators" json:"administrators" xml:"administrators"`
}

// Validate validates the ProjectAdministrators media type instance.
func (mt *ProjectAdministrators) Validate() (err error) {
	if mt.Administrators == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "administrators"))
	}
	return
}

// Document media type (default view)
//
// Identifier: application/vnd.app.document+json; view=default
type Document struct {
	Data      interface{} `form:"data" json:"data" xml:"data"`
	ID        string      `form:"id" json:"id" xml:"id"`
	InputID   int         `form:"input_id" json:"input_id" xml:"input_id"`
	Location  *Point      `form:"location" json:"location" xml:"location"`
	TeamID    *int        `form:"team_id,omitempty" json:"team_id,omitempty" xml:"team_id,omitempty"`
	Timestamp int         `form:"timestamp" json:"timestamp" xml:"timestamp"`
	UserID    *int        `form:"user_id,omitempty" json:"user_id,omitempty" xml:"user_id,omitempty"`
}

// Validate validates the Document media type instance.
func (mt *Document) Validate() (err error) {
	if mt.ID == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "id"))
	}

	if mt.Location == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "location"))
	}

	if mt.Location != nil {
		if err2 := mt.Location.Validate(); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// DocumentCollection is the media type for an array of Document (default view)
//
// Identifier: application/vnd.app.document+json; type=collection; view=default
type DocumentCollection []*Document

// Validate validates the DocumentCollection media type instance.
func (mt DocumentCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// Documents media type (default view)
//
// Identifier: application/vnd.app.documents+json; view=default
type Documents struct {
	Documents DocumentCollection `form:"documents" json:"documents" xml:"documents"`
}

// Validate validates the Documents media type instance.
func (mt *Documents) Validate() (err error) {
	if mt.Documents == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "documents"))
	}
	if err2 := mt.Documents.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

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

// FieldkitInput media type (default view)
//
// Identifier: application/vnd.app.fieldkit_input+json; view=default
type FieldkitInput struct {
	ExpeditionID int    `form:"expedition_id" json:"expedition_id" xml:"expedition_id"`
	ID           int    `form:"id" json:"id" xml:"id"`
	Name         string `form:"name" json:"name" xml:"name"`
	TeamID       *int   `form:"team_id,omitempty" json:"team_id,omitempty" xml:"team_id,omitempty"`
	UserID       *int   `form:"user_id,omitempty" json:"user_id,omitempty" xml:"user_id,omitempty"`
}

// Validate validates the FieldkitInput media type instance.
func (mt *FieldkitInput) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	return
}

// FieldkitInputCollection is the media type for an array of FieldkitInput (default view)
//
// Identifier: application/vnd.app.fieldkit_input+json; type=collection; view=default
type FieldkitInputCollection []*FieldkitInput

// Validate validates the FieldkitInputCollection media type instance.
func (mt FieldkitInputCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// FieldkitBinary media type (default view)
//
// Identifier: application/vnd.app.fieldkit_input_binary+json; view=default
type FieldkitBinary struct {
	Fields  []string `form:"fields" json:"fields" xml:"fields"`
	ID      int      `form:"id" json:"id" xml:"id"`
	InputID int      `form:"input_id" json:"input_id" xml:"input_id"`
}

// Validate validates the FieldkitBinary media type instance.
func (mt *FieldkitBinary) Validate() (err error) {

	if mt.Fields == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "fields"))
	}
	for _, e := range mt.Fields {
		if !(e == "varint" || e == "uvarint" || e == "float32" || e == "float64") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError(`response.fields[*]`, e, []interface{}{"varint", "uvarint", "float32", "float64"}))
		}
	}
	return
}

// FieldkitInputs media type (default view)
//
// Identifier: application/vnd.app.fieldkit_inputs+json; view=default
type FieldkitInputs struct {
	FieldkitInputs FieldkitInputCollection `form:"fieldkit_inputs" json:"fieldkit_inputs" xml:"fieldkit_inputs"`
}

// Validate validates the FieldkitInputs media type instance.
func (mt *FieldkitInputs) Validate() (err error) {
	if mt.FieldkitInputs == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "fieldkit_inputs"))
	}
	if err2 := mt.FieldkitInputs.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// Input media type (default view)
//
// Identifier: application/vnd.app.input+json; view=default
type Input struct {
	ExpeditionID int    `form:"expedition_id" json:"expedition_id" xml:"expedition_id"`
	ID           int    `form:"id" json:"id" xml:"id"`
	Name         string `form:"name" json:"name" xml:"name"`
	TeamID       *int   `form:"team_id,omitempty" json:"team_id,omitempty" xml:"team_id,omitempty"`
	UserID       *int   `form:"user_id,omitempty" json:"user_id,omitempty" xml:"user_id,omitempty"`
}

// Validate validates the Input media type instance.
func (mt *Input) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	return
}

// Inputs media type (default view)
//
// Identifier: application/vnd.app.inputs+json; view=default
type Inputs struct {
	FieldkitInputs       FieldkitInputCollection       `form:"fieldkit_inputs,omitempty" json:"fieldkit_inputs,omitempty" xml:"fieldkit_inputs,omitempty"`
	TwitterAccountInputs TwitterAccountInputCollection `form:"twitter_account_inputs,omitempty" json:"twitter_account_inputs,omitempty" xml:"twitter_account_inputs,omitempty"`
}

// Validate validates the Inputs media type instance.
func (mt *Inputs) Validate() (err error) {
	if err2 := mt.FieldkitInputs.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	if err2 := mt.TwitterAccountInputs.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// Location media type (default view)
//
// Identifier: application/vnd.app.location+json; view=default
type Location struct {
	Location string `form:"location" json:"location" xml:"location"`
}

// Validate validates the Location media type instance.
func (mt *Location) Validate() (err error) {
	if mt.Location == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "location"))
	}
	if err2 := goa.ValidateFormat(goa.FormatURI, mt.Location); err2 != nil {
		err = goa.MergeErrors(err, goa.InvalidFormatError(`response.location`, mt.Location, goa.FormatURI, err2))
	}
	return
}

// TeamMember media type (default view)
//
// Identifier: application/vnd.app.member+json; view=default
type TeamMember struct {
	Role   string `form:"role" json:"role" xml:"role"`
	TeamID int    `form:"team_id" json:"team_id" xml:"team_id"`
	UserID int    `form:"user_id" json:"user_id" xml:"user_id"`
}

// Validate validates the TeamMember media type instance.
func (mt *TeamMember) Validate() (err error) {

	if mt.Role == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "role"))
	}
	return
}

// TeamMemberCollection is the media type for an array of TeamMember (default view)
//
// Identifier: application/vnd.app.member+json; type=collection; view=default
type TeamMemberCollection []*TeamMember

// Validate validates the TeamMemberCollection media type instance.
func (mt TeamMemberCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// TeamMembers media type (default view)
//
// Identifier: application/vnd.app.members+json; view=default
type TeamMembers struct {
	Members TeamMemberCollection `form:"members" json:"members" xml:"members"`
}

// Validate validates the TeamMembers media type instance.
func (mt *TeamMembers) Validate() (err error) {
	if mt.Members == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "members"))
	}
	if err2 := mt.Members.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// Point media type (default view)
//
// Identifier: application/vnd.app.point+json; view=default
type Point struct {
	Coordinates []float64 `form:"coordinates" json:"coordinates" xml:"coordinates"`
	Type        string    `form:"type" json:"type" xml:"type"`
}

// Validate validates the Point media type instance.
func (mt *Point) Validate() (err error) {
	if mt.Type == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "type"))
	}
	if mt.Coordinates == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "coordinates"))
	}
	if len(mt.Coordinates) < 2 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.coordinates`, mt.Coordinates, len(mt.Coordinates), 2, true))
	}
	if len(mt.Coordinates) > 2 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.coordinates`, mt.Coordinates, len(mt.Coordinates), 2, false))
	}
	if !(mt.Type == "Point") {
		err = goa.MergeErrors(err, goa.InvalidEnumValueError(`response.type`, mt.Type, []interface{}{"Point"}))
	}
	return
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

// Schema media type (default view)
//
// Identifier: application/vnd.app.schema+json; view=default
type Schema struct {
	ID         int         `form:"id" json:"id" xml:"id"`
	JSONSchema interface{} `form:"json_schema" json:"json_schema" xml:"json_schema"`
	ProjectID  *int        `form:"project_id,omitempty" json:"project_id,omitempty" xml:"project_id,omitempty"`
}

// Validate validates the Schema media type instance.
func (mt *Schema) Validate() (err error) {

	return
}

// SchemaCollection is the media type for an array of Schema (default view)
//
// Identifier: application/vnd.app.schema+json; type=collection; view=default
type SchemaCollection []*Schema

// Validate validates the SchemaCollection media type instance.
func (mt SchemaCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// Schemas media type (default view)
//
// Identifier: application/vnd.app.schemas+json; view=default
type Schemas struct {
	Schemas SchemaCollection `form:"schemas,omitempty" json:"schemas,omitempty" xml:"schemas,omitempty"`
}

// Team media type (default view)
//
// Identifier: application/vnd.app.team+json; view=default
type Team struct {
	Description string `form:"description" json:"description" xml:"description"`
	ID          int    `form:"id" json:"id" xml:"id"`
	Name        string `form:"name" json:"name" xml:"name"`
	Slug        string `form:"slug" json:"slug" xml:"slug"`
}

// Validate validates the Team media type instance.
func (mt *Team) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Slug == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "slug"))
	}
	if mt.Description == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "description"))
	}
	if ok := goa.ValidatePattern(`\S`, mt.Name); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.name`, mt.Name, `\S`))
	}
	if utf8.RuneCountInString(mt.Name) > 256 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, mt.Name, utf8.RuneCountInString(mt.Name), 256, false))
	}
	if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, mt.Slug); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.slug`, mt.Slug, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
	}
	if utf8.RuneCountInString(mt.Slug) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.slug`, mt.Slug, utf8.RuneCountInString(mt.Slug), 40, false))
	}
	return
}

// TeamCollection is the media type for an array of Team (default view)
//
// Identifier: application/vnd.app.team+json; type=collection; view=default
type TeamCollection []*Team

// Validate validates the TeamCollection media type instance.
func (mt TeamCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// Teams media type (default view)
//
// Identifier: application/vnd.app.teams+json; view=default
type Teams struct {
	Teams TeamCollection `form:"teams" json:"teams" xml:"teams"`
}

// Validate validates the Teams media type instance.
func (mt *Teams) Validate() (err error) {
	if mt.Teams == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "teams"))
	}
	if err2 := mt.Teams.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// TwitterAccountInput media type (default view)
//
// Identifier: application/vnd.app.twitter_account_input+json; view=default
type TwitterAccountInput struct {
	ExpeditionID     int    `form:"expedition_id" json:"expedition_id" xml:"expedition_id"`
	ID               int    `form:"id" json:"id" xml:"id"`
	Name             string `form:"name" json:"name" xml:"name"`
	ScreenName       string `form:"screen_name" json:"screen_name" xml:"screen_name"`
	TeamID           *int   `form:"team_id,omitempty" json:"team_id,omitempty" xml:"team_id,omitempty"`
	TwitterAccountID int    `form:"twitter_account_id" json:"twitter_account_id" xml:"twitter_account_id"`
	UserID           *int   `form:"user_id,omitempty" json:"user_id,omitempty" xml:"user_id,omitempty"`
}

// Validate validates the TwitterAccountInput media type instance.
func (mt *TwitterAccountInput) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}

	if mt.ScreenName == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "screen_name"))
	}
	return
}

// TwitterAccountInputCollection is the media type for an array of TwitterAccountInput (default view)
//
// Identifier: application/vnd.app.twitter_account_input+json; type=collection; view=default
type TwitterAccountInputCollection []*TwitterAccountInput

// Validate validates the TwitterAccountInputCollection media type instance.
func (mt TwitterAccountInputCollection) Validate() (err error) {
	for _, e := range mt {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// TwitterAccountInputs media type (default view)
//
// Identifier: application/vnd.app.twitter_account_intputs+json; view=default
type TwitterAccountInputs struct {
	TwitterAccountInputs TwitterAccountInputCollection `form:"twitter_account_inputs" json:"twitter_account_inputs" xml:"twitter_account_inputs"`
}

// Validate validates the TwitterAccountInputs media type instance.
func (mt *TwitterAccountInputs) Validate() (err error) {
	if mt.TwitterAccountInputs == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "twitter_account_inputs"))
	}
	if err2 := mt.TwitterAccountInputs.Validate(); err2 != nil {
		err = goa.MergeErrors(err, err2)
	}
	return
}

// User media type (default view)
//
// Identifier: application/vnd.app.user+json; view=default
type User struct {
	Bio      string `form:"bio" json:"bio" xml:"bio"`
	Email    string `form:"email" json:"email" xml:"email"`
	ID       int    `form:"id" json:"id" xml:"id"`
	Name     string `form:"name" json:"name" xml:"name"`
	Username string `form:"username" json:"username" xml:"username"`
}

// Validate validates the User media type instance.
func (mt *User) Validate() (err error) {

	if mt.Name == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "name"))
	}
	if mt.Username == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "username"))
	}
	if mt.Email == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "email"))
	}
	if mt.Bio == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "bio"))
	}
	if err2 := goa.ValidateFormat(goa.FormatEmail, mt.Email); err2 != nil {
		err = goa.MergeErrors(err, goa.InvalidFormatError(`response.email`, mt.Email, goa.FormatEmail, err2))
	}
	if ok := goa.ValidatePattern(`\S`, mt.Name); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.name`, mt.Name, `\S`))
	}
	if utf8.RuneCountInString(mt.Name) > 256 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.name`, mt.Name, utf8.RuneCountInString(mt.Name), 256, false))
	}
	if ok := goa.ValidatePattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`, mt.Username); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`response.username`, mt.Username, `^[[:alnum:]]+(-[[:alnum:]]+)*$`))
	}
	if utf8.RuneCountInString(mt.Username) > 40 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`response.username`, mt.Username, utf8.RuneCountInString(mt.Username), 40, false))
	}
	return
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
