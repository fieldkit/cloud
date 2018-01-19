package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var Input = MediaType("application/vnd.app.input+json", func() {
	TypeName("Input")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("expedition_id", Integer)
		Attribute("name", String)
		Attribute("team_id", Integer)
		Attribute("user_id", Integer)
		Attribute("active", Boolean)
		Required("id", "expedition_id", "name")
	})
	View("default", func() {
		Attribute("id")
		Attribute("expedition_id")
		Attribute("name")
		Attribute("team_id")
		Attribute("user_id")
		Attribute("active")
	})
})

var UpdateInputPayload = Type("UpdateInputPayload", func() {
	Reference(Input)
	Attribute("name")
	Attribute("team_id")
	Attribute("user_id")
	Attribute("active")
	Required("name")
})

var Inputs = MediaType("application/vnd.app.inputs+json", func() {
	TypeName("Inputs")
	Attributes(func() {
		Attribute("twitter_account_inputs", CollectionOf(TwitterAccountInput))
		Attribute("device_inputs", CollectionOf(DeviceInput))
	})
	View("default", func() {
		Attribute("twitter_account_inputs")
		Attribute("device_inputs")
	})
})

var _ = Resource("input", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("update", func() {
		Routing(PATCH("inputs/:input_id"))
		Description("Update an input")
		Params(func() {
			Param("input_id", Integer)
			Required("input_id")
		})
		Payload(UpdateInputPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Input)
		})
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/inputs"))
		Description("List a project's inputs")
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("expedition", String, ExpeditionSlug)
			Required("project", "expedition")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Inputs)
		})
	})

	Action("list id", func() {
		NoSecurity() // TOOD: Fix this.

		Routing(GET("inputs/:input_id"))
		Description("List an input")
		Params(func() {
			Param("input_id", Integer)
			Required("input_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Inputs)
		})
	})

	Action("list expedition id", func() {
		Routing(GET("expeditions/:expedition_id/inputs"))
		Description("List an expedition's inputs")
		Params(func() {
			Param("expedition_id", Integer)
			Required("expedition_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Inputs)
		})
	})
})
