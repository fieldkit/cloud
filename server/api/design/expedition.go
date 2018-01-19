package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddExpeditionPayload = Type("AddExpeditionPayload", func() {
	Attribute("name", String)
	Attribute("slug", String, func() {
		Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
		MaxLength(40)
	})
	Attribute("description", String)
	Required("name", "slug", "description")
})

var Expedition = MediaType("application/vnd.app.expedition+json", func() {
	TypeName("Expedition")
	Reference(AddExpeditionPayload)
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("name")
		Attribute("slug")
		Attribute("description")
		Required("id", "name", "slug", "description")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("slug")
		Attribute("description")
	})
})

var Expeditions = MediaType("application/vnd.app.expeditions+json", func() {
	TypeName("Expeditions")
	Attributes(func() {
		Attribute("expeditions", CollectionOf(Expedition))
		Required("expeditions")
	})
	View("default", func() {
		Attribute("expeditions")
	})
})

var _ = Resource("expedition", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("projects/:projectId/expeditions"))
		Description("Add an expedition")
		Params(func() {
			Param("projectId", Integer)
		})
		Payload(AddExpeditionPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Expedition)
		})
	})

	Action("update", func() {
		Routing(PATCH("expeditions/:expeditionId"))
		Description("Update an expedition")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
		})
		Payload(AddExpeditionPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Expedition)
		})
	})

	Action("get", func() {
		NoSecurity()
		Routing(GET("projects/@/:project/expeditions/@/:expedition"))
		Description("Add an expedition")
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("expedition", String, ExpeditionSlug)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Expedition)
		})
	})

	Action("get id", func() {
		Routing(GET("expeditions/:expeditionId"))
		Description("Add an expedition")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Expedition)
		})
	})

	Action("list", func() {
		NoSecurity()
		Routing(GET("projects/@/:project/expeditions"))
		Description("List a project's expeditions")
		Params(func() {
			Param("project", String, ProjectSlug)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Expeditions)
		})
	})

	Action("list id", func() {
		Routing(GET("projects/:projectId/expeditions"))
		Description("List a project's expeditions")
		Params(func() {
			Param("projectId", Integer)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Expeditions)
		})
	})
})
