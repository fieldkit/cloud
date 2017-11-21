package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var GeoJSONGeometry = MediaType("application/vnd.app.geojson-geometry+json", func() {
	TypeName("GeoJSONGeometry")
	Attributes(func() {
		Attribute("type", String)
		Attribute("coordinates", ArrayOf(Number))
		Required("type")
		Required("coordinates")
	})
	View("default", func() {
		Attribute("type")
		Attribute("coordinates")
	})
})

var GeoJSONFeature = MediaType("application/vnd.app.geojson-feature+json", func() {
	TypeName("GeoJSONFeature")
	Attributes(func() {
		Attribute("type", String)
		Attribute("geometry", GeoJSONGeometry)
		Attribute("properties", HashOf(String, Any))
		Required("type")
		Required("geometry")
		Required("properties")
	})
	View("default", func() {
		Attribute("type")
		Attribute("geometry")
		Attribute("properties")
	})
})

var GeoJSON = MediaType("application/vnd.app.geojson+json", func() {
	TypeName("GeoJSON")
	Attributes(func() {
		Attribute("type", String)
		Attribute("features", CollectionOf(GeoJSONFeature))
		Required("type")
		Required("features")
	})
	View("default", func() {
		Attribute("type")
		Attribute("features")
	})
})

var PagedGeoJSON = MediaType("application/vnd.app.paged-geojson+json", func() {
	TypeName("PagedGeoJSON")
	Attributes(func() {
		Attribute("nextUrl", String)
		Attribute("previousUrl", String)
		Attribute("geo", GeoJSON)
		Attribute("hasMore", Boolean)
		Required("nextUrl", "geo", "hasMore")
	})
	View("default", func() {
		Attribute("nextUrl")
		Attribute("previousUrl")
		Attribute("hasMore")
		Attribute("geo")
	})
})

var _ = Resource("GeoJSON", func() {
	Action("list", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/geojson"))
		Description("List a expedition's features in a GeoJSON.")
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("expedition", String, ExpeditionSlug)
			Required("project", "expedition")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(PagedGeoJSON)
		})
	})
})
