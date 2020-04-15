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

var LastUpload = MediaType("application/vnd.app.upload+json", func() {
	TypeName("LastUpload")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("time", DateTime)
		Attribute("upload_id", String)
		Attribute("size", Integer)
		Attribute("url", String)
		Attribute("type", String)
		Attribute("blocks", ArrayOf(Integer))
		Required("id")
		Required("time")
		Required("upload_id")
		Required("size")
		Required("url")
		Required("type")
		Required("blocks")
	})
	View("default", func() {
		Attribute("id")
		Attribute("time")
		Attribute("upload_id")
		Attribute("size")
		Attribute("url")
		Attribute("type")
		Attribute("blocks")
	})
})

var ImageRef = MediaType("application/vnd.app.imageref+json", func() {
	TypeName("ImageRef")
	Attributes(func() {
		Attribute("url", String)
		Required("url")
	})
	View("default", func() {
		Attribute("url")
	})
})

var StationPhotos = MediaType("application/vnd.app.station.photos+json", func() {
	TypeName("StationPhotos")
	Attributes(func() {
		Attribute("small", String)
		Required("small")
	})
	View("default", func() {
		Attribute("small")
	})
})

var StationOwner = MediaType("application/vnd.app.station.owner+json", func() {
	TypeName("StationOwner")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("name", String)
		Required("id", "name")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
	})
})

var Station = MediaType("application/vnd.app.station+json", func() {
	TypeName("Station")
	Reference(AddStationPayload)
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("name")
		Attribute("owner", StationOwner)
		Attribute("device_id", String)
		Attribute("last_uploads", CollectionOf(LastUpload))
		Attribute("status_json", HashOf(String, Any))
		Attribute("images", CollectionOf(ImageRef))
		Attribute("photos", StationPhotos)
		Attribute("read_only", Boolean)
		Required("id", "name", "owner", "device_id", "status_json", "images", "photos", "read_only")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("owner")
		Attribute("device_id")
		Attribute("last_uploads")
		Attribute("status_json")
		Attribute("images")
		Attribute("photos")
		Attribute("read_only")
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

var BadRequestResponse = MediaType("application/vnd.brr+json", func() {
	TypeName("BadRequestResponse")
	Attributes(func() {
		Attribute("key", String)
		Attribute("message", String)
		Required("key", "message")
	})
	View("default", func() {
		Attribute("key")
		Attribute("message")
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
		Response(BadRequest, func() {
			Media(BadRequestResponse)
		})
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
		Response(NotFound)
		Response(BadRequest)
		Response(OK, func() {
			Media(Station)
		})
	})

	Action("get", func() {
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

	Action("list project", func() {
		Routing(GET("projects/:projectId/stations"))
		Description("List project stations")
		Response(BadRequest)
		Response(OK, func() {
			Media(Stations)
		})
	})

	Action("delete", func() {
		Routing(DELETE("stations/:stationId"))
		Description("Delete station")
		Params(func() {
			Param("stationId", Integer)
		})
		Response(BadRequest)
		Response(OK, func() {
			Status(204)
		})
	})

	Action("photo", func() {
		NoSecurity()
		Routing(GET("/stations/:stationId/photo"))
		Params(func() {
			Param("stationId", Integer)
			Required("stationId")
		})
		Response(OK, func() {
			Media("image/jpeg")
		})
	})
})
