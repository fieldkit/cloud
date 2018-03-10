package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var SeriesData = MediaType("application/vnd.app.series+json", func() {
	TypeName("SeriesData")
	Attributes(func() {
		Attribute("name", String)
		Attribute("rows", ArrayOf(Any))
		Required("name")
		Required("rows")
	})
	View("default", func() {
		Attribute("name")
		Attribute("rows")
	})
})

var QueryData = MediaType("application/vnd.app.queried+json", func() {
	TypeName("QueryData")
	Attributes(func() {
		Attribute("series", CollectionOf(SeriesData))
		Required("series")
	})
	View("default", func() {
		Attribute("series")
	})
})

var _ = Resource("Query", func() {
	Action("list by source", func() {
		Routing(GET("sources/:sourceId/query"))
		Description("Query data for an source.")
		Params(func() {
			Param("sourceId", Integer)
			Required("sourceId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(QueryData)
		})
	})
})
