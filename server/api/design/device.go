package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var DeviceSource = MediaType("application/vnd.app.device_source+json", func() {
	TypeName("DeviceSource")
	Reference(Source)
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
		Attribute("startTime", DateTime)
		Attribute("endTime", DateTime)
		Attribute("centroid", ArrayOf(Number))
		Attribute("radius", Number)
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
		Attribute("startTime")
		Attribute("endTime")
		Attribute("centroid")
		Attribute("radius")
		Required("id", "expeditionId", "active", "name", "numberOfFeatures", "lastFeatureId", "startTime", "endTime", "centroid", "radius")
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

var DeviceSources = MediaType("application/vnd.app.device_sources+json", func() {
	TypeName("DeviceSources")
	Attributes(func() {
		Attribute("deviceSources", CollectionOf(DeviceSource))
		Required("deviceSources")
	})
	View("default", func() {
		Attribute("deviceSources")
	})
})

var AddDeviceSourcePayload = Type("AddDeviceSourcePayload", func() {
	Reference(Source)
	Attribute("name")
	Required("name")
	Attribute("key")
	Required("key")
})

var UpdateDeviceSourcePayload = Type("UpdateDeviceSourcePayload", func() {
	Reference(Source)
	Attribute("name")
	Required("name")
	Attribute("key")
	Required("key")
})

var UpdateDeviceSourceSchemaPayload = Type("UpdateDeviceSourceSchemaPayload", func() {
	Reference(Source)
	Attribute("key")
	Required("key")
	Attribute("active")
	Required("active")
	Attribute("jsonSchema")
	Required("jsonSchema")
})

var UpdateDeviceSourceLocationPayload = Type("UpdateDeviceSourceLocationPayload", func() {
	Reference(Source)
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
		Routing(POST("expeditions/:expeditionId/sources/devices"))
		Description("Add a Device source")
		Params(func() {
			Param("expeditionId", Integer)
			Required("expeditionId")
		})
		Payload(AddDeviceSourcePayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(DeviceSource)
		})
	})

	Action("update", func() {
		Routing(PATCH("sources/devices/:id"))
		Description("Update an Device source")
		Params(func() {
			Param("id", Integer)
			Required("id")
		})
		Payload(UpdateDeviceSourcePayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(DeviceSource)
		})
	})

	Action("update location", func() {
		Routing(PATCH("sources/devices/:id/location"))
		Description("Update an Device source location")
		Params(func() {
			Param("id", Integer)
			Required("id")
		})
		Payload(UpdateDeviceSourceLocationPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(DeviceSource)
		})
	})

	Action("update schema", func() {
		Routing(PATCH("sources/devices/:id/schemas"))
		Description("Update an Device source schema")
		Params(func() {
			Param("id", Integer)
			Required("id")
		})
		Payload(UpdateDeviceSourceSchemaPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(DeviceSchemas)
		})
	})

	Action("get id", func() {
		Routing(GET("sources/devices/:id"))
		Description("Get a Device source")
		Params(func() {
			Param("id", Integer)
			Required("id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(DeviceSource)
		})
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/sources/devices"))
		Description("List an expedition's Device sources")
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("expedition", String, ExpeditionSlug)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(DeviceSources)
		})
	})
})
