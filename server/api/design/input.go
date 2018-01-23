package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var Input = MediaType("application/vnd.app.input+json", func() {
	TypeName("Input")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("expeditionId", Integer)
		Attribute("name", String)
		Attribute("teamId", Integer)
		Attribute("userId", Integer)
		Attribute("active", Boolean)
		Required("id", "expeditionId", "name")
	})
	View("default", func() {
		Attribute("id")
		Attribute("expeditionId")
		Attribute("name")
		Attribute("teamId")
		Attribute("userId")
		Attribute("active")
	})
})

var UpdateInputPayload = Type("UpdateInputPayload", func() {
	Reference(Input)
	Attribute("name")
	Attribute("teamId")
	Attribute("userId")
	Attribute("active")
	Required("name")
})

var Inputs = MediaType("application/vnd.app.inputs+json", func() {
	TypeName("Inputs")
	Attributes(func() {
		Attribute("twitterAccountInputs", CollectionOf(TwitterAccountInput))
		Attribute("deviceInputs", CollectionOf(DeviceInput))
	})
	View("default", func() {
		Attribute("twitterAccountInputs")
		Attribute("deviceInputs")
	})
})

var _ = Resource("input", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("update", func() {
		Routing(PATCH("inputs/:inputId"))
		Description("Update an input")
		Params(func() {
			Param("inputId", Integer)
			Required("inputId")
		})
		Payload(UpdateInputPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Input)
		})
	})

	Action("list", func() {
		NoSecurity()

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
		NoSecurity()

		Routing(GET("inputs/:inputId"))
		Description("List an input")
		Params(func() {
			Param("inputId", Integer)
			Required("inputId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(DeviceInput)
		})
	})

	Action("list expedition id", func() {
		Routing(GET("expeditions/:expeditionId/inputs"))
		Description("List an expedition's inputs")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Inputs)
		})
	})
})
