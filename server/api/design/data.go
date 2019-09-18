package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var DeviceMetaRecord = MediaType("application/vnd.app.device.meta.record+json", func() {
	TypeName("DeviceMetaRecord")
	Attributes(func() {
		Attribute("time", DateTime)
		Attribute("record", Integer)
		Attribute("data", HashOf(String, Any))
		Required("time", "record", "data")
	})
	View("default", func() {
		Attribute("time")
		Attribute("record")
		Attribute("data")
	})
})

var DeviceDataRecord = MediaType("application/vnd.app.device.data.record+json", func() {
	TypeName("DeviceDataRecord")
	Attributes(func() {
		Attribute("time", DateTime)
		Attribute("record", Integer)
		Attribute("data", HashOf(String, Any))
		Required("time", "record", "data")
	})
	View("default", func() {
		Attribute("time")
		Attribute("record")
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
		Attribute("summary", DeviceDataStreamsSummary)
		Attribute("meta", CollectionOf(DeviceMetaRecord))
		Attribute("data", CollectionOf(DeviceDataRecord))

		Required("summary")
		Required("meta")
		Required("data")
	})
	View("default", func() {
		Attribute("summary")
		Attribute("meta")
		Attribute("data")
	})
})

var _ = Resource("data", func() {
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

	Action("device", func() {
		Routing(GET("data/devices/:deviceId"))
		Description("Retrieve data")
		Params(func() {
			Param("first_block", Integer)
			Param("last_block", Integer)
			Param("page", Integer)
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
})
