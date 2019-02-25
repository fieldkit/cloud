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

var _ = Resource("Streams", func() {
	Action("list device", func() {
		Routing(GET("devices/:deviceId/streams"))
		Description("List device streams")
		Params(func() {
			Param("page", Integer)
			Param("file_id", Integer)
		})
		Response(OK, func() {
			Media(DeviceStreams)
		})
	})

	Action("list all", func() {
		Routing(GET("devices/streams"))
		Description("List streams")
		Params(func() {
			Param("page", Integer)
			Param("file_id", Integer)
		})
		Response(OK, func() {
			Media(DeviceStreams)
		})
	})

	Action("raw", func() {
		Routing(GET("streams/:streamId/raw"))
		Description("Export stream")
		Params(func() {
		})
		Response(NotFound)
		Response(OK, func() {
			Status(200)
		})
	})

	Action("csv", func() {
		Routing(GET("streams/:streamId/csv"))
		Description("Export stream")
		Params(func() {
		})
		Response(NotFound)
		Response(OK, func() {
			Status(200)
		})
	})
})
