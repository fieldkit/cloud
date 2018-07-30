package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddFirmwarePayload = Type("AddFirmwarePayload", func() {
	Attribute("etag")
	Required("etag")
	Attribute("module")
	Required("module")
	Attribute("url")
	Required("url")
	Attribute("meta")
	Required("meta")
})

var UpdateDeviceFirmwarePayload = Type("UpdateDeviceFirmwarePayload", func() {
	Reference(Source)
	Attribute("deviceId", Integer)
	Required("deviceId")
	Attribute("module")
	Required("module")
	Attribute("etag")
	Required("etag")
	Attribute("url")
	Required("url")
})

var _ = Resource("Firmware", func() {
	Action("check", func() {
		Routing(GET("devices/:deviceId/:module/firmware"))
		Description("Return firmware for a device")
		Headers(func() {
			Header("If-None-Match")
		})
		Params(func() {
			Param("deviceId", String)
			Required("deviceId")
			Param("module", String)
			Required("module")
		})
		Response(NotFound)
		Response(NotModified)
		Response(OK, func() {
			Status(200)
			Headers(func() {
				Header("ETag")
			})
		})
	})

	Action("update", func() {
		Routing(PATCH("devices/firmware"))
		Description("Update an Device firmware")
		Params(func() {
		})
		Payload(UpdateDeviceFirmwarePayload)
		Response(BadRequest)
		Response(OK)
	})

	Action("add", func() {
		Routing(PATCH("firmware"))
		Description("Add firmware")
		Params(func() {
		})
		Payload(AddFirmwarePayload)
		Response(BadRequest)
		Response(OK)
	})
})
