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
		Attribute("expeditionId", Integer)
		Required("id", "token", "expeditionId")
	})
	View("default", func() {
		Attribute("id")
		Attribute("token")
		Attribute("expeditionId")
	})
})

var InputTokens = MediaType("application/vnd.app.input_tokens+json", func() {
	TypeName("InputTokens")
	Attributes(func() {
		Attribute("inputTokens", CollectionOf(InputToken))
		Required("inputTokens")
	})
	View("default", func() {
		Attribute("inputTokens")
	})
})

var _ = Resource("input_token", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("expeditions/:expeditionId/input-tokens"))
		Description("Add an input token")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(InputToken)
		})
	})

	Action("delete", func() {
		Routing(DELETE("input-tokens/:inputTokenId"))
		Description("Delete an input token")
		Params(func() {
			Param("inputTokenId", Integer)
			Required("inputTokenId")
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
		Routing(GET("expeditions/:expeditionId/input-tokens"))
		Description("Update an expedition's input tokens")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(InputTokens)
		})
	})
})
