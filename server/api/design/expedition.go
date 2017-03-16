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
		Routing(POST("project/:project/expedition"))
		Description("Add a expedition")
		Payload(AddExpeditionPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Expedition)
		})
	})

	// Action("get", func() {
	// 	Routing(GET("expeditions/:expedition"))
	// 	Description("Get a expedition")
	// 	Response(BadRequest)
	// 	Response(OK, func() {
	// 		Media(Expedition)
	// 	})
	// })

	// Action("list", func() {
	// 	Routing(GET("expeditions"))
	// 	Description("List expeditions")
	// 	Response(BadRequest)
	// 	Response(OK, func() {
	// 		Media(Expeditions)
	// 	})
	// })

	Action("list project", func() {
		Routing(GET("project/:project/expeditions"))
		Description("List a project's expeditions")
		Response(BadRequest)
		Response(OK, func() {
			Media(Expeditions)
		})
	})
})
