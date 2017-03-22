package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var Inputs = MediaType("application/vnd.app.inputs+json", func() {
	TypeName("Inputs")
	Attributes(func() {
		Attribute("twitter_accounts", CollectionOf(TwitterAccount))
	})
	View("default", func() {
		Attribute("twitter_accounts")
	})
})

var _ = Resource("input", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/inputs"))
		Description("List a project's inputs")
		Params(func() {
			Param("project", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Project slug")
			})
			Param("expedition", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Expedition slug")
			})
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Inputs)
		})
	})

	Action("list id", func() {
		Routing(GET("expeditions/:expedition_id/inputs"))
		Description("List a project's inputs")
		Params(func() {
			Param("expedition_id", Integer)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Inputs)
		})
	})
})
