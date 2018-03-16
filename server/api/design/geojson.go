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

var PagedGeoJSON = MediaType("application/vnd.app.paged_geojson+json", func() {
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

var MapFeatures = MediaType("application/vnd.app.map_features+json", func() {
	TypeName("MapFeatures")
	Attributes(func() {
		Attribute("geoJSON", PagedGeoJSON)
		Attribute("temporal", CollectionOf(ClusterSummary))
		Attribute("spatial", CollectionOf(ClusterSummary))
		Attribute("readings", CollectionOf(ReadingSummary))
		Attribute("geometries", CollectionOf(ClusterGeometrySummary))
		Required("geoJSON", "temporal", "spatial", "readings", "geometries")
	})
	View("default", func() {
		Attribute("geoJSON")
		Attribute("temporal")
		Attribute("spatial")
		Attribute("readings")
		Attribute("geometries")
	})
})

var _ = Resource("GeoJSON", func() {
	Action("list by source", func() {
		Routing(GET("sources/:sourceId/geojson"))
		Description("List an source's features in a GeoJSON.")
		Params(func() {
			Param("sourceId", Integer)
			Required("sourceId")
			Param("descending", Boolean)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(PagedGeoJSON)
		})
	})

	Action("list by id", func() {
		Routing(GET("features/:featureId/geojson"))
		Description("List a feature's GeoJSON by id.")
		Params(func() {
			Param("featureId", Integer)
			Required("featureId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(PagedGeoJSON)
		})
	})

	Action("geographical query", func() {
		Routing(GET("features"))
		Description("List features in a geographical area.")
		Params(func() {
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(MapFeatures)
		})
	})
})
