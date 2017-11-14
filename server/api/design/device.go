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
		Attribute("expedition_id")
		Attribute("team_id")
		Attribute("user_id")
		Attribute("active", Boolean)
		Attribute("name", String)
		Attribute("token", String)
		Attribute("key", String)
		Required("id", "expedition_id", "active", "name", "token", "key")
	})
	View("default", func() {
		Attribute("id")
		Attribute("expedition_id")
		Attribute("team_id")
		Attribute("user_id")
		Attribute("active")
		Attribute("name")
		Attribute("token")
		Attribute("key")
	})
})

var DeviceSchema = MediaType("application/vnd.app.device_schema+json", func() {
	TypeName("DeviceSchema")
	Attributes(func() {
		Attribute("device_id", Integer)
		Attribute("schema_id", Integer)
		Attribute("project_id", Integer)
		Attribute("active", Boolean)
		Attribute("json_schema", String)
		Attribute("key", String)
		Required("device_id", "schema_id", "project_id", "active", "json_schema", "key")
	})
	View("default", func() {
		Attribute("device_id")
		Attribute("schema_id")
		Attribute("project_id")
		Attribute("active")
		Attribute("json_schema")
		Attribute("key")
	})
})

var DeviceSchemas = MediaType("application/vnd.app.device_schemas+json", func() {
	TypeName("DeviceSchemas")
	Attributes(func() {
		Attribute("schemas", CollectionOf(DeviceSchema))
		Required("schemas")
	})
	View("default", func() {
		Attribute("schemas")
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
	Attribute("name")
	Required("name")
	Attribute("key")
	Required("key")
})

var UpdateDeviceInputSchemaPayload = Type("UpdateDeviceInputSchemaPayload", func() {
	Reference(Input)
	Attribute("key")
	Required("key")
	Attribute("active")
	Required("active")
	Attribute("json_schema")
	Required("json_schema")
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
			Media(DeviceInput)
		})
	})

	Action("update", func() {
		Routing(PATCH("inputs/devices/:id"))
		Description("Update an Device input")
		Params(func() {
			Param("id", Integer)
			Required("id")
		})
		Payload(UpdateDeviceInputPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(DeviceInput)
		})
	})

	Action("update schema", func() {
		Routing(PATCH("inputs/devices/:id/schemas"))
		Description("Update an Device input schema")
		Params(func() {
			Param("id", Integer)
			Required("id")
		})
		Payload(UpdateDeviceInputSchemaPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(DeviceSchemas)
		})
	})

	Action("get id", func() {
		Routing(GET("inputs/devices/:id"))
		Description("Get a Device input")
		Params(func() {
			Param("id", Integer)
			Required("id")
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
