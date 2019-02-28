package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var DeviceFileSummary = MediaType("application/vnd.app.device.file+json", func() {
	TypeName("DeviceFile")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("time", DateTime)
		Attribute("file_id", String)
		Attribute("device_id", String)
		Attribute("firmware", String)
		Attribute("meta", String)
		Attribute("file_type_id", String)
		Attribute("size", Integer)
		Attribute("url", String)
		Required("id")
		Required("time")
		Required("file_id")
		Required("device_id")
		Required("firmware")
		Required("meta")
		Required("file_type_id")
		Required("size")
		Required("url")
	})
	View("default", func() {
		Attribute("id")
		Attribute("time")
		Attribute("file_id")
		Attribute("device_id")
		Attribute("firmware")
		Attribute("meta")
		Attribute("file_type_id")
		Attribute("size")
		Attribute("url")
	})
})

var DeviceFiles = MediaType("application/vnd.app.device.files+json", func() {
	TypeName("DeviceFiles")
	Attributes(func() {
		Attribute("files", CollectionOf(DeviceFileSummary))
		Required("files")
	})
	View("default", func() {
		Attribute("files")
	})
})

var DeviceSummary = MediaType("application/vnd.app.device+json", func() {
	TypeName("Device")
	Attributes(func() {
		Attribute("device_id", String)
		Attribute("number_of_files", Integer)
		Attribute("last_file_time", DateTime)
		Attribute("last_file_id", String)
		Required("device_id")
		Required("number_of_files")
		Required("last_file_time")
		Required("last_file_id")
	})
	View("default", func() {
		Attribute("device_id")
		Attribute("number_of_files")
		Attribute("last_file_time")
		Attribute("last_file_id")
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

var _ = Resource("Files", func() {
	Action("list devices", func() {
		Routing(GET("files/devices"))
		Description("List devices")
		Response(NotFound)
		Response(OK, func() {
			Media(Devices)
		})
	})

	Action("list device data files", func() {
		Routing(GET("devices/:deviceId/files/data"))
		Description("List device files")
		Params(func() {
			Param("page", Integer)
		})
		Response(OK, func() {
			Media(DeviceFiles)
		})
	})

	Action("list device log files", func() {
		Routing(GET("devices/:deviceId/files/logs"))
		Description("List device files")
		Params(func() {
			Param("page", Integer)
		})
		Response(OK, func() {
			Media(DeviceFiles)
		})
	})

	/*
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
	*/

	Action("csv", func() {
		Routing(GET("files/:fileId/csv"))
		Description("Export file")
		Response(NotFound)
		Response(OK, func() {
			Status(200)
		})
	})

	Action("json", func() {
		Routing(GET("files/:fileId/json"))
		Description("Export file")
		Response(NotFound)
		Response(OK, func() {
			Status(200)
		})
	})
})
