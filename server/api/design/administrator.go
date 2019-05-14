package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddProjectAdministratorPayload = Type("AddAdministratorPayload", func() {
	Attribute("userId", Integer, func() {})
	Required("userId")
})

var ProjectAdministrator = MediaType("application/vnd.app.administrator+json", func() {
	TypeName("ProjectAdministrator")
	Reference(AddProjectAdministratorPayload)
	Attributes(func() {
		Attribute("projectId", Integer)
		Attribute("userId")
		Required("projectId", "userId")
	})
	View("default", func() {
		Attribute("projectId")
		Attribute("userId")
	})
})

var ProjectAdministrators = MediaType("application/vnd.app.administrators+json", func() {
	TypeName("ProjectAdministrators")
	Attributes(func() {
		Attribute("administrators", CollectionOf(ProjectAdministrator))
		Required("administrators")
	})
	View("default", func() {
		Attribute("administrators")
	})
})

var _ = Resource("administrator", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("projects/:projectId/administrators"))
		Description("Add an administrator to a project")
		Params(func() {
			Param("projectId", Integer)
		})
		Payload(AddProjectAdministratorPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(ProjectAdministrator)
		})
	})

	Action("delete", func() {
		Routing(DELETE("projects/:projectId/administrators/:userId"))
		Description("Remove an administrator from a project")
		Params(func() {
			Param("projectId", Integer)
			Param("userId", Integer)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(ProjectAdministrator)
		})
	})

	Action("get", func() {
		Routing(GET("projects/@/:project/administrators/@/:email"))
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("email", String)
		})
		Description("Get an administrator")
		Response(OK, func() {
			Media(ProjectAdministrator)
		})
	})

	Action("get id", func() {
		Routing(GET("projects/:projectId/administrators/:userId"))
		Description("Get an administrator")
		Params(func() {
			Param("projectId", Integer)
			Param("userId", Integer)
			Required("projectId", "userId")
		})
		Response(OK, func() {
			Media(ProjectAdministrator)
		})
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/administrators"))
		Description("List an project's administrators")
		Params(func() {
			Param("project", String, ProjectSlug)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(ProjectAdministrators)
		})
	})

	Action("list id", func() {
		Routing(GET("projects/:projectId/administrators"))
		Description("List an projects's administrators")
		Params(func() {
			Param("projectId", Integer)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(ProjectAdministrators)
		})
	})
})
