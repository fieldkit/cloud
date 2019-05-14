package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AddMemberPayload = Type("AddMemberPayload", func() {
	Attribute("userId", Integer)
	Attribute("role", String)
	Required("userId", "role")
})

var UpdateMemberPayload = Type("UpdateMemberPayload", func() {
	Attribute("role", String)
	Required("role")
})

var TeamMember = MediaType("application/vnd.app.member+json", func() {
	TypeName("TeamMember")
	Reference(AddMemberPayload)
	Attributes(func() {
		Attribute("teamId", Integer)
		Attribute("userId")
		Attribute("role")
		Required("teamId", "userId", "role")
	})
	View("default", func() {
		Attribute("teamId")
		Attribute("userId")
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
		Routing(POST("teams/:teamId/members"))
		Description("Add a member to a team")
		Params(func() {
			Param("teamId", Integer)
		})
		Payload(AddMemberPayload)
		Response(BadRequest)
		Response(OK, func() {
			Media(TeamMember)
		})
	})

	Action("update", func() {
		Routing(PATCH("teams/:teamId/members/:userId"))
		Description("Update a member")
		Params(func() {
			Param("teamId", Integer)
			Param("userId", Integer)
			Required("teamId", "userId")
		})
		Payload(UpdateMemberPayload)
		Response(OK, func() {
			Media(TeamMember)
		})
	})

	Action("delete", func() {
		Routing(DELETE("teams/:teamId/members/:userId"))
		Description("Remove a member from a team")
		Params(func() {
			Param("teamId", Integer)
			Param("userId", Integer)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(TeamMember)
		})
	})

	Action("get", func() {
		Routing(GET("projects/@/:project/expeditions/@/:expedition/teams/@/:team/members/@/:email"))
		Params(func() {
			Param("project", String, ProjectSlug)
			Param("expedition", String, ExpeditionSlug)
			Param("team", String, TeamSlug)
			Param("email", String)
		})
		Description("Get a member")
		Response(OK, func() {
			Media(TeamMember)
		})
	})

	Action("get id", func() {
		Routing(GET("teams/:teamId/members/:userId"))
		Description("Get a member")
		Params(func() {
			Param("teamId", Integer)
			Param("userId", Integer)
			Required("teamId", "userId")
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
		Routing(GET("teams/:teamId/members"))
		Description("List an teams's members")
		Params(func() {
			Param("teamId", Integer)
		})
		Response(BadRequest)
		Response(OK, func() {
			Media(TeamMembers)
		})
	})
})
