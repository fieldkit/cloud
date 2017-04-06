package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddProjectAdministratorPayload = Type("AddAdministratorPayload", func() {
	Attribute("user_id", Integer, func() {})
	Required("user_id")
})

var ProjectAdministrator = MediaType("application/vnd.app.administrator+json", func() {
	TypeName("ProjectAdministrator")
	Reference(AddProjectAdministratorPayload)
	Attributes(func() {
		Attribute("project_id", Integer)
		Attribute("user_id")
		Required("project_id", "user_id")
	})
	View("default", func() {
		Attribute("project_id")
		Attribute("user_id")
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
		Routing(POST("projects/:project_id/administrators"))
		Description("Add an administrator to a project")
		Params(func() {
			Param("project_id", Integer)
		})
		Payload(AddProjectAdministratorPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(ProjectAdministrator)
		})
	})

	Action("delete", func() {
		Routing(DELETE("projects/:project_id/administrators/:user_id"))
		Description("Remove an administrator from a project")
		Params(func() {
			Param("project_id", Integer)
			Param("user_id", Integer)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(ProjectAdministrator)
		})
	})

	Action("get", func() {
		Routing(GET("projects/@/:project/administrators/@/:username"))
		Params(func() {
			Param("project", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Project slug")
			})
			Param("username", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Username")
			})
		})
		Description("Get an administrator")
		Response(OK, func() {
			Media(ProjectAdministrator)
		})
	})

	Action("get id", func() {
		Routing(GET("projects/:project_id/administrators/:user_id"))
		Description("Get an administrator")
		Params(func() {
			Param("project_id", Integer)
			Param("user_id", Integer)
			Required("project_id", "user_id")
		})
		Response(OK, func() {
			Media(ProjectAdministrator)
		})
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/administrators"))
		Description("List an project's administrators")
		Params(func() {
			Param("project", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Project slug")
			})
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(ProjectAdministrators)
		})
	})

	Action("list id", func() {
		Routing(GET("projects/:project_id/administrators"))
		Description("List an projects's administrators")
		Params(func() {
			Param("project_id", Integer)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(ProjectAdministrators)
		})
	})
})
