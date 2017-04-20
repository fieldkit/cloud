package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddMemberPayload = Type("AddMemberPayload", func() {
	Attribute("user_id", Integer)
	Attribute("role", String)
	Required("user_id", "role")
})

var UpdateMemberPayload = Type("UpdateMemberPayload", func() {
	Attribute("role", String)
	Required("role")
})

var TeamMember = MediaType("application/vnd.app.member+json", func() {
	TypeName("TeamMember")
	Reference(AddMemberPayload)
	Attributes(func() {
		Attribute("team_id", Integer)
		Attribute("user_id")
		Attribute("role")
		Required("team_id", "user_id", "role")
	})
	View("default", func() {
		Attribute("team_id")
		Attribute("user_id")
		Attribute("role")
	})
})

var TeamMembers = MediaType("application/vnd.app.members+json", func() {
	TypeName("TeamMembers")
	Attributes(func() {
		Attribute("members", CollectionOf(TeamMember))
		Required("members")
	})
	View("default", func() {
		Attribute("members")
	})
})

var _ = Resource("member", func() {
	Security(JWT, func() { // Use JWT to auth requests to this endpoint
		Scope("api:access") // Enforce presence of "api" scope in JWT claims.
	})

	Action("add", func() {
		Routing(POST("teams/:team_id/members"))
		Description("Add a member to a team")
		Params(func() {
			Param("team_id", Integer)
		})
		Payload(AddMemberPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(TeamMember)
		})
	})

	Action("update", func() {
		Routing(PATCH("teams/:team_id/members/:user_id"))
		Description("Update a member")
		Params(func() {
			Param("team_id", Integer)
			Param("user_id", Integer)
			Required("team_id", "user_id")
		})
		Payload(UpdateMemberPayload)
		Response(OK, func() {
			Media(TeamMember)
		})
	})

	Action("delete", func() {
		Routing(DELETE("teams/:team_id/members/:user_id"))
		Description("Remove a member from a team")
		Params(func() {
			Param("team_id", Integer)
			Param("user_id", Integer)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(TeamMember)
		})
	})

	Action("get", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/teams/@/:team/members/@/:username"))
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("expedition", String, ExpeditionSlug)
			Param("team", String, TeamSlug)
			Param("username", String, Username)
		})
		Description("Get a member")
		Response(OK, func() {
			Media(TeamMember)
		})
	})

	Action("get id", func() {
		Routing(GET("teams/:team_id/members/:user_id"))
		Description("Get a member")
		Params(func() {
			Param("team_id", Integer)
			Param("user_id", Integer)
			Required("team_id", "user_id")
		})
		Response(OK, func() {
			Media(TeamMember)
		})
	})

	Action("list", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/teams/@/:team/members"))
		Description("List an team's members")
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("expedition", String, ExpeditionSlug)
			Param("team", String, TeamSlug)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(TeamMembers)
		})
	})

	Action("list id", func() {
		Routing(GET("teams/:team_id/members"))
		Description("List an teams's members")
		Params(func() {
			Param("team_id", Integer)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(TeamMembers)
		})
	})
})
