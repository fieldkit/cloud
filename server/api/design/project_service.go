package design

import (
	. "goa.design/goa/v3/dsl"
)

var PendingInvite = Type("PendingInvite", func() {
	Attribute("id", Int64)
	Attribute("project", ProjectSummary)
	Attribute("time", Int64)
	Attribute("role", Int32)
	Required("id")
	Required("project")
	Required("time")
	Required("role")
})

var PendingInvites = ResultType("application/vnd.app.invites.pending", func() {
	TypeName("PendingInvites")
	Attributes(func() {
		Attribute("pending", ArrayOf(PendingInvite))
		Required("pending")
	})
	View("default", func() {
		Attribute("pending")
	})
})

var ProjectUpdate = ResultType("application/vnd.app.project.update", func() {
	TypeName("ProjectUpdate")
	Attributes(func() {
		Attribute("id", Int64)
		Attribute("body", String)
		Attribute("created_at", Int64)
		Required("id")
		Required("body")
		Required("created_at")
	})
	View("default", func() {
		Attribute("id")
		Attribute("body")
	})
})

var _ = Service("project", func() {
	Method("add update", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("projectId", Int32)
			Attribute("body", String)
			Required("auth")
			Required("projectId")
			Required("body")
		})

		Result(ProjectUpdate)

		HTTP(func() {
			POST("projects/{projectId}/updates")

			httpAuthentication()
		})
	})

	Method("delete update", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("projectId", Int32)
			Attribute("updateId", Int64)
			Required("auth")
			Required("projectId")
			Required("updateId")
		})

		Result(Empty)

		HTTP(func() {
			DELETE("projects/{projectId}/updates/{updateId}")

			httpAuthentication()
		})
	})

	Method("modify update", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Attribute("projectId", Int32)
			Attribute("updateId", Int64)
			Attribute("body", String)
			Required("auth")
			Required("projectId")
			Required("updateId")
			Required("body")
		})

		Result(ProjectUpdate)

		HTTP(func() {
			POST("projects/{projectId}/updates/{updateId}")

			httpAuthentication()
		})
	})

	Method("invites", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
		})

		Result(PendingInvites)

		HTTP(func() {
			GET("projects/invites/pending")

			httpAuthentication()
		})
	})

	Method("lookup invite", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Token("token")
			Required("token")
		})

		Result(PendingInvites)

		HTTP(func() {
			GET("projects/invites/{token}")

			httpAuthentication()
		})
	})

	Method("accept invite", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("id", Int64)
			Required("id")
			Attribute("token", String)
		})

		HTTP(func() {
			POST("projects/invites/{id}/accept")

			Params(func() {
				Param("token")
			})

			httpAuthentication()
		})
	})

	Method("reject invite", func() {
		Security(JWTAuth, func() {
			Scope("api:access")
		})

		Payload(func() {
			Token("auth")
			Required("auth")
			Attribute("id", Int64)
			Required("id")
			Attribute("token", String)
		})

		HTTP(func() {
			POST("projects/invites/{id}/reject")

			Params(func() {
				Param("token")
			})

			httpAuthentication()
		})
	})

	commonOptions()
})
