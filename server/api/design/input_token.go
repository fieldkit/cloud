package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var InputToken = MediaType("application/vnd.app.input_token+json", func() {
	TypeName("InputToken")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("token", String)
		Attribute("expedition_id", Integer)
		Required("id", "token", "expedition_id")
	})
	View("default", func() {
		Attribute("id")
		Attribute("token")
		Attribute("expedition_id")
	})
})

var InputTokens = MediaType("application/vnd.app.input_tokens+json", func() {
	TypeName("InputTokens")
	Attributes(func() {
		Attribute("input_tokens", CollectionOf(InputToken))
		Required("input_tokens")
	})
	View("default", func() {
		Attribute("input_tokens")
	})
})

var _ = Resource("input_token", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("expeditions/:expedition_id/input-tokens"))
		Description("Add an input token")
		Params(func() {
			Param("expedition_id", Integer)
			Required("expedition_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(InputToken)
		})
	})

	Action("delete", func() {
		Routing(DELETE("input-tokens/:input_token_id"))
		Description("Delete an input token")
		Params(func() {
			Param("input_token_id", Integer)
			Required("input_token_id")
		})
		Response(BadRequest)
		Response(NoContent)
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/input-tokens"))
		Description("List an expedition's input tokens")
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("expedition", String, ExpeditionSlug)
			Required("project", "expedition")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(InputTokens)
		})
	})

	Action("list id", func() {
		Routing(GET("expeditions/:expedition_id/input-tokens"))
		Description("Update an expedition's input tokens")
		Params(func() {
			Param("expedition_id", Integer)
			Required("expedition_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(InputTokens)
		})
	})
})
