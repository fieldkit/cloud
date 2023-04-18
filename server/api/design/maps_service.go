package design

import (
	. "goa.design/goa/v3/dsl"
)

var MapGeometry = ResultType("application/vnd.app.maps.geometry", func() {
	TypeName("MapGeometry")
	Attributes(func() {
		Attribute("type", String)
		Attribute("coordinates", ArrayOf(Float64))
		Required("type")
		Required("coordinates")
	})
	View("default", func() {
		Attribute("type")
		Attribute("coordinates")
	})
})

var MapGeoJson = ResultType("application/vnd.app.maps.geo", func() {
	TypeName("MapGeoJson")
	Attributes(func() {
		Attribute("type", String)
		Attribute("geometry", MapGeometry)
		Attribute("properties", MapOf(String, String))
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

var Map = ResultType("application/vnd.app.maps", func() {
	TypeName("Map")
	Attributes(func() {
		Attribute("features", ArrayOf(MapGeoJson))
		Required("features")
	})
	View("default", func() {
		Attribute("features")
	})
})

var _ = Service("maps", func() {
	Method("coverage", func() {
		Result(Map)

		HTTP(func() {
			GET("maps/coverage")
		})
	})

	Error("unauthorized", String, "credentials are invalid")

	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
	})

	commonOptions()
})

