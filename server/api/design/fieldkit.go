package design

import (
	"math"

	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddFieldkitInputPayload = Type("AddFieldkitInputPayload", func() {
	Reference(Input)
	Attribute("name")
	Required("name")
})

var UpdateFieldkitInputPayload = Type("UpdateFieldkitInputPayload", func() {
	Reference(Input)
	Attribute("name")
	Attribute("team_id")
	Attribute("user_id")
})

var FieldkitInput = MediaType("application/vnd.app.fieldkit_input+json", func() {
	TypeName("FieldkitInput")
	Reference(Input)
	Attributes(func() {
		Attribute("id")
		Attribute("expedition_id")
		Attribute("name")
		Attribute("team_id")
		Attribute("user_id")
		Required("id", "expedition_id", "name")
	})
	View("default", func() {
		Attribute("id")
		Attribute("expedition_id")
		Attribute("name")
		Attribute("team_id")
		Attribute("user_id")
	})
})

var FieldkitInputs = MediaType("application/vnd.app.fieldkit_inputs+json", func() {
	TypeName("FieldkitInputs")
	Attributes(func() {
		Attribute("fieldkit_inputs", CollectionOf(FieldkitInput))
		Required("fieldkit_inputs")
	})
	View("default", func() {
		Attribute("fieldkit_inputs")
	})
})

var FieldkitBinary = MediaType("application/vnd.app.fieldkit_input_binary+json", func() {
	TypeName("FieldkitBinary")
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("input_id", Integer)
		Attribute("fields", ArrayOf(String, func() {
			Enum("varint", "uvarint", "float32", "float64")
		}))
		Required("id", "input_id", "fields")
	})
	View("default", func() {
		Attribute("id")
		Attribute("input_id")
		Attribute("fields")
	})
})

var SetFieldkitBinaryPayload = Type("SetFieldkitBinaryPayload", func() {
	Reference(FieldkitBinary)
	Attribute("fields")
	Required("fields")
})

var _ = Resource("fieldkit", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("expeditions/:expedition_id/inputs/fieldkits"))
		Description("Add a Fieldkit input")
		Params(func() {
			Param("expedition_id", Integer)
			Required("expedition_id")
		})
		Payload(AddFieldkitInputPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(FieldkitInput)
		})
	})

	Action("get id", func() {
		Routing(GET("inputs/fieldkits/:input_id"))
		Description("Get a Fieldkit input")
		Params(func() {
			Param("input_id", Integer)
			Required("input_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(FieldkitInput)
		})
	})

	Action("set binary", func() {
		Routing(PUT("inputs/fieldkits/:input_id/binary/:binary_id"))
		Description("Set a Fieldkit binary format")
		Params(func() {
			Param("input_id", Integer)
			Param("binary_id", Integer, func() {
				Maximum(math.MaxInt16)
			})
			Required("input_id", "binary_id")
		})
		Payload(SetFieldkitBinaryPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(FieldkitBinary)
		})
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/inputs/fieldkits"))
		Description("List an expedition's Fieldkit inputs")
		Params(func() {
			Param("project", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Project slug")
			})
			Param("expedition", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Expedition slug")
			})
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(FieldkitInputs)
		})
	})

	Action("list id", func() {
		Routing(GET("expeditions/:expedition_id/inputs/fieldkits"))
		Description("List an expedition's Fieldkit inputs")
		Params(func() {
			Param("expedition_id", Integer)
			Required("expedition_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(FieldkitInputs)
		})
	})
})
