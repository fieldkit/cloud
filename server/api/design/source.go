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

var ClusterSummary = MediaType("application/vnd.app.geometry_cluster_summary+json", func() {
	TypeName("ClusterSummary")
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

var ClusterGeometrySummary = MediaType("application/vnd.app.cluster_geometry_summary+json", func() {
	TypeName("ClusterGeometrySummary")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("sourceId", Integer)
		Attribute("geometry", ArrayOf(ArrayOf(Number)))
		Required("id", "sourceId", "geometry")
	})
	View("default", func() {
		Attribute("id")
		Attribute("sourceId")
		Attribute("geometry")
	})
})

var SourceSummary = MediaType("application/vnd.app.source_summary+json", func() {
	TypeName("SourceSummary")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("name", String)
		Attribute("temporal", CollectionOf(ClusterSummary))
		Attribute("spatial", CollectionOf(ClusterSummary))
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

	Action("temporal cluster geometry by id", func() {
		NoSecurity()

		Routing(GET("sources/:sourceId/temporal/:clusterId/geometry"))
		Description("Retrieve temporal cluster geometry")
		Params(func() {
			Param("sourceId", Integer)
			Param("clusterId", Integer)
			Required("sourceId")
			Required("clusterId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(ClusterGeometrySummary)
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
