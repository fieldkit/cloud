package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var DeviceFileUrls = Type("DeviceFileUrls", func() {
	Attribute("csv", String)
	Attribute("fkpb", String)
	Attribute("json", String)
	Required("csv", "fkpb", "json")
})

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
		Attribute("file_type_name", String)
		Attribute("size", Integer)
		Attribute("url", String)
		Attribute("urls", DeviceFileUrls)
		Required("id")
		Required("time")
		Required("file_id")
		Required("device_id")
		Required("firmware")
		Required("meta")
		Required("file_type_id")
		Required("file_type_name")
		Required("size")
		Required("url")
		Required("urls")
	})
	View("default", func() {
		Attribute("id")
		Attribute("time")
		Attribute("file_id")
		Attribute("device_id")
		Attribute("firmware")
		Attribute("meta")
		Attribute("file_type_id")
		Attribute("file_type_name")
		Attribute("size")
		Attribute("url")
		Attribute("urls")
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

var DeviceFileTypeUrls = Type("DeviceFileTypeUrls", func() {
	Attribute("generate", String)
	Attribute("info", String)
	Attribute("csv", String)
	Attribute("fkpb", String)
	Required("generate", "info", "csv", "fkpb")
})

var DeviceSummaryUrls = Type("DeviceSummaryUrls", func() {
	Attribute("logs", DeviceFileTypeUrls)
	Attribute("data", DeviceFileTypeUrls)
	Required("logs", "data")
})

var DeviceSummary = MediaType("application/vnd.app.device+json", func() {
	TypeName("Device")
	Attributes(func() {
		Attribute("device_id", String)
		Attribute("number_of_files", Integer)
		Attribute("last_file_time", DateTime)
		Attribute("last_file_id", String)
		Attribute("urls", DeviceSummaryUrls)
		Required("device_id")
		Required("number_of_files")
		Required("last_file_time")
		Required("last_file_id")
		Required("urls")
	})
	View("default", func() {
		Attribute("device_id")
		Attribute("number_of_files")
		Attribute("last_file_time")
		Attribute("last_file_id")
		Attribute("urls")
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

var _ = Resource("device_logs", func() {
	Action("all", func() {
		Routing(GET("devices/:deviceId/logs"))
		Description("Get all device logs")
		Response(NotFound)
		Response("Busy", func() {
			Status(503)
		})
		Response(Found, func() {
			Headers(func() {
				Header("Location", String)
			})
		})
		Response(OK, func() {
			Status(200)
		})
	})
})

var _ = Resource("device_data", func() {
	Action("all", func() {
		Routing(GET("devices/:deviceId/data"))
		Description("Get all device data")
		Response(NotFound)
		Response("Busy", func() {
			Status(503)
		})
		Response(Found, func() {
			Headers(func() {
				Header("Location", String)
			})
		})
		Response(OK, func() {
			Status(200)
		})
	})
})

var _ = Resource("files", func() {
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

	Action("raw", func() {
		Routing(GET("files/:fileId/data.fkpb"))
		Description("Export file")
		Params(func() {
			Param("dl", Boolean, func() {
				Default(true)
			})
		})
		Response(NotFound)
		Response(OK, func() {
			Status(200)
		})
	})

	Action("file", func() {
		Routing(GET("files/:fileId"))
		Description("File info")
		Response(NotFound)
		Response(OK, func() {
			Media(DeviceFileSummary)
		})
	})

	Action("csv", func() {
		Routing(GET("files/:fileId/data.csv"))
		Description("Export file")
		Params(func() {
			Param("dl", Boolean, func() {
				Default(true)
			})
		})
		Response(NotFound)
		Response(OK, func() {
			Status(200)
		})
	})

	Action("json", func() {
		Routing(GET("files/:fileId/data.json"))
		Description("Export file")
		Params(func() {
			Param("dl", Boolean, func() {
				Default(true)
			})
		})
		Response(NotFound)
		Response(OK, func() {
			Status(200)
		})
	})
})
