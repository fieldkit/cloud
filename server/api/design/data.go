package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("data", func() {
	Security(JWT, func() {
		Scope("api:access")
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
			Param("page", Integer)
			Param("pageSize", Integer)
		})
		Response(NotFound)
		Response(OK, func() {
			Media(DeviceDataRecordsResponse)
		})
	})
})

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
