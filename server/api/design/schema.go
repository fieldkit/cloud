package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var Schema = MediaType("application/vnd.app.schema+json", func() {
	TypeName("Schema")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("project_id", Integer)
		Attribute("json_schema", Any)
		Required("id", "json_schema")
	})
	View("default", func() {
		Attribute("id")
		Attribute("project_id")
		Attribute("json_schema")
	})
})

var AddSchemaPayload = Type("AddSchemaPayload", func() {
	Reference(Schema)
	Attribute("json_schema")
	Required("json_schema")
})

var UpdateSchemaPayload = Type("UpdateSchemaPayload", func() {
	Reference(Schema)
	Attribute("project_id")
	Attribute("json_schema")
	Required("json_schema")
})

var Schemas = MediaType("application/vnd.app.schemas+json", func() {
	TypeName("Schemas")
	Attributes(func() {
		Attribute("schemas", CollectionOf(Schema))
	})
	View("default", func() {
		Attribute("schemas")
	})
})

var _ = Resource("schema", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("projects/:project_id/schemas"))
		Description("Add a schema")
		Params(func() {
			Param("project_id", Integer)
			Required("project_id")
		})
		Payload(AddSchemaPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Schema)
		})
	})

	Action("update", func() {
		Routing(PATCH("schemas/:schema_id"))
		Description("Update a schema")
		Params(func() {
			Param("schema_id", Integer)
			Required("schema_id")
		})
		Payload(UpdateSchemaPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Schema)
		})
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/schemas"))
		Description("List a project's schemas")
		Params(func() {
			Param("project", String, func() {
				Pattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`)
				Description("Project slug")
			})
			Required("project")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Schemas)
		})
	})

	Action("list id", func() {
		Routing(GET("projects/:project_id/schemas"))
		Description("List a project's schemas")
		Params(func() {
			Param("project_id", Integer)
			Required("project_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Schemas)
		})
	})
})
