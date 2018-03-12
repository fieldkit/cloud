package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var Source = MediaType("application/vnd.app.source+json", func() {
	TypeName("Source")
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

var UpdateSourcePayload = Type("UpdateSourcePayload", func() {
	Reference(Source)
	Attribute("name")
	Attribute("teamId")
	Attribute("userId")
	Attribute("active")
	Required("name")
})

var Sources = MediaType("application/vnd.app.sources+json", func() {
	TypeName("Sources")
	Attributes(func() {
		Attribute("twitterAccountSources", CollectionOf(TwitterAccountSource))
		Attribute("deviceSources", CollectionOf(DeviceSource))
	})
	View("default", func() {
		Attribute("twitterAccountSources")
		Attribute("deviceSources")
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
		Attribute("geometry", ArrayOf(ArrayOf(Number)))
		Required("id", "startTime", "endTime", "numberOfFeatures", "centroid", "radius", "geometry")
	})
	View("default", func() {
		Attribute("id")
		Attribute("startTime")
		Attribute("endTime")
		Attribute("numberOfFeatures")
		Attribute("centroid")
		Attribute("radius")
		Attribute("geometry")
	})
})

var ReadingSummary = MediaType("application/vnd.app.reading_summary+json", func() {
	TypeName("ReadingSummary")
	Attributes(func() {
		Attribute("name", String)
		Required("name")
	})
	View("default", func() {
		Attribute("name")
	})
})

var SourceSummary = MediaType("application/vnd.app.source_summary+json", func() {
	TypeName("SourceSummary")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("name", String)
		Attribute("temporal", CollectionOf(GeometryClusterSummary))
		Attribute("spatial", CollectionOf(GeometryClusterSummary))
		Attribute("readings", CollectionOf(ReadingSummary))
		Required("id", "name", "temporal", "spatial")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("temporal")
		Attribute("spatial")
		Attribute("readings")
	})
})

var _ = Resource("source", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("update", func() {
		Routing(PATCH("sources/:sourceId"))
		Description("Update an source")
		Params(func() {
			Param("sourceId", Integer)
			Required("sourceId")
		})
		Payload(UpdateSourcePayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Source)
		})
	})

	Action("list", func() {
		NoSecurity()

		Routing(GET("projects/@/:project/expeditions/@/:expedition/sources"))
		Description("List a project's sources")
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("expedition", String, ExpeditionSlug)
			Required("project", "expedition")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Sources)
		})
	})

	Action("list id", func() {
		NoSecurity()

		Routing(GET("sources/:sourceId"))
		Description("List an source")
		Params(func() {
			Param("sourceId", Integer)
			Required("sourceId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(DeviceSource)
		})
	})

	Action("summary by id", func() {
		NoSecurity()

		Routing(GET("sources/:sourceId/summary"))
		Description("List an source")
		Params(func() {
			Param("sourceId", Integer)
			Required("sourceId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(SourceSummary)
		})
	})

	Action("list expedition id", func() {
		Routing(GET("expeditions/:expeditionId/sources"))
		Description("List an expedition's sources")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Sources)
		})
	})
})
