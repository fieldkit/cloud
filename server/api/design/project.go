package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddProjectPayload = Type("AddProjectPayload", func() {
	Attribute("name", String)
	Attribute("slug", String, func() {
		Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
		MaxLength(40)
	})
	Attribute("description", String)
	Required("name", "slug", "description")
})

var Project = MediaType("application/vnd.app.project+json", func() {
	TypeName("Project")
	Reference(AddProjectPayload)
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("name")
		Attribute("slug")
		Attribute("description")
		Required("id", "name", "slug", "description")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("slug")
		Attribute("description")
	})
})

var Projects = MediaType("application/vnd.app.projects+json", func() {
	TypeName("Projects")
	Attributes(func() {
		Attribute("projects", CollectionOf(Project))
		Required("projects")
	})
	View("default", func() {
		Attribute("projects")
	})
})

var _ = Resource("project", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("project"))
		Description("Add a project")
		Payload(AddProjectPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Project)
		})
	})

	Action("get", func() {
		Routing(GET("projects/:project"))
		Description("Get a project")
		Response(BadRequest)
		Response(OK, func() {
			Media(Project)
		})
	})

	Action("list", func() {
		Routing(GET("projects"))
		Description("List projects")
		Response(BadRequest)
		Response(OK, func() {
			Media(Projects)
		})
	})

	Action("list current", func() {
		Routing(GET("user/projects"))
		Description("List the authenticated user's projects")
		Response(BadRequest)
		Response(OK, func() {
			Media(Projects)
		})
	})
})
