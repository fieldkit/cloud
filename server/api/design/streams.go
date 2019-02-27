package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var DeviceStreamSummary = MediaType("application/vnd.app.device.stream+json", func() {
	TypeName("DeviceStream")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("time", DateTime)
		Attribute("stream_id", String)
		Attribute("device_id", String)
		Attribute("firmware", String)
		Attribute("meta", String)
		Attribute("file_id", String)
		Attribute("size", Integer)
		Attribute("url", String)
		Required("id")
		Required("time")
		Required("stream_id")
		Required("device_id")
		Required("firmware")
		Required("meta")
		Required("file_id")
		Required("size")
		Required("url")
	})
	View("default", func() {
		Attribute("id")
		Attribute("time")
		Attribute("stream_id")
		Attribute("device_id")
		Attribute("firmware")
		Attribute("meta")
		Attribute("file_id")
		Attribute("size")
		Attribute("url")
	})
})

var DeviceStreams = MediaType("application/vnd.app.device.streams+json", func() {
	TypeName("DeviceStreams")
	Attributes(func() {
		Attribute("streams", CollectionOf(DeviceStreamSummary))
		Required("streams")
	})
	View("default", func() {
		Attribute("streams")
	})
})

var DeviceSummary = MediaType("application/vnd.app.device+json", func() {
	TypeName("Device")
	Attributes(func() {
		Attribute("device_id", String)
		Attribute("number_of_streams", Integer)
		Attribute("last_stream_time", DateTime)
		Attribute("last_stream_id", String)
		Required("device_id")
		Required("number_of_streams")
		Required("last_stream_time")
		Required("last_stream_id")
	})
	View("default", func() {
		Attribute("device_id")
		Attribute("number_of_streams")
		Attribute("last_stream_time")
		Attribute("last_stream_id")
	})
})

var Devices = MediaType("application/vnd.app.devices+json", func() {
	TypeName("Devices")
	Attributes(func() {
		Attribute("devices", CollectionOf(DeviceSummary))
		Required("devices")
	})
	View("default", func() {
		Attribute("devices")
	})
})

var _ = Resource("Streams", func() {
	Action("list device", func() {
		Routing(GET("devices/:deviceId/streams"))
		Description("List device streams")
		Params(func() {
			Param("page", Integer)
			Param("file_id", String)
		})
		Response(OK, func() {
			Media(DeviceStreams)
		})
	})

	Action("list devices", func() {
		Routing(GET("streams/devices"))
		Description("List devices")
		Response(NotFound)
		Response(OK, func() {
			Media(Devices)
		})
	})

	Action("list all", func() {
		Routing(GET("streams"))
		Description("List streams")
		Params(func() {
			Param("page", Integer)
			Param("file_id", String)
		})
		Response(OK, func() {
			Media(DeviceStreams)
		})
	})

	Action("device data", func() {
		Routing(GET("devices/:deviceId/data"))
		Description("Export device data")
		Response(NotFound)
		Response(OK, func() {
			Status(200)
		})
	})

	Action("device logs", func() {
		Routing(GET("devices/:deviceId/logs"))
		Description("Export device logs")
		Response(NotFound)
		Response(OK, func() {
			Status(200)
		})
	})

	Action("csv", func() {
		Routing(GET("streams/:streamId/csv"))
		Description("Export stream")
		Response(NotFound)
		Response(OK, func() {
			Status(200)
		})
	})

	Action("json", func() {
		Routing(GET("streams/:streamId/json"))
		Description("Export stream")
		Response(NotFound)
		Response(OK, func() {
			Status(200)
		})
	})
})
