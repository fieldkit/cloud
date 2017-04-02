package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddTeamPayload = Type("AddTeamPayload", func() {
	Attribute("name", String, func() {
		Pattern(`\S`)
		MaxLength(256)
	})
	Attribute("slug", String, func() {
		Pattern(`^[[:alnum:]]+(-[[:alnum:]]+)*$`)
		MaxLength(40)
	})
	Attribute("description", String)
	Required("name", "slug", "description")
})

var Team = MediaType("application/vnd.app.team+json", func() {
	TypeName("Team")
	Reference(AddTeamPayload)
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

var Teams = MediaType("application/vnd.app.teams+json", func() {
	TypeName("Teams")
	Attributes(func() {
		Attribute("teams", CollectionOf(Team))
		Required("teams")
	})
	View("default", func() {
		Attribute("teams")
	})
})

var _ = Resource("team", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("expeditions/:expedition_id/teams"))
		Description("Add a team")
		Params(func() {
			Param("expedition_id", Integer)
			Required("expeditions_id")
		})
		Payload(AddTeamPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Team)
		})
	})

	Action("update", func() {
		Routing(PATCH("teams/:team_id"))
		Description("Update a team")
		Params(func() {
			Param("team_id", Integer)
			Required("team_id")
		})
		Payload(AddTeamPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(Team)
		})
	})

	Action("delete", func() {
		Routing(DELETE("teams/:team_id"))
		Description("Update a team")
		Params(func() {
			Param("team_id", Integer)
			Required("team_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Team)
		})
	})

	Action("get", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/teams/@/:team"))
		Description("Add a team")
		Params(func() {
			Param("project", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Project slug")
			})
			Param("expedition", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Expedition slug")
			})
			Param("team", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Project slug")
			})
			Required("project", "expedition", "team")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Team)
		})
	})

	Action("get id", func() {
		Routing(GET("teams/:team_id"))
		Description("Add a team")
		Params(func() {
			Param("team_id", Integer)
			Required("team_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Team)
		})
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/teams"))
		Description("List an expedition's teams")
		Params(func() {
			Param("project", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Project slug")
			})
			Param("expedition", String, func() {
				Pattern("^[[:alnum:]]+(-[[:alnum:]]+)*$")
				Description("Expedition slug")
			})
			Required("project", "expedition")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Teams)
		})
	})

	Action("list id", func() {
		Routing(GET("expeditions/:expedition_id/teams"))
		Description("List an expedition's teams")
		Params(func() {
			Param("expedition_id", Integer)
			Required("expedition_id")
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(Teams)
		})
	})
})
