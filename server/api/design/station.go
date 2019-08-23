package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddStationPayload = Type("AddStationPayload", func() {
	Attribute("name", String)
	Attribute("device_id", String)
	Attribute("status_json", HashOf(String, Any))
	Required("name", "device_id", "status_json")
})

var UpdateStationPayload = Type("UpdateStationPayload", func() {
	Attribute("name", String)
	Attribute("status_json", HashOf(String, Any))
	Required("name", "status_json")
})

var Station = MediaType("application/vnd.app.station+json", func() {
	TypeName("Station")
	Reference(AddStationPayload)
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("name")
		Attribute("owner_id", Integer)
		Attribute("device_id", String)
		Attribute("status_json", HashOf(String, Any))
		Required("id", "name", "owner_id", "device_id", "status_json")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("owner_id")
		Attribute("device_id")
		Attribute("status_json")
	})
})

var Stations = MediaType("application/vnd.app.stations+json", func() {
	TypeName("Stations")
	Attributes(func() {
		Attribute("stations", CollectionOf(Station))
		Required("stations")
	})
	View("default", func() {
		Attribute("stations")
	})
})

var _ = Resource("station", func() {
	Security(JWT, func() {
		Scope("api:access")
	})

	Action("add", func() {
		Routing(POST("stations"))
		Description("Add a station")
		Payload(AddStationPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Station)
		})
	})

	Action("update", func() {
		Routing(PATCH("stations/:stationId"))
		Description("Update a station")
		Params(func() {
			Param("stationId", Integer)
			Required("stationId")
		})
		Payload(UpdateStationPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Station)
		})
	})

	Action("get", func() {
		NoSecurity()
		Routing(GET("stations/@/:stationId"))
		Description("Get a station")
		Params(func() {
			Param("stationId", Integer)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Station)
		})
	})

	Action("list", func() {
		Routing(GET("stations"))
		Description("List stations")
		Response(BadRequest)
		Response(OK, func() {
			Media(Stations)
		})
	})
})
