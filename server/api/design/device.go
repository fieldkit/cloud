package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var DeviceInput = MediaType("application/vnd.app.device_input+json", func() {
	TypeName("DeviceInput")
	Reference(Input)
	Attributes(func() {
		Attribute("id")
		Attribute("token")
		Attribute("key")
	})
	View("default", func() {
		Attribute("id")
		Attribute("token")
		Attribute("key")
	})
})

var DeviceInputs = MediaType("application/vnd.app.device_inputs+json", func() {
	TypeName("DeviceInputs")
	Attributes(func() {
		Attribute("device_inputs", CollectionOf(DeviceInput))
		Required("device_inputs")
	})
	View("default", func() {
		Attribute("device_inputs")
	})
})

var AddDeviceInputPayload = Type("AddDeviceInputPayload", func() {
	Reference(Input)
	Attribute("name")
	Required("name")
	Attribute("key")
	Required("key")
})

var UpdateDeviceInputPayload = Type("UpdateDeviceInputPayload", func() {
	Reference(Input)
	Attribute("schema")
	Required("schema")
})

var _ = Resource("device", func() {
	Security(JWT, func() {
		Scope("api:access")
	})

	Action("add", func() {
		Routing(POST("expeditions/:expedition_id/inputs/devices"))
		Description("Add a Device input")
		Params(func() {
			Param("expedition_id", Integer)
			Required("expedition_id")
		})
		Payload(AddDeviceInputPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Location)
		})
	})

	Action("update", func() {
		Routing(PATCH("inputs/devices/:input_id"))
		Description("Update an Device input")
		Params(func() {
			Param("input_id", Integer)
			Required("input_id")
		})
		Payload(UpdateDeviceInputPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(DeviceInput)
		})
	})

	Action("get id", func() {
		Routing(GET("inputs/devices/:input_id"))
		Description("Get a Device input")
		Params(func() {
			Param("input_id", Integer)
			Required("input_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(DeviceInput)
		})
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/inputs/devices"))
		Description("List an expedition's Device inputs")
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("expedition", String, ExpeditionSlug)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(DeviceInputs)
		})
	})
})
