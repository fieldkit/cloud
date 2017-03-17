package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddInputPayload = Type("AddInputPayload", func() {
	Attribute("name", String)
	Attribute("slug", String, func() {
		Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
		MaxLength(40)
	})
	Required("name", "slug")
})

var Input = MediaType("application/vnd.app.input+json", func() {
	TypeName("Input")
	Reference(AddInputPayload)
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("name")
		Attribute("slug")
		Required("id", "name", "slug")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("slug")
	})
})

var Inputs = MediaType("application/vnd.app.inputs+json", func() {
	TypeName("Inputs")
	Attributes(func() {
		Attribute("inputs", CollectionOf(Input))
		Required("inputs")
	})
	View("default", func() {
		Attribute("inputs")
	})
})

var _ = Resource("input", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("projects/:project_id/input"))
		Description("Add a input")
		Params(func() {
			Param("project_id", Integer)
		})
		Payload(AddInputPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Input)
		})
	})

	Action("get", func() {
		Routing(GET("projects/@/:project/inputs/@/:input"))
		Description("Add a input")
		Params(func() {
			Param("project", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Project slug")
			})
			Param("input", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Input slug")
			})
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Input)
		})
	})

	Action("get id", func() {
		Routing(GET("inputs/:input_id"))
		Description("Add a input")
		Params(func() {
			Param("input_id", Integer)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Input)
		})
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/inputs"))
		Description("List a project's inputs")
		Params(func() {
			Param("project", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Project slug")
			})
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Inputs)
		})
	})

	Action("list id", func() {
		Routing(GET("projects/:project_id/inputs"))
		Description("List a project's inputs")
		Params(func() {
			Param("project_id", Integer)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Inputs)
		})
	})
})
