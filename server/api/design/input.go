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

var GeometryClusterSummary = MediaType("application/vnd.app.geometry_cluster_summary+json", func() {
	TypeName("GeometryClusterSummary")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("startTime", DateTime)
		Attribute("endTime", DateTime)
		Attribute("numberOfFeatures", Integer)
		Attribute("centroid", ArrayOf(Number))
		Attribute("radius", Number)
		Required("id", "startTime", "endTime", "numberOfFeatures", "centroid", "radius")
	})
	View("default", func() {
		Attribute("id")
		Attribute("startTime")
		Attribute("endTime")
		Attribute("numberOfFeatures")
		Attribute("centroid")
		Attribute("radius")
	})
})

var InputSummary = MediaType("application/vnd.app.input_summary+json", func() {
	TypeName("InputSummary")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("name", String)
		Attribute("temporal", CollectionOf(GeometryClusterSummary))
		Attribute("spatial", CollectionOf(GeometryClusterSummary))
		Required("id", "name", "temporal", "spatial")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("temporal")
		Attribute("spatial")
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

	Action("summary by id", func() {
		NoSecurity()

		Routing(GET("inputs/:inputId/summary"))
		Description("List an input")
		Params(func() {
			Param("inputId", Integer)
			Required("inputId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(InputSummary)
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
