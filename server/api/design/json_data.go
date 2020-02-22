package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("jsonData", func() {
	Security(JWT, func() {
		Scope("api:access")
	})

	Action("get", func() {
		Routing(GET("data/devices/:deviceId/data/json"))
		Description("Retrieve data")
		Params(func() {
			Param("page", Integer)
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

	Action("get lines", func() {
		Routing(GET("data/devices/:deviceId/data/lines"))
		Description("Retrieve data")
		Params(func() {
			Param("page", Integer)
			Param("pageSize", Integer)
			Param("start", Integer)
			Param("end", Integer)
			Param("internal", Boolean)
		})
		Response(NotFound)
		Response(OK)
	})

	Action("summary", func() {
		Routing(GET("data/devices/:deviceId/summary/json"))
		Description("Retrieve summarized data")
		Params(func() {
			Param("page", Integer)
			Param("pageSize", Integer)
			Param("start", Integer)
			Param("end", Integer)
			Param("resolution", Integer)
			Param("internal", Boolean)
		})
		Response(NotFound)
		Response(OK, func() {
			Media(JSONDataSummaryResponse)
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
	Attribute("version", String)
	Attribute("build", String)
	Attribute("number", String)
	Attribute("timestamp", Integer)
	Attribute("hash", String)

	Required("version", "build", "number", "timestamp", "hash")
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
	Attribute("metas", ArrayOf(Integer))
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

var JSONDataSummaryResponse = MediaType("application/vnd.app.device.json.data.summary+json", func() {
	TypeName("JSONDataSummaryResponse")
	Attributes(func() {
		Attribute("modules", ArrayOf(JSONDataMetaModule))
		Attribute("data", ArrayOf(JSONDataRow))
		Required("modules", "data")
	})
	View("default", func() {
		Attribute("modules")
		Attribute("data")
	})
})
