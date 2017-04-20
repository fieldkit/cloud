package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var Point = MediaType("application/vnd.app.point+json", func() {
	TypeName("Point")
	Attributes(func() {
		Attribute("type", String, func() {
			Enum("Point")
		})
		Attribute("coordinates", ArrayOf(Number), func() {
			MaxLength(2)
			MinLength(2)
		})
		Required("type", "coordinates")
	})
	View("default", func() {
		Attribute("type")
		Attribute("coordinates")
	})
})

var Document = MediaType("application/vnd.app.document+json", func() {
	TypeName("Document")
	Attributes(func() {
		Attribute("id", String)
		Attribute("input_id", Integer)
		Attribute("team_id", Integer)
		Attribute("user_id", Integer)
		Attribute("timestamp", Integer)
		Attribute("location", Point)
		Attribute("data", Any)
		Required("id", "input_id", "timestamp", "location", "data")
	})
	View("default", func() {
		Attribute("id")
		Attribute("input_id")
		Attribute("team_id")
		Attribute("user_id")
		Attribute("timestamp")
		Attribute("location")
		Attribute("data")
	})
})

var Documents = MediaType("application/vnd.app.documents+json", func() {
	TypeName("Documents")
	Attributes(func() {
		Attribute("documents", CollectionOf(Document))
		Required("documents")
	})
	View("default", func() {
		Attribute("documents")
	})
})

var _ = Resource("document", func() {
	Action("list", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/documents"))
		Description("List a expedition's documents")
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("expedition", String, ExpeditionSlug)
			Required("project", "expedition")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Documents)
		})
	})

	Action("list id", func() {
		Routing(GET("expeditions/:expedition_id/documents"))
		Description("List a expedition's documents")
		Params(func() {
			Param("expedition_id", Integer)
			Required("expedition_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Documents)
		})
	})
})
