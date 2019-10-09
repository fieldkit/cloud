package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var DeviceMetaRecord = MediaType("application/vnd.app.device.meta.record+json", func() {
	TypeName("DeviceMetaRecord")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("time", DateTime)
		Attribute("record", Integer)
		Attribute("data", HashOf(String, Any))
		Required("id", "time", "record", "data")
	})
	View("default", func() {
		Attribute("id")
		Attribute("time")
		Attribute("record")
		Attribute("data")
	})
})

var DeviceDataRecord = MediaType("application/vnd.app.device.data.record+json", func() {
	TypeName("DeviceDataRecord")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("time", DateTime)
		Attribute("record", Integer)
		Attribute("meta", Integer)
		Attribute("location", ArrayOf(Number))
		Attribute("data", HashOf(String, Any))
		Required("id", "time", "record", "meta", "location", "data")
	})
	View("default", func() {
		Attribute("id")
		Attribute("time")
		Attribute("record")
		Attribute("meta")
		Attribute("location")
		Attribute("data")
	})
})

var DeviceStreamSummary = Type("DeviceStreamSummary", func() {
	Attribute("records", Integer)
	Attribute("size", Integer)
})

var DeviceDataStreamsSummary = Type("DeviceDataStreamsSummary", func() {
	Attribute("meta", DeviceStreamSummary)
	Attribute("data", DeviceStreamSummary)
})

var DeviceDataRecordsResponse = MediaType("application/vnd.app.device.data+json", func() {
	TypeName("DeviceDataRecordsResponse")
	Attributes(func() {
		Attribute("meta", CollectionOf(DeviceMetaRecord))
		Attribute("data", CollectionOf(DeviceDataRecord))
		Required("meta")
		Required("data")
	})
	View("default", func() {
		Attribute("meta")
		Attribute("data")
	})
})

var DeviceProvisionSummary = MediaType("application/vnd.app.device.provision.summary+json", func() {
	TypeName("DeviceProvisionSummary")
	Attributes(func() {
		Attribute("generation", String)
		Attribute("created", DateTime)
		Attribute("updated", DateTime)
		Attribute("meta", DeviceMetaSummary)
		Attribute("data", DeviceDataSummary)

		Required("generation")
		Required("created")
		Required("updated")
		Required("meta")
		Required("data")
	})
	View("default", func() {
		Attribute("generation")
		Attribute("created")
		Attribute("updated")
		Attribute("meta")
		Attribute("data")
	})
})

var DeviceMetaSummary = MediaType("application/vnd.app.device.meta.summary+json", func() {
	TypeName("DeviceMetaSummary")
	Attributes(func() {
		Attribute("size", Integer)
		Attribute("first", Integer)
		Attribute("last", Integer)
		Required("size")
		Required("first")
		Required("last")
	})
	View("default", func() {
		Attribute("size")
		Attribute("first")
		Attribute("last")
	})
})

var DeviceDataSummary = MediaType("application/vnd.app.device.data.summary+json", func() {
	TypeName("DeviceDataSummary")
	Attributes(func() {
		Attribute("size", Integer)
		Attribute("first", Integer)
		Attribute("last", Integer)
		Required("size")
		Required("first")
		Required("last")
	})
	View("default", func() {
		Attribute("size")
		Attribute("first")
		Attribute("last")
	})
})

var DeviceDataSummaryResponse = MediaType("application/vnd.app.device.summary+json", func() {
	TypeName("DeviceDataSummaryResponse")
	Attributes(func() {
		Attribute("provisions", CollectionOf(DeviceProvisionSummary))
		Required("provisions")
	})
	View("default", func() {
		Attribute("provisions")
	})
})

var _ = Resource("data", func() {
	Security(JWT, func() {
		Scope("api:access")
	})

	Action("process", func() {
		Routing(GET("data/process"))
		Description("Process data")
		Response("Busy", func() {
			Status(503)
		})
		Response(OK, func() {
			Status(200)
		})
	})

	Action("device summary", func() {
		Routing(GET("data/devices/:deviceId/summary"))
		Description("Retrieve summary")
		Response(NotFound)
		Response(OK, func() {
			Media(DeviceDataSummaryResponse)
		})
	})

	Action("device data", func() {
		Routing(GET("data/devices/:deviceId/data"))
		Description("Retrieve data")
		Params(func() {
			Param("firstBlock", Integer)
			Param("lastBlock", Integer)
			Param("pageNumber", Integer)
			Param("pageSize", Integer)
		})
		Response(NotFound)
		Response(OK, func() {
			Media(DeviceDataRecordsResponse)
		})
	})

	Action("delete", func() {
		Routing(DELETE("data/ingestions/:ingestionId"))
		Description("Delete data")
		Params(func() {
			Param("ingestionId", Integer)
		})
		Response(NotFound)
		Response(OK, func() {
			Status(200)
		})
	})

	Action("process ingestion", func() {
		Routing(POST("data/ingestions/:ingestionId/process"))
		Description("Process ingestion")
		Params(func() {
			Param("ingestionId", Integer)
		})
		Response(NotFound)
		Response(OK, func() {
			Status(200)
		})
	})
})

var JSONDataMetaSensor = Type("JSONDataMetaSensor", func() {
	Attribute("name", String)
	Attribute("key", String)
	Attribute("units", String)

	Required("name", "key", "units")
})

var JSONDataMetaModule = Type("JSONDataMetaModule", func() {
	Attribute("id", String)
	Attribute("name", String)

	Attribute("manufacturer", Integer)
	Attribute("kind", Integer)
	Attribute("version", Integer)
	Attribute("sensors", ArrayOf(JSONDataMetaSensor))

	Required("id", "name", "kind", "version", "manufacturer")
})

var JSONDataMetaStationFirmware = Type("JSONDataMetaStationFirmware", func() {
	Attribute("git", String)
	Attribute("build", String)

	Required("git", "build")
})

var JSONDataMetaStation = Type("JSONDataMetaStation", func() {
	Attribute("id", String)
	Attribute("name", String)
	Attribute("modules", ArrayOf(JSONDataMetaModule))
	Attribute("firmware", JSONDataMetaStationFirmware)

	Required("id", "name", "modules", "firmware")
})

var JSONDataRow = Type("JSONDataRow", func() {
	Attribute("id", Integer)
	Attribute("time", Integer)
	Attribute("location", ArrayOf(Number))
	Attribute("d", HashOf(String, Any))

	Required("id", "time", "location", "d")
})

var JSONDataMeta = Type("JSONDataMeta", func() {
	Attribute("id", Integer)
	Attribute("station", JSONDataMetaStation)

	Required("id")
})

var JSONDataVersion = Type("JSONDataVersion", func() {
	Attribute("meta", JSONDataMeta)
	Attribute("data", ArrayOf(JSONDataRow))
})

var JSONDataResponse = MediaType("application/vnd.app.device.json.data+json", func() {
	TypeName("JSONDataResponse")
	Attributes(func() {
		Attribute("versions", ArrayOf(JSONDataVersion))
		Required("versions")
	})
	View("default", func() {
		Attribute("versions")
	})
})

var _ = Resource("jsonData", func() {
	Security(JWT, func() {
		Scope("api:access")
	})

	Action("get", func() {
		Routing(GET("data/devices/:deviceId/data/json"))
		Description("Retrieve data")
		Params(func() {
			Param("pageNumber", Integer)
			Param("pageSize", Integer)
			Param("start", Integer)
			Param("end", Integer)
			Param("internal", Boolean)
		})
		Response(NotFound)
		Response(OK, func() {
			Media(JSONDataResponse)
		})
	})
})
