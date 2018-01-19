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
		Attribute("expeditionId")
		Attribute("teamId")
		Attribute("userId")
		Attribute("active", Boolean)
		Attribute("name", String)
		Attribute("token", String)
		Attribute("key", String)
		Attribute("numberOfFeatures", Integer)
		Attribute("lastFeatureId", Integer)
		Required("id", "expeditionId", "active", "name", "token", "key")
	})
	View("default", func() {
		Attribute("id")
		Attribute("expeditionId")
		Attribute("teamId")
		Attribute("userId")
		Attribute("active")
		Attribute("name")
		Attribute("token")
		Attribute("key")
	})
	View("public", func() {
		Attribute("id")
		Attribute("expeditionId")
		Attribute("teamId")
		Attribute("userId")
		Attribute("active")
		Attribute("name")
		Attribute("numberOfFeatures")
		Attribute("lastFeatureId")
		Required("id", "expeditionId", "active", "name", "numberOfFeatures", "lastFeatureId")
	})
})

var DeviceSchema = MediaType("application/vnd.app.device_schema+json", func() {
	TypeName("DeviceSchema")
	Attributes(func() {
		Attribute("deviceId", Integer)
		Attribute("schemaId", Integer)
		Attribute("projectId", Integer)
		Attribute("active", Boolean)
		Attribute("jsonSchema", String)
		Attribute("key", String)
		Required("deviceId", "schemaId", "projectId", "active", "jsonSchema", "key")
	})
	View("default", func() {
		Attribute("deviceId")
		Attribute("schemaId")
		Attribute("projectId")
		Attribute("active")
		Attribute("jsonSchema")
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
		Attribute("deviceInputs", CollectionOf(DeviceInput))
		Required("deviceInputs")
	})
	View("default", func() {
		Attribute("deviceInputs")
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
	Attribute("jsonSchema")
	Required("jsonSchema")
})

var UpdateDeviceInputLocationPayload = Type("UpdateDeviceInputLocationPayload", func() {
	Reference(Input)
	Attribute("key")
	Required("key")
	Attribute("longitude", Number)
	Required("longitude")
	Attribute("latitude", Number)
	Required("latitude")
})

var _ = Resource("device", func() {
	Security(JWT, func() {
		Scope("api:access")
	})

	Action("add", func() {
		Routing(POST("expeditions/:expeditionId/inputs/devices"))
		Description("Add a Device input")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
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

	Action("update location", func() {
		Routing(PATCH("inputs/devices/:id/location"))
		Description("Update an Device input location")
		Params(func() {
			Param("id", Integer)
			Required("id")
		})
		Payload(UpdateDeviceInputLocationPayload)
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
