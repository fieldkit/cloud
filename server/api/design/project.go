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
	Attribute("goal", String)
	Attribute("location", String)
	Attribute("tags", String)
	Attribute("private", Boolean)
	Attribute("start_time", DateTime)
	Attribute("end_time", DateTime)
	Required("name", "slug", "description")
})

var InviteUserPayload = Type("InviteUserPayload", func() {
	Attribute("email", String)
	Required("email")
})

var Project = MediaType("application/vnd.app.project+json", func() {
	TypeName("Project")
	Reference(AddProjectPayload)
	Attributes(func() {
		Attribute("id", Integer)
		Attribute("name")
		Attribute("slug")
		Attribute("description")
		Attribute("goal")
		Attribute("location")
		Attribute("tags")
		Attribute("private", Boolean)
		Attribute("start_time", DateTime)
		Attribute("end_time", DateTime)
		Attribute("media_url")
		Attribute("media_content_type")
		Required("id", "name", "slug", "description", "goal", "location", "private", "tags")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("slug")
		Attribute("description")
		Attribute("goal")
		Attribute("location")
		Attribute("tags")
		Attribute("private")
		Attribute("start_time")
		Attribute("end_time")
		Attribute("media_url")
		Attribute("media_content_type")
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
		Routing(POST("projects"))
		Description("Add a project")
		Payload(AddProjectPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Project)
		})
	})

	Action("update", func() {
		Routing(PATCH("projects/:projectId"))
		Description("Update a project")
		Params(func() {
			Param("projectId", Integer)
			Required("projectId")
		})
		Payload(AddProjectPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Project)
		})
	})

	Action("get", func() {
		NoSecurity()
		Routing(GET("projects/@/:project"))
		Description("Get a project")
		Params(func() {
			Param("project", String, ProjectSlug)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Project)
		})
	})

	Action("get id", func() {
		Routing(GET("projects/:projectId"))
		Description("Get a project")
		Params(func() {
			Param("projectId", Integer)
			Required("projectId")
		})
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

	Action("save image", func() {
		Routing(POST("/projects/:projectId/media"))
		Description("Save a project image")
		Params(func() {
			Param("projectId", Integer)
			Required("projectId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Project)
		})
	})

	Action("get image", func() {
		NoSecurity()
		Routing(GET("/projects/:projectId/media"))
		Description("Get a project image")
		Params(func() {
			Param("projectId", Integer)
			Required("projectId")
		})
		Response(OK, func() {
			Media("image/png")
		})
	})

	Action("invite user", func() {
		Routing(POST("/projects/:projectId/invite"))
		Description("Invite a user to project")
		Payload(InviteUserPayload)
		Params(func() {
			Param("projectId", Integer)
			Required("projectId")
		})
		Response(BadRequest)
		Response(OK, func() {
			Status(204)
		})
	})
})
