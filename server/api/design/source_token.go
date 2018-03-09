package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var SourceToken = MediaType("application/vnd.app.source_token+json", func() {
	TypeName("SourceToken")
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

var SourceTokens = MediaType("application/vnd.app.source_tokens+json", func() {
	TypeName("SourceTokens")
	Attributes(func() {
		Attribute("sourceTokens", CollectionOf(SourceToken))
		Required("sourceTokens")
	})
	View("default", func() {
		Attribute("sourceTokens")
	})
})

var _ = Resource("source_token", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("expeditions/:expeditionId/source-tokens"))
		Description("Add an source token")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(SourceToken)
		})
	})

	Action("delete", func() {
		Routing(DELETE("source-tokens/:sourceTokenId"))
		Description("Delete an source token")
		Params(func() {
			Param("sourceTokenId", Integer)
			Required("sourceTokenId")
		})
		Response(BadRequest)
		Response(NoContent)
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/source-tokens"))
		Description("List an expedition's source tokens")
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("expedition", String, ExpeditionSlug)
			Required("project", "expedition")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(SourceTokens)
		})
	})

	Action("list id", func() {
		Routing(GET("expeditions/:expeditionId/source-tokens"))
		Description("Update an expedition's source tokens")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(SourceTokens)
		})
	})
})
